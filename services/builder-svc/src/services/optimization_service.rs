use anyhow::Result;
use std::sync::Arc;
use tracing::{info, instrument};

use crate::config_manager::AppConfig;

#[derive(Clone)]
pub struct OptimizationService {
    config: Arc<AppConfig>,
}

impl OptimizationService {
    pub async fn new(config: Arc<AppConfig>) -> Result<Self> {
        info!("Initializing optimization_service");
        Ok(Self { config })
    }
}
