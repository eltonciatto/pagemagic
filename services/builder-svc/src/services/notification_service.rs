use anyhow::Result;
use std::sync::Arc;
use tracing::{info, instrument};

use crate::config_manager::AppConfig;

#[derive(Clone)]
pub struct NotificationService {
    config: Arc<AppConfig>,
}

impl NotificationService {
    pub async fn new(config: Arc<AppConfig>) -> Result<Self> {
        info!("Initializing notification_service");
        Ok(Self { config })
    }
}
