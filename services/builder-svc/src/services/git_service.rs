use anyhow::Result;
use std::sync::Arc;
use tracing::{info, instrument};

use crate::config_manager::AppConfig;

#[derive(Clone)]
pub struct GitService {
    config: Arc<AppConfig>,
}

impl GitService {
    pub async fn new(config: Arc<AppConfig>) -> Result<Self> {
        info!("Initializing git_service");
        Ok(Self { config })
    }
}
