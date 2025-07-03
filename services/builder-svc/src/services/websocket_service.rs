use anyhow::Result;
use std::sync::Arc;
use tracing::{info, instrument};

use crate::config_manager::AppConfig;

#[derive(Clone)]
pub struct WebSocketService {
    config: Arc<AppConfig>,
}

impl WebSocketService {
    pub async fn new(config: Arc<AppConfig>) -> Result<Self> {
        info!("Initializing websocket_service");
        Ok(Self { config })
    }
}
