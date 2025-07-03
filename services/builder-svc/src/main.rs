use anyhow::Result;
use axum::{
    extract::{State, WebSocketUpgrade},
    http::Method,
    middleware,
    response::Response,
    routing::{get, post},
    Router,
};
use clap::Parser;
use config::Config;
use dotenv::dotenv;
use metrics_exporter_prometheus::PrometheusBuilder;
use serde::Deserialize;
use sqlx::postgres::PgPoolOptions;
use std::{net::SocketAddr, sync::Arc, time::Duration};
use tokio::signal;
use tower::ServiceBuilder;
use tower_http::{
    compression::CompressionLayer,
    cors::{Any, CorsLayer},
    limit::RequestBodyLimitLayer,
    timeout::TimeoutLayer,
    trace::TraceLayer,
};
use tracing::{info, warn};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod config_manager;
mod handlers;
mod middleware_custom;
mod models;
mod services;
mod utils;
mod websocket;

use crate::{
    config_manager::AppConfig,
    handlers::{
        analytics, auth, builder, health, projects, templates, themes, version_control,
    },
    middleware_custom::{metrics_middleware, rate_limit_middleware},
    services::{
        analytics_service::AnalyticsService,
        ast_service::ASTService,
        auth_service::AuthService,
        builder_service::BuilderService,
        cache_service::CacheService,
        code_gen_service::CodeGenService,
        git_service::GitService,
        notification_service::NotificationService,
        optimization_service::OptimizationService,
        project_service::ProjectService,
        template_service::TemplateService,
        theme_service::ThemeService,
        websocket_service::WebSocketService,
    },
    websocket::websocket_handler,
};

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// Configuration file path
    #[arg(short, long, default_value = "configs/config.yaml")]
    config: String,
    
    /// Port to bind to
    #[arg(short, long)]
    port: Option<u16>,
    
    /// Environment
    #[arg(short, long, default_value = "development")]
    env: String,
    
    /// Enable debug mode
    #[arg(short, long)]
    debug: bool,
}

#[derive(Clone)]
pub struct AppState {
    pub config: Arc<AppConfig>,
    pub db: sqlx::PgPool,
    pub redis: redis::aio::ConnectionManager,
    pub ast_service: Arc<ASTService>,
    pub builder_service: Arc<BuilderService>,
    pub project_service: Arc<ProjectService>,
    pub template_service: Arc<TemplateService>,
    pub theme_service: Arc<ThemeService>,
    pub auth_service: Arc<AuthService>,
    pub git_service: Arc<GitService>,
    pub code_gen_service: Arc<CodeGenService>,
    pub optimization_service: Arc<OptimizationService>,
    pub analytics_service: Arc<AnalyticsService>,
    pub notification_service: Arc<NotificationService>,
    pub websocket_service: Arc<WebSocketService>,
    pub cache_service: Arc<CacheService>,
}

#[tokio::main]
async fn main() -> Result<()> {
    // Parse command line arguments
    let args = Args::parse();
    
    // Load environment variables
    dotenv().ok();
    
    // Initialize configuration
    let config = Arc::new(AppConfig::new(&args.config, &args.env)?);
    
    // Initialize tracing
    init_tracing(&config).await?;
    
    info!("Starting Page Magic Builder Service v{}", env!("CARGO_PKG_VERSION"));
    info!("Environment: {}", args.env);
    info!("Configuration loaded from: {}", args.config);
    
    // Initialize metrics
    init_metrics(&config).await?;
    
    // Initialize database
    let db = init_database(&config).await?;
    
    // Initialize Redis
    let redis = init_redis(&config).await?;
    
    // Initialize services
    let services = init_services(&config, &db, &redis).await?;
    
    // Create application state
    let state = AppState {
        config: config.clone(),
        db,
        redis,
        ast_service: services.ast_service,
        builder_service: services.builder_service,
        project_service: services.project_service,
        template_service: services.template_service,
        theme_service: services.theme_service,
        auth_service: services.auth_service,
        git_service: services.git_service,
        code_gen_service: services.code_gen_service,
        optimization_service: services.optimization_service,
        analytics_service: services.analytics_service,
        notification_service: services.notification_service,
        websocket_service: services.websocket_service,
        cache_service: services.cache_service,
    };
    
    // Build application router
    let app = build_router(state.clone()).await?;
    
    // Determine port
    let port = args.port.unwrap_or(config.server.port);
    let addr = SocketAddr::from(([0, 0, 0, 0], port));
    
    info!("Server starting on {}", addr);
    
    // Start server with graceful shutdown
    let listener = tokio::net::TcpListener::bind(addr).await?;
    
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal())
        .await?;
    
    info!("Server shutting down gracefully");
    Ok(())
}

async fn init_tracing(config: &AppConfig) -> Result<()> {
    let env_filter = std::env::var("RUST_LOG")
        .unwrap_or_else(|_| format!("{}=debug,tower_http=debug", env!("CARGO_PKG_NAME")));
    
    let subscriber = tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| env_filter.into()),
        )
        .with(tracing_subscriber::fmt::layer());
    
    // Add OpenTelemetry if configured
    if config.observability.jaeger.enabled {
        let tracer = opentelemetry_jaeger::new_agent_pipeline()
            .with_service_name("builder-svc")
            .with_endpoint(&config.observability.jaeger.endpoint)
            .install_simple()?;
        
        let telemetry = tracing_opentelemetry::layer().with_tracer(tracer);
        subscriber.with(telemetry).init();
    } else {
        subscriber.init();
    }
    
    Ok(())
}

async fn init_metrics(config: &AppConfig) -> Result<()> {
    if config.observability.metrics.enabled {
        let builder = PrometheusBuilder::new();
        let handle = builder
            .with_http_listener(([0, 0, 0, 0], config.observability.metrics.port))
            .install()?;
        
        tokio::spawn(async move {
            handle;
        });
        
        info!("Metrics server started on port {}", config.observability.metrics.port);
    }
    
    Ok(())
}

async fn init_database(config: &AppConfig) -> Result<sqlx::PgPool> {
    info!("Connecting to database...");
    
    let pool = PgPoolOptions::new()
        .max_connections(config.database.max_connections)
        .acquire_timeout(Duration::from_secs(config.database.acquire_timeout))
        .connect(&config.database.url)
        .await?;
    
    // Run migrations
    sqlx::migrate!("./migrations").run(&pool).await?;
    
    info!("Database connected and migrations applied");
    Ok(pool)
}

async fn init_redis(config: &AppConfig) -> Result<redis::aio::ConnectionManager> {
    info!("Connecting to Redis...");
    
    let client = redis::Client::open(config.redis.url.as_str())?;
    let conn = redis::aio::ConnectionManager::new(client).await?;
    
    info!("Redis connected");
    Ok(conn)
}

struct Services {
    ast_service: Arc<ASTService>,
    builder_service: Arc<BuilderService>,
    project_service: Arc<ProjectService>,
    template_service: Arc<TemplateService>,
    theme_service: Arc<ThemeService>,
    auth_service: Arc<AuthService>,
    git_service: Arc<GitService>,
    code_gen_service: Arc<CodeGenService>,
    optimization_service: Arc<OptimizationService>,
    analytics_service: Arc<AnalyticsService>,
    notification_service: Arc<NotificationService>,
    websocket_service: Arc<WebSocketService>,
    cache_service: Arc<CacheService>,
}

async fn init_services(
    config: &AppConfig,
    db: &sqlx::PgPool,
    redis: &redis::aio::ConnectionManager,
) -> Result<Services> {
    info!("Initializing services...");
    
    // Initialize cache service first (other services depend on it)
    let cache_service = Arc::new(CacheService::new(redis.clone()).await?);
    
    // Initialize core services
    let config_arc = Arc::new(config.clone());
    let ast_service = Arc::new(ASTService::new(config_arc.clone(), cache_service.clone()).await?);
    let auth_service = Arc::new(AuthService::new(config_arc.clone(), db.clone()).await?);
    let git_service = Arc::new(GitService::new(config_arc.clone()).await?);
    let code_gen_service = Arc::new(CodeGenService::new(config_arc.clone(), cache_service.clone()).await?);
    let optimization_service = Arc::new(OptimizationService::new(config_arc.clone()).await?);
    let analytics_service = Arc::new(AnalyticsService::new(config_arc.clone(), db.clone()).await?);
    let notification_service = Arc::new(NotificationService::new(config_arc.clone()).await?);
    let websocket_service = Arc::new(WebSocketService::new(config_arc.clone()).await?);
    
    // Initialize business logic services
    let template_service = Arc::new(
        TemplateService::new(
            config_arc.clone(),
            db.clone(),
            cache_service.clone(),
            code_gen_service.clone(),
        ).await?
    );
    
    let theme_service = Arc::new(
        ThemeService::new(
            config_arc.clone(),
            db.clone(),
            cache_service.clone(),
        ).await?
    );
    
    let project_service = Arc::new(
        ProjectService::new(
            config_arc.clone(),
            db.clone(),
            cache_service.clone(),
            git_service.clone(),
            analytics_service.clone(),
        ).await?
    );
    
    let builder_service = Arc::new(
        BuilderService::new(
            config_arc.clone(),
            ast_service.clone(),
            template_service.clone(),
            theme_service.clone(),
            code_gen_service.clone(),
            optimization_service.clone(),
            websocket_service.clone(),
        ).await?
    );
    
    info!("All services initialized successfully");
    
    Ok(Services {
        ast_service,
        builder_service,
        project_service,
        template_service,
        theme_service,
        auth_service,
        git_service,
        code_gen_service,
        optimization_service,
        analytics_service,
        notification_service,
        websocket_service,
        cache_service,
    })
}

async fn build_router(state: AppState) -> Result<Router> {
    // Create CORS layer
    let cors = CorsLayer::new()
        .allow_methods([Method::GET, Method::POST, Method::PUT, Method::DELETE, Method::PATCH])
        .allow_headers(Any)
        .allow_origin(Any);
    
    // Build main router
    let app = Router::new()
        // Health and system endpoints
        .route("/health", get(health::health_check))
        .route("/health/ready", get(health::readiness_check))
        .route("/health/live", get(health::liveness_check))
        .route("/metrics", get(health::metrics))
        .route("/version", get(health::version))
        
        // Authentication endpoints
        .route("/auth/login", post(auth::login))
        .route("/auth/refresh", post(auth::refresh_token))
        .route("/auth/validate", post(auth::validate_token))
        
        // Project management
        .nest("/api/v1/projects", projects::router())
        
        // Builder endpoints
        .nest("/api/v1/builder", builder::router())
        
        // Template management
        .nest("/api/v1/templates", templates::router())
        
        // Theme management
        .nest("/api/v1/themes", themes::router())
        
        // Version control
        .nest("/api/v1/git", version_control::router())
        
        // Analytics
        .nest("/api/v1/analytics", analytics::router())
        
        // WebSocket for real-time features
        .route("/ws", get(websocket_upgrade))
        
        // Static file serving
        .nest_service("/static", tower_http::services::ServeDir::new("static"))
        
        // Apply middleware and state
        .layer(TraceLayer::new_for_http())
        .layer(CompressionLayer::new())
        .layer(TimeoutLayer::new(Duration::from_secs(30)))
        .layer(RequestBodyLimitLayer::new(100 * 1024 * 1024)) // 100MB
        .layer(cors)
        .with_state(state);
    
    Ok(app)
}

async fn websocket_upgrade(
    ws: WebSocketUpgrade,
    State(state): State<AppState>,
) -> Response {
    ws.on_upgrade(move |socket| websocket_handler(socket, state))
}

async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };
    
    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };
    
    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();
    
    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }
    
    warn!("Shutdown signal received");
}
