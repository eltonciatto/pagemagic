use anyhow::Result;
use std::sync::Arc;
use tracing::{info, instrument};

use crate::{config_manager::AppConfig, services::cache_service::CacheService};

#[derive(Clone)]
pub struct ThemeService {
    config: Arc<AppConfig>,
    db: sqlx::PgPool,
    cache: Arc<CacheService>,
}

impl ThemeService {
    pub async fn new(
        config: Arc<AppConfig>,
        db: sqlx::PgPool,
        cache: Arc<CacheService>,
    ) -> Result<Self> {
        info!("Initializing theme service");
        Ok(Self { config, db, cache })
    }
}
