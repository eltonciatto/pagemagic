use axum::{
    extract::{Request, State},
    http::StatusCode,
    middleware::Next,
    response::Response,
};
use dashmap::DashMap;
use std::{
    net::IpAddr,
    sync::Arc,
    time::{Duration, Instant},
};
use tokio::time::sleep;
use tracing::{info, warn, instrument};

use crate::AppState;

// Rate limiting middleware
pub async fn rate_limit_middleware(
    State(state): State<AppState>,
    request: Request,
    next: Next,
) -> Result<Response, StatusCode> {
    if !state.config.rate_limiting.enabled {
        return Ok(next.run(request).await);
    }

    // Extract client IP
    let client_ip = extract_client_ip(&request);
    
    // Check rate limit
    if is_rate_limited(&state, &client_ip).await {
        warn!("Rate limit exceeded for IP: {}", client_ip);
        return Err(StatusCode::TOO_MANY_REQUESTS);
    }

    Ok(next.run(request).await)
}

// Metrics middleware
pub async fn metrics_middleware(
    State(_state): State<AppState>,
    request: Request,
    next: Next,
) -> Response {
    let start = Instant::now();
    let method = request.method().clone();
    let path = request.uri().path().to_string();
    
    let response = next.run(request).await;
    
    let duration = start.elapsed();
    let status = response.status();
    
    // Record metrics (simplified - would use proper metrics crate in production)
    info!(
        method = %method,
        path = %path,
        status = %status,
        duration_ms = %duration.as_millis(),
        "HTTP request processed"
    );
    
    response
}

// Authentication middleware
#[instrument(skip(state, request, next))]
pub async fn auth_middleware(
    State(state): State<AppState>,
    mut request: Request,
    next: Next,
) -> Result<Response, StatusCode> {
    // Extract Authorization header
    let auth_header = request
        .headers()
        .get("Authorization")
        .and_then(|header| header.to_str().ok());

    if let Some(auth_header) = auth_header {
        if let Some(token) = auth_header.strip_prefix("Bearer ") {
            // Validate token
            match state.auth_service.validate_token(token).await {
                Ok(user) => {
                    // Add user to request extensions
                    request.extensions_mut().insert(user);
                    return Ok(next.run(request).await);
                }
                Err(e) => {
                    warn!("Token validation failed: {}", e);
                    return Err(StatusCode::UNAUTHORIZED);
                }
            }
        }
    }

    Err(StatusCode::UNAUTHORIZED)
}

// Helper functions
fn extract_client_ip(request: &Request) -> IpAddr {
    // Try to get real IP from headers (reverse proxy)
    if let Some(forwarded_for) = request.headers().get("X-Forwarded-For") {
        if let Ok(forwarded_str) = forwarded_for.to_str() {
            if let Some(first_ip) = forwarded_str.split(',').next() {
                if let Ok(ip) = first_ip.trim().parse() {
                    return ip;
                }
            }
        }
    }

    if let Some(real_ip) = request.headers().get("X-Real-IP") {
        if let Ok(ip_str) = real_ip.to_str() {
            if let Ok(ip) = ip_str.parse() {
                return ip;
            }
        }
    }

    // Fallback to connection IP (not available in this context, using localhost)
    IpAddr::V4(std::net::Ipv4Addr::new(127, 0, 0, 1))
}

// Simple in-memory rate limiter (would use Redis in production)
static RATE_LIMITER: once_cell::sync::Lazy<Arc<DashMap<IpAddr, RateLimitInfo>>> =
    once_cell::sync::Lazy::new(|| Arc::new(DashMap::new()));

#[derive(Debug, Clone)]
struct RateLimitInfo {
    count: u32,
    window_start: Instant,
}

async fn is_rate_limited(state: &AppState, ip: &IpAddr) -> bool {
    let config = &state.config.rate_limiting;
    let window_duration = Duration::from_secs(config.window_size);
    let now = Instant::now();

    let mut entry = RATE_LIMITER.entry(*ip).or_insert(RateLimitInfo {
        count: 0,
        window_start: now,
    });

    // Reset window if expired
    if now.duration_since(entry.window_start) >= window_duration {
        entry.count = 0;
        entry.window_start = now;
    }

    entry.count += 1;

    // Check if limit exceeded
    if entry.count > config.requests_per_minute {
        return true;
    }

    // Clean up old entries periodically
    if RATE_LIMITER.len() > 10000 {
        tokio::spawn(cleanup_rate_limiter(window_duration));
    }

    false
}

async fn cleanup_rate_limiter(window_duration: Duration) {
    sleep(Duration::from_secs(60)).await; // Wait a bit before cleanup
    
    let now = Instant::now();
    RATE_LIMITER.retain(|_, info| {
        now.duration_since(info.window_start) < window_duration * 2
    });
}
