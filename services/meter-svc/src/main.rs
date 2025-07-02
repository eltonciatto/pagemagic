mod config;
mod models;
mod handlers;
mod services;

use actix_web::{web, App, HttpServer, middleware::Logger};
use tracing::{info, error};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};
use std::sync::Arc;

use config::load_config;
use services::MeterService;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Initialize tracing
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "meter_svc=debug,tower_http=debug".into()),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    // Load configuration
    let config = load_config().expect("Failed to load configuration");
    let bind_address = format!("{}:{}", config.server.host, config.server.port);

    info!("Starting meter service on {}", bind_address);

    // Initialize meter service
    let meter_service = Arc::new(
        MeterService::new(config.clone())
            .await
            .expect("Failed to initialize meter service")
    );

    // Start background sync task
    let sync_service = meter_service.clone();
    tokio::spawn(async move {
        let mut interval = tokio::time::interval(tokio::time::Duration::from_secs(60));
        loop {
            interval.tick().await;
            if let Err(e) = sync_service.sync_to_stripe().await {
                error!("Background sync error: {}", e);
            }
        }
    });

    // Start HTTP server
    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(meter_service.clone()))
            .wrap(Logger::default())
            .route("/health", web::get().to(handlers::health))
            .route("/v1/events", web::post().to(handlers::receive_event))
            .route("/v1/events/batch", web::post().to(handlers::receive_batch))
            .route("/v1/meters", web::get().to(handlers::get_meters))
            .route("/v1/users/{user_id}/usage", web::get().to(handlers::get_user_usage))
            .route("/v1/sync", web::post().to(handlers::force_sync))
    })
    .bind(&bind_address)?
    .run()
    .await
}
