use anyhow::Result;
use std::sync::Arc;
use tracing::{info, instrument};

use crate::{config_manager::AppConfig, services::{cache_service::CacheService, code_gen_service::CodeGenService}};

#[derive(Clone)]
pub struct TemplateService {
    config: Arc<AppConfig>,
    db: sqlx::PgPool,
    cache: Arc<CacheService>,
    code_gen: Arc<CodeGenService>,
}

impl TemplateService {
    pub async fn new(
        config: Arc<AppConfig>,
        db: sqlx::PgPool,
        cache: Arc<CacheService>,
        code_gen: Arc<CodeGenService>,
    ) -> Result<Self> {
        info!("Initializing template service");
        Ok(Self { config, db, cache, code_gen })
    }
}
