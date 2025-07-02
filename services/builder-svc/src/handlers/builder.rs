use axum::{
    extract::{Path, Query, State},
    http::StatusCode,
    response::Json,
};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tracing::{error, info, instrument};
use uuid::Uuid;
use validator::Validate;

use crate::{
    models::{
        BuildRequest, BuildResult, BuildStatus, BuildType, ComponentTemplate,
        BuildJob, WebsiteOutput
    },
    AppState,
};

#[derive(Debug, Deserialize, Validate)]
pub struct BuildWebsiteRequest {
    #[validate(length(min = 1))]
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
            info!("Created build job {}", build_job.id);

            // Start build process asynchronously
            let service = state.builder_service.clone();
            let job_id = build_job.id;
            tokio::spawn(async move {
                if let Err(e) = service.process_build_job(job_id).await {
                    error!("Build job {} failed: {}", job_id, e);
                }
            });

            Ok(Json(BuildResponse {
                build_id: build_job.id,
                status: build_job.status,
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
        .build_component(&request.component_type, &request.props, &request.framework)
        .await
    {
        Ok(component) => {
            info!("Successfully built component {}", request.component_type);
            Ok(Json(ComponentResponse {
                component_id: format!("{}_{}", request.component_type, Uuid::new_v4()),
                html: component.html,
                css: component.css,
                js: component.js,
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
            info!("Successfully built static site: {}", build_result.build_id);
            Ok(Json(BuildResponse {
                build_id: build_result.build_id,
                status: build_result.status,
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
        Ok(Some(mut build_result)) => {
            // Filter response based on query parameters
            if !query.include_logs.unwrap_or(false) {
                build_result.logs.clear();
            }

            if !query.include_metrics.unwrap_or(true) {
                // Keep metrics by default, only remove if explicitly requested
            }

            info!("Retrieved build status for {}: {:?}", build_id, build_result.status);
            Ok(Json(build_result))
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
        Ok(true) => {
            info!("Successfully cancelled build {}", build_id);
            Ok(Json(serde_json::json!({
                "message": "Build cancelled successfully",
                "build_id": build_id
            })))
        }
        Ok(false) => {
            error!("Build {} cannot be cancelled", build_id);
            Err((
                StatusCode::BAD_REQUEST,
                Json(serde_json::json!({
                    "error": "Build cannot be cancelled",
                    "build_id": build_id
                })),
            ))
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
            info!("Retrieved {} log entries for build {}", logs.len(), build_id);
            Ok(Json(serde_json::json!({
                "build_id": build_id,
                "logs": logs
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
            info!("Retrieved {} artifacts for build {}", artifacts.len(), build_id);
            Ok(Json(serde_json::json!({
                "build_id": build_id,
                "artifacts": artifacts
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
pub async fn preview_build(
    State(state): State<AppState>,
    Path(build_id): Path<Uuid>,
) -> Result<Json<serde_json::Value>, (StatusCode, Json<serde_json::Value>)> {
    info!("Creating preview for build {}", build_id);

    match state.builder_service.create_preview(build_id).await {
        Ok(preview_url) => {
            info!("Created preview for build {}: {}", build_id, preview_url);
            Ok(Json(serde_json::json!({
                "build_id": build_id,
                "preview_url": preview_url
            })))
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
