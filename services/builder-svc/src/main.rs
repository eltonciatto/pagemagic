use axum::{
    extract::State,
    http::StatusCode,
    response::Json,
    routing::{get, post},
    Router,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::net::TcpListener;
use tower_http::{cors::CorsLayer, trace::TraceLayer};
use tracing::{info, instrument};
use uuid::Uuid;

mod config;
mod handlers;
mod models;
mod services;
mod utils;

use config::Config;
use handlers::*;
use services::BuilderService;

#[derive(Clone)]
pub struct AppState {
    builder_service: Arc<BuilderService>,
    config: Arc<Config>,
}

#[derive(Serialize)]
struct HealthResponse {
    status: String,
    service: String,
    version: String,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Inicializar tracing
    tracing_subscriber::fmt()
        .with_env_filter(tracing_subscriber::EnvFilter::from_default_env())
        .init();

    // Carregar configuração
    let config = Arc::new(Config::load()?);
    
    // Inicializar serviços
    let builder_service = Arc::new(BuilderService::new(config.clone()).await?);
    
    // Estado da aplicação
    let state = AppState {
        builder_service,
        config: config.clone(),
    };

    // Configurar rotas
    let app = Router::new()
        .route("/health", get(health_check))
        .route("/api/v1/build/website", post(build_website))
        .route("/api/v1/build/component", post(build_component))
        .route("/api/v1/build/static", post(build_static))
        .route("/api/v1/build/status/:build_id", get(get_build_status))
        .with_state(state)
        .layer(CorsLayer::permissive())
        .layer(TraceLayer::new_for_http());

    // Inicializar servidor
    let listener = TcpListener::bind(&format!("{}:{}", config.server.host, config.server.port)).await?;
    
    info!(
        "Builder service starting on {}:{}",
        config.server.host, config.server.port
    );

    axum::serve(listener, app).await?;

    Ok(())
}

#[instrument]
async fn health_check(State(_): State<AppState>) -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "ok".to_string(),
        service: "builder-svc".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
    })
}
