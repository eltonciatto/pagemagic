use actix_web::{web, HttpResponse, Result};
use crate::models::{UsageEvent, EventBatch, MeterRecord};
use crate::services::MeterService;
use tracing::{info, error};
use std::sync::Arc;

pub async fn health() -> Result<HttpResponse> {
    Ok(HttpResponse::Ok().json(serde_json::json!({
        "status": "healthy",
        "service": "meter-svc"
    })))
}

pub async fn receive_event(
    event: web::Json<UsageEvent>,
    meter_service: web::Data<Arc<MeterService>>,
) -> Result<HttpResponse> {
    info!("Received usage event: {:?}", event.event_type);
    
    match meter_service.process_event(&event).await {
        Ok(_) => {
            Ok(HttpResponse::Ok().json(serde_json::json!({
                "message": "Event processed successfully"
            })))
        },
        Err(e) => {
            error!("Failed to process event: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to process event"
            })))
        }
    }
}

pub async fn receive_batch(
    batch: web::Json<EventBatch>,
    meter_service: web::Data<Arc<MeterService>>,
) -> Result<HttpResponse> {
    info!("Received batch of {} events", batch.events.len());
    
    match meter_service.process_batch(&batch).await {
        Ok(processed) => {
            Ok(HttpResponse::Ok().json(serde_json::json!({
                "message": "Batch processed successfully",
                "processed_count": processed
            })))
        },
        Err(e) => {
            error!("Failed to process batch: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to process batch"
            })))
        }
    }
}

pub async fn get_meters(
    meter_service: web::Data<Arc<MeterService>>,
) -> Result<HttpResponse> {
    match meter_service.get_active_meters().await {
        Ok(meters) => Ok(HttpResponse::Ok().json(meters)),
        Err(e) => {
            error!("Failed to get meters: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to get meters"
            })))
        }
    }
}

pub async fn get_user_usage(
    path: web::Path<String>,
    meter_service: web::Data<Arc<MeterService>>,
) -> Result<HttpResponse> {
    let user_id = path.into_inner();
    
    match meter_service.get_user_usage(&user_id).await {
        Ok(usage) => Ok(HttpResponse::Ok().json(usage)),
        Err(e) => {
            error!("Failed to get user usage for {}: {}", user_id, e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to get user usage"
            })))
        }
    }
}

pub async fn force_sync(
    meter_service: web::Data<Arc<MeterService>>,
) -> Result<HttpResponse> {
    info!("Forcing sync to Stripe");
    
    match meter_service.sync_to_stripe().await {
        Ok(synced) => {
            Ok(HttpResponse::Ok().json(serde_json::json!({
                "message": "Sync completed",
                "synced_records": synced
            })))
        },
        Err(e) => {
            error!("Failed to sync to Stripe: {}", e);
            Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to sync to Stripe"
            })))
        }
    }
}
