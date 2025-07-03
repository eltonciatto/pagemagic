use anyhow::Result;
use std::sync::Arc;
use tracing::{info, instrument};

use crate::{config_manager::AppConfig, services::cache_service::CacheService};

#[derive(Clone)]
pub struct CodeGenService {
    config: Arc<AppConfig>,
    cache: Arc<CacheService>,
}

impl CodeGenService {
    pub async fn new(config: Arc<AppConfig>, cache: Arc<CacheService>) -> Result<Self> {
        info!("Initializing code generation service");
        Ok(Self { config, cache })
    }
}
