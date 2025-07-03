use axum::{
    extract::State,
    http::StatusCode,
    response::Json,
    routing::{get, post},
    Router,
};
use serde::{Deserialize, Serialize};
use tracing::{error, info, instrument};

use crate::AppState;

#[derive(Debug, Deserialize)]
pub struct LoginRequest {
    pub email: String,
    pub password: String,
}

#[derive(Debug, Deserialize)]
pub struct RefreshTokenRequest {
    pub refresh_token: String,
}

#[derive(Debug, Deserialize)]
pub struct ValidateTokenRequest {
    pub token: String,
}

#[derive(Debug, Serialize)]
pub struct AuthResponse {
    pub access_token: String,
    pub refresh_token: String,
    pub expires_in: u64,
    pub user: UserInfo,
}

#[derive(Debug, Clone, Serialize)]
pub struct UserInfo {
    pub id: uuid::Uuid,
    pub email: String,
    pub username: String,
    pub role: String,
}

#[derive(Debug, Serialize)]
pub struct ValidateResponse {
    pub valid: bool,
    pub user: Option<UserInfo>,
}

pub fn router() -> Router<AppState> {
    Router::new()
        .route("/login", post(login))
        .route("/refresh", post(refresh_token))
        .route("/validate", post(validate_token))
}

#[instrument(skip(state, request))]
pub async fn login(
    State(state): State<AppState>,
    Json(request): Json<LoginRequest>,
) -> Result<Json<AuthResponse>, StatusCode> {
    info!("Login attempt for email: {}", request.email);
    
    match state.auth_service.authenticate(&request.email, &request.password).await {
        Ok(auth_response) => {
            info!("Login successful for user: {}", auth_response.user.id);
            Ok(Json(auth_response))
        }
        Err(e) => {
            error!("Login failed: {}", e);
            Err(StatusCode::UNAUTHORIZED)
        }
    }
}

#[instrument(skip(state, request))]
pub async fn refresh_token(
    State(state): State<AppState>,
    Json(request): Json<RefreshTokenRequest>,
) -> Result<Json<AuthResponse>, StatusCode> {
    info!("Token refresh attempt");
    
    match state.auth_service.refresh_token(&request.refresh_token).await {
        Ok(auth_response) => {
            info!("Token refresh successful");
            Ok(Json(auth_response))
        }
        Err(e) => {
            error!("Token refresh failed: {}", e);
            Err(StatusCode::UNAUTHORIZED)
        }
    }
}

#[instrument(skip(state, request))]
pub async fn validate_token(
    State(state): State<AppState>,
    Json(request): Json<ValidateTokenRequest>,
) -> Result<Json<ValidateResponse>, StatusCode> {
    match state.auth_service.validate_token(&request.token).await {
        Ok(user) => Ok(Json(ValidateResponse {
            valid: true,
            user: Some(user),
        })),
        Err(_) => Ok(Json(ValidateResponse {
            valid: false,
            user: None,
        })),
    }
}
