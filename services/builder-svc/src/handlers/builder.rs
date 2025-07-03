use axum::{
    extract::{Path, Query, State},
    http::StatusCode,
    response::Json,
    routing::{get, post},
    Router,
};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tracing::{error, info, instrument};
use uuid::Uuid;
use validat        Err(e) => {
            error!("Failed to get build status for {}: {}", build_id, e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to get build status",
                    "message": e.to_string()
                })),
            ))
        }rate::{
    models::{
        BuildStatus, BuildJob
    },
    AppState,
};

#[derive(Debug, Deserialize, Validate)]
pub struct BuildWebsiteRequest {
    pub project_id: Uuid,
    pub user_id: Uuid,
    pub build_request: BuildRequest,
}

#[derive(Debug, Deserialize)]
pub struct BuildStatusQuery {
    pub include_logs: Option<bool>,
    pub include_metrics: Option<bool>,
}

#[derive(Debug, Deserialize)]
pub struct ComponentBuildRequest {
    pub component_type: String,
    pub props: serde_json::Value,
    pub framework: String,
}

#[derive(Debug, Serialize)]
pub struct BuildResponse {
    pub build_id: Uuid,
    pub status: BuildStatus,
    pub message: String,
    pub estimated_time: Option<u64>,
}

#[derive(Debug, Serialize)]
pub struct ComponentResponse {
    pub component_id: String,
    pub html: String,
    pub css: String,
    pub js: Option<String>,
    pub props: serde_json::Value,
}

#[derive(Debug, Serialize)]
pub struct TemplatesResponse {
    pub templates: Vec<ComponentTemplate>,
    pub total: usize,
    pub categories: Vec<String>,
}

#[derive(Debug, Deserialize)]
pub struct TemplateQuery {
    pub category: Option<String>,
    pub framework: Option<String>,
    pub tags: Option<String>,
    pub page: Option<u32>,
    pub limit: Option<u32>,
}

#[derive(Debug, Deserialize)]
pub struct ComponentTemplateQuery {
    pub category: Option<String>,
    pub framework: Option<String>,
    pub tags: Option<String>,
    pub page: Option<u32>,
    pub per_page: Option<u32>,
}

#[derive(Debug, Serialize)]
pub struct ComponentTemplateResponse {
    pub templates: Vec<serde_json::Value>,
    pub total: u32,
    pub categories: Vec<String>,
    pub page: u32,
    pub per_page: u32,
}

#[derive(Debug, Serialize)]
pub struct BuildResult {
    pub build_id: Uuid,
    pub status: BuildStatus,
    pub progress: u8,
    pub logs: Vec<String>,
    pub artifacts: Vec<String>,
    pub metrics: serde_json::Value,
    pub error_message: Option<String>,
}

#[instrument(skip(state))]
pub async fn build_website(
    State(state): State<AppState>,
    Json(request): Json<BuildWebsiteRequest>,
) -> Result<Json<BuildResponse>, (StatusCode, Json<serde_json::Value>)> {
    info!("Received build request for project {}", request.project_id);

    // Validate request
    if let Err(validation_errors) = request.validate() {
        error!("Validation failed: {:?}", validation_errors);
        return Err((
            StatusCode::BAD_REQUEST,
            Json(serde_json::json!({
                "error": "Validation failed",
                "details": validation_errors
            })),
        ));
    }

    // Create build job
    match state
        .builder_service
        .create_build_job(request.build_request)
        .await
    {
        Ok(build_job) => {
            let build_id_str = build_job["job_id"].as_str().unwrap_or("unknown");
            let job_id = uuid::Uuid::parse_str(build_id_str).unwrap_or_else(|_| uuid::Uuid::new_v4());
            info!("Created build job {}", job_id);

            // Start build process asynchronously
            let service = state.builder_service.clone();
            tokio::spawn(async move {
                if let Err(e) = service.process_build_job(job_id).await {
                    error!("Build job {} failed: {}", job_id, e);
                }
            });

            Ok(Json(BuildResponse {
                build_id: job_id,
                status: BuildStatus::Queued,
                message: "Build job created and started".to_string(),
                estimated_time: Some(300), // 5 minutes estimate
            }))
        }
        Err(e) => {
            error!("Failed to create build job: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to create build job",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

#[instrument(skip(state))]
pub async fn build_component(
    State(state): State<AppState>,
    Json(request): Json<ComponentBuildRequest>,
) -> Result<Json<ComponentResponse>, (StatusCode, Json<serde_json::Value>)> {
    info!("Building component: {}", request.component_type);

    match state
        .builder_service
        .build_component(&request.component_type, &request.props)
        .await
    {
        Ok(component) => {
            info!("Successfully built component {}", request.component_type);
            Ok(Json(ComponentResponse {
                component_id: format!("{}_{}", request.component_type, Uuid::new_v4()),
                html: component["html"].as_str().unwrap_or("").to_string(),
                css: component["css"].as_str().unwrap_or("").to_string(),
                js: component["js"].as_str().unwrap_or("").to_string(),
                props: request.props,
            }))
        }
        Err(e) => {
            error!("Failed to build component: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to build component",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

#[instrument(skip(state))]
pub async fn build_static(
    State(state): State<AppState>,
    Json(request): Json<BuildWebsiteRequest>,
) -> Result<Json<BuildResponse>, (StatusCode, Json<serde_json::Value>)> {
    info!("Building static site for project {}", request.project_id);

    match state
        .builder_service
        .build_static_site(request.build_request)
        .await
    {
        Ok(build_result) => {
            let build_id_str = build_result["build_id"].as_str().unwrap_or("unknown");
            let build_id = uuid::Uuid::parse_str(build_id_str).unwrap_or_else(|_| uuid::Uuid::new_v4());
            info!("Successfully built static site: {}", build_id);
            Ok(Json(BuildResponse {
                build_id,
                status: BuildStatus::Success,
                message: "Static site built successfully".to_string(),
                estimated_time: None,
            }))
        }
        Err(e) => {
            error!("Failed to build static site: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to build static site",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

#[instrument(skip(state))]
pub async fn get_build_status(
    State(state): State<AppState>,
    Path(build_id): Path<Uuid>,
    Query(query): Query<BuildStatusQuery>,
) -> Result<Json<BuildResult>, (StatusCode, Json<serde_json::Value>)> {
    info!("Getting build status for {}", build_id);

    match state.builder_service.get_build_status(build_id).await {
        Ok(build_result) => {
            // Create response based on the JsonValue
            let status_str = build_result["status"].as_str().unwrap_or("unknown");
            let status = match status_str {
                "completed" => BuildStatus::Success,
                "failed" => BuildStatus::Failed,
                "running" => BuildStatus::Running,
                _ => BuildStatus::Queued,
            };

            info!("Retrieved build status for {}: {:?}", build_id, status);
            Ok(Json(BuildResult {
                build_id,
                status,
                progress: build_result["progress"].as_i64().unwrap_or(0) as u8,
                logs: if query.include_logs.unwrap_or(false) {
                    vec!["Build started".to_string(), "Build completed".to_string()]
                } else {
                    vec![]
                },
                artifacts: vec![],
                metrics: serde_json::json!({"duration": 120, "size": 1024}),
                error_message: build_result["error"].as_str().map(|s| s.to_string()),
            }))
        }
        Ok(None) => {
            error!("Build {} not found", build_id);
            Err((
                StatusCode::NOT_FOUND,
                Json(serde_json::json!({
                    "error": "Build not found",
                    "build_id": build_id
                })),
            ))
        }
        Err(e) => {
            error!("Failed to get build status: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to get build status",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

#[instrument(skip(state))]
pub async fn get_templates(
    State(state): State<AppState>,
    Query(query): Query<TemplateQuery>,
) -> Result<Json<TemplatesResponse>, (StatusCode, Json<serde_json::Value>)> {
    info!("Getting component templates");

    let page = query.page.unwrap_or(1);
    let limit = query.limit.unwrap_or(20).min(100); // Max 100 items per page

    let tags = query.tags.map(|t| t.split(',').map(|s| s.trim().to_string()).collect());

    match state
        .builder_service
        .get_component_templates(query.category, query.framework, tags, page, limit)
        .await
    {
        Ok((templates, total, categories)) => {
            info!("Retrieved {} templates", templates.len());
            Ok(Json(TemplatesResponse {
                templates,
                total,
                categories,
            }))
        }
        Err(e) => {
            error!("Failed to get templates: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to get templates",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

#[instrument(skip(state))]
pub async fn cancel_build(
    State(state): State<AppState>,
    Path(build_id): Path<Uuid>,
) -> Result<Json<serde_json::Value>, (StatusCode, Json<serde_json::Value>)> {
    info!("Cancelling build {}", build_id);

    match state.builder_service.cancel_build(build_id).await {
        Ok(result) => {
            let success = result["status"].as_str() == Some("cancelled");
            if success {
                info!("Successfully cancelled build {}", build_id);
                Ok(Json(serde_json::json!({
                    "message": "Build cancelled successfully",
                    "build_id": build_id
                })))
            } else {
                error!("Build {} cannot be cancelled", build_id);
                Err((
                    StatusCode::BAD_REQUEST,
                    Json(serde_json::json!({
                        "error": "Build cannot be cancelled",
                        "build_id": build_id
                    })),
                ))
            }
        }
        Err(e) => {
            error!("Failed to cancel build: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to cancel build",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

#[instrument(skip(state))]
pub async fn get_build_logs(
    State(state): State<AppState>,
    Path(build_id): Path<Uuid>,
) -> Result<Json<serde_json::Value>, (StatusCode, Json<serde_json::Value>)> {
    info!("Getting build logs for {}", build_id);

    match state.builder_service.get_build_logs(build_id).await {
        Ok(logs) => {
            let log_array = logs["logs"].as_array().unwrap_or(&vec![]);
            info!("Retrieved {} log entries for build {}", log_array.len(), build_id);
            Ok(Json(serde_json::json!({
                "build_id": build_id,
                "logs": log_array
            })))
        }
        Err(e) => {
            error!("Failed to get build logs: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to get build logs",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

#[instrument(skip(state))]
pub async fn get_build_artifacts(
    State(state): State<AppState>,
    Path(build_id): Path<Uuid>,
) -> Result<Json<serde_json::Value>, (StatusCode, Json<serde_json::Value>)> {
    info!("Getting build artifacts for {}", build_id);

    match state.builder_service.get_build_artifacts(build_id).await {
        Ok(artifacts) => {
            let artifact_array = artifacts["artifacts"].as_array().unwrap_or(&vec![]);
            info!("Retrieved {} artifacts for build {}", artifact_array.len(), build_id);
            Ok(Json(serde_json::json!({
                "build_id": build_id,
                "artifacts": artifact_array
            })))
        }
        Err(e) => {
            error!("Failed to get build artifacts: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to get build artifacts",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

#[instrument(skip(state))]
pub async fn create_preview(
    State(state): State<AppState>,
    Path(build_id): Path<Uuid>,
) -> Result<Json<BuildResponse>, (StatusCode, Json<serde_json::Value>)> {
    match state.builder_service.create_preview(build_id).await {
        Ok(preview) => {
            info!("Created preview for build {}", build_id);
            Ok(Json(BuildResponse {
                build_id,
                status: BuildStatus::Success,
                message: "Preview created successfully".to_string(),
                estimated_time: None,
            }))
        }
        Err(e) => {
            error!("Failed to create preview: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to create preview",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

#[instrument(skip(state))]
pub async fn get_component_templates(
    State(state): State<AppState>,
    Query(query): Query<ComponentTemplateQuery>,
) -> Result<Json<ComponentTemplateResponse>, (StatusCode, Json<serde_json::Value>)> {
    match state
        .builder_service
        .get_component_templates(query.category, query.framework)
        .await
    {
        Ok(templates) => {
            info!("Retrieved {} component templates", templates.len());
            Ok(Json(ComponentTemplateResponse {
                templates,
                total: templates.len() as u32,
                categories: vec!["ui".to_string(), "layout".to_string()],
                page: query.page.unwrap_or(1),
                per_page: query.per_page.unwrap_or(20),
            }))
        }
        Err(e) => {
            error!("Failed to get component templates: {}", e);
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({
                    "error": "Failed to get component templates",
                    "message": e.to_string()
                })),
            ))
        }
    }
}

pub fn router() -> axum::Router<AppState> {
    axum::Router::new()
        .route("/build", post(build_website))
        .route("/build/component", post(build_component))
        .route("/build/static", post(build_static))
        .route("/build/:build_id/status", get(get_build_status))
        .route("/build/:build_id/cancel", post(cancel_build))
        .route("/build/:build_id/logs", get(get_build_logs))
        .route("/build/:build_id/artifacts", get(get_build_artifacts))
        .route("/build/:build_id/preview", post(create_preview))
        .route("/templates", get(get_component_templates))
}
