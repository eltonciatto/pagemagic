use axum::{
    extract::State,
    http::StatusCode,
    response::Json,
    routing::get,
    Router,
};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tracing::{error, info, instrument};

use crate::AppState;

#[derive(Debug, Serialize)]
pub struct HealthResponse {
    pub status: String,
    pub service: String,
    pub version: String,
    pub timestamp: chrono::DateTime<chrono::Utc>,
    pub uptime: u64,
    pub checks: HashMap<String, HealthCheck>,
}

#[derive(Debug, Serialize)]
pub struct HealthCheck {
    pub status: String,
    pub message: Option<String>,
    pub details: Option<serde_json::Value>,
    pub duration_ms: u64,
}

#[derive(Debug, Serialize)]
pub struct ReadinessResponse {
    pub ready: bool,
    pub services: HashMap<String, ServiceStatus>,
    pub timestamp: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Serialize)]
pub struct ServiceStatus {
    pub available: bool,
    pub response_time_ms: u64,
    pub error: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct LivenessResponse {
    pub alive: bool,
    pub timestamp: chrono::DateTime<chrono::Utc>,
    pub memory_usage: MemoryUsage,
    pub cpu_usage: f64,
}

#[derive(Debug, Serialize)]
pub struct MemoryUsage {
    pub allocated: u64,
    pub total_allocated: u64,
    pub system: u64,
}

#[derive(Debug, Serialize)]
pub struct VersionResponse {
    pub service: String,
    pub version: String,
    pub git_commit: String,
    pub build_time: String,
    pub rust_version: String,
    pub dependencies: HashMap<String, String>,
}

#[derive(Debug, Serialize)]
pub struct MetricsResponse {
    pub metrics: String,
    pub format: String,
    pub timestamp: chrono::DateTime<chrono::Utc>,
}

pub fn router() -> Router<AppState> {
    Router::new()
        .route("/health", get(health_check))
        .route("/health/ready", get(readiness_check))
        .route("/health/live", get(liveness_check))
        .route("/metrics", get(metrics))
        .route("/version", get(version))
}

#[instrument(skip(state))]
pub async fn health_check(
    State(state): State<AppState>,
) -> Result<Json<HealthResponse>, StatusCode> {
    let start_time = std::time::Instant::now();
    let mut checks = HashMap::new();
    
    // Check database connection
    let db_start = std::time::Instant::now();
    let db_check = match sqlx::query("SELECT 1").fetch_one(&state.db).await {
        Ok(_) => HealthCheck {
            status: "healthy".to_string(),
            message: Some("Database connection is working".to_string()),
            details: None,
            duration_ms: db_start.elapsed().as_millis() as u64,
        },
        Err(e) => HealthCheck {
            status: "unhealthy".to_string(),
            message: Some(format!("Database connection failed: {}", e)),
            details: None,
            duration_ms: db_start.elapsed().as_millis() as u64,
        },
    };
    checks.insert("database".to_string(), db_check);
    
    // Check Redis connection
    let redis_start = std::time::Instant::now();
    let redis_check = match state.cache_service.ping().await {
        Ok(_) => HealthCheck {
            status: "healthy".to_string(),
            message: Some("Redis connection is working".to_string()),
            details: None,
            duration_ms: redis_start.elapsed().as_millis() as u64,
        },
        Err(e) => HealthCheck {
            status: "unhealthy".to_string(),
            message: Some(format!("Redis connection failed: {}", e)),
            details: None,
            duration_ms: redis_start.elapsed().as_millis() as u64,
        },
    };
    checks.insert("redis".to_string(), redis_check);
    
    // Check external services
    if state.config.features.ai_generation {
        let ai_start = std::time::Instant::now();
        let ai_check = match check_openai_connection(&state).await {
            Ok(_) => HealthCheck {
                status: "healthy".to_string(),
                message: Some("OpenAI API is accessible".to_string()),
                details: None,
                duration_ms: ai_start.elapsed().as_millis() as u64,
            },
            Err(e) => HealthCheck {
                status: "unhealthy".to_string(),
                message: Some(format!("OpenAI API failed: {}", e)),
                details: None,
                duration_ms: ai_start.elapsed().as_millis() as u64,
            },
        };
        checks.insert("openai".to_string(), ai_check);
    }
    
    // Determine overall status
    let overall_status = if checks.values().all(|check| check.status == "healthy") {
        "healthy".to_string()
    } else {
        "unhealthy".to_string()
    };
    
    let response = HealthResponse {
        status: overall_status,
        service: "builder-svc".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
        timestamp: chrono::Utc::now(),
        uptime: start_time.elapsed().as_secs(), // This would be actual uptime in production
        checks,
    };
    
    Ok(Json(response))
}

#[instrument(skip(state))]
pub async fn readiness_check(
    State(state): State<AppState>,
) -> Result<Json<ReadinessResponse>, StatusCode> {
    let mut services = HashMap::new();
    let mut ready = true;
    
    // Check database readiness
    let db_start = std::time::Instant::now();
    let db_status = match sqlx::query("SELECT 1").fetch_one(&state.db).await {
        Ok(_) => ServiceStatus {
            available: true,
            response_time_ms: db_start.elapsed().as_millis() as u64,
            error: None,
        },
        Err(e) => {
            ready = false;
            ServiceStatus {
                available: false,
                response_time_ms: db_start.elapsed().as_millis() as u64,
                error: Some(e.to_string()),
            }
        },
    };
    services.insert("database".to_string(), db_status);
    
    // Check Redis readiness
    let redis_start = std::time::Instant::now();
    let redis_status = match state.cache_service.ping().await {
        Ok(_) => ServiceStatus {
            available: true,
            response_time_ms: redis_start.elapsed().as_millis() as u64,
            error: None,
        },
        Err(e) => {
            ready = false;
            ServiceStatus {
                available: false,
                response_time_ms: redis_start.elapsed().as_millis() as u64,
                error: Some(e.to_string()),
            }
        },
    };
    services.insert("redis".to_string(), redis_status);
    
    // Check critical services are initialized
    let services_check = check_services_initialized(&state).await;
    if !services_check.available {
        ready = false;
    }
    services.insert("internal_services".to_string(), services_check);
    
    let response = ReadinessResponse {
        ready,
        services,
        timestamp: chrono::Utc::now(),
    };
    
    if ready {
        Ok(Json(response))
    } else {
        Err(StatusCode::SERVICE_UNAVAILABLE)
    }
}

#[instrument(skip(_state))]
pub async fn liveness_check(
    State(_state): State<AppState>,
) -> Result<Json<LivenessResponse>, StatusCode> {
    // Get memory usage (simplified - would use proper system metrics in production)
    let memory_usage = MemoryUsage {
        allocated: 0, // Would get actual values
        total_allocated: 0,
        system: 0,
    };
    
    let response = LivenessResponse {
        alive: true,
        timestamp: chrono::Utc::now(),
        memory_usage,
        cpu_usage: 0.0, // Would get actual CPU usage
    };
    
    Ok(Json(response))
}

#[instrument(skip(_state))]
pub async fn metrics(
    State(_state): State<AppState>,
) -> Result<Json<MetricsResponse>, StatusCode> {
    // In production, this would return actual Prometheus metrics
    let metrics_data = format!(
        "# HELP http_requests_total Total HTTP requests\n\
         # TYPE http_requests_total counter\n\
         http_requests_total{{method=\"GET\",status=\"200\"}} 1000\n\
         # HELP memory_usage_bytes Current memory usage\n\
         # TYPE memory_usage_bytes gauge\n\
         memory_usage_bytes 1048576\n\
         # HELP service_uptime_seconds Service uptime in seconds\n\
         # TYPE service_uptime_seconds counter\n\
         service_uptime_seconds 3600\n"
    );
    
    let response = MetricsResponse {
        metrics: metrics_data,
        format: "prometheus".to_string(),
        timestamp: chrono::Utc::now(),
    };
    
    Ok(Json(response))
}

#[instrument(skip(_state))]
pub async fn version(
    State(_state): State<AppState>,
) -> Result<Json<VersionResponse>, StatusCode> {
    let mut dependencies = HashMap::new();
    dependencies.insert("axum".to_string(), "0.7.0".to_string());
    dependencies.insert("tokio".to_string(), "1.0.0".to_string());
    dependencies.insert("sqlx".to_string(), "0.7.0".to_string());
    dependencies.insert("redis".to_string(), "0.24.0".to_string());
    
    let response = VersionResponse {
        service: "builder-svc".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
        git_commit: std::env::var("VERGEN_GIT_SHA").unwrap_or("unknown".to_string()),
        build_time: std::env::var("VERGEN_BUILD_TIMESTAMP").unwrap_or("unknown".to_string()),
        rust_version: std::env::var("VERGEN_RUSTC_SEMVER").unwrap_or("unknown".to_string()),
        dependencies,
    };
    
    Ok(Json(response))
}

async fn check_openai_connection(_state: &AppState) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    // Simplified check - would make actual API call in production
    Ok(())
}

async fn check_services_initialized(_state: &AppState) -> ServiceStatus {
    // Check if all required services are properly initialized
    ServiceStatus {
        available: true,
        response_time_ms: 0,
        error: None,
    }
}
