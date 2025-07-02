use serde::{Deserialize, Serialize};
use config::{Config as ConfigBuilder, ConfigError, Environment, File};
use std::env;

#[derive(Debug, Deserialize, Clone)]
pub struct Config {
    pub server: ServerConfig,
    pub database: DatabaseConfig,
    pub redis: RedisConfig,
    pub services: ServicesConfig,
    pub build: BuildConfig,
}

#[derive(Debug, Deserialize, Clone)]
pub struct ServerConfig {
    pub host: String,
    pub port: u16,
    pub environment: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct DatabaseConfig {
    pub url: String,
    pub max_connections: u32,
}

#[derive(Debug, Deserialize, Clone)]
pub struct RedisConfig {
    pub url: String,
    pub pool_size: u32,
}

#[derive(Debug, Deserialize, Clone)]
pub struct ServicesConfig {
    pub auth_svc_url: String,
    pub prompt_svc_url: String,
    pub build_svc_url: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct BuildConfig {
    pub output_dir: String,
    pub template_dir: String,
    pub assets_dir: String,
    pub max_concurrent_builds: usize,
    pub build_timeout_seconds: u64,
}

impl Config {
    pub fn load() -> Result<Self, ConfigError> {
        let run_mode = env::var("RUN_MODE").unwrap_or_else(|_| "development".into());

        let s = ConfigBuilder::builder()
            // Arquivo de configuração padrão
            .add_source(File::with_name("config/default").required(false))
            // Arquivo específico do ambiente
            .add_source(File::with_name(&format!("config/{}", run_mode)).required(false))
            // Variáveis de ambiente (com prefixo BUILDER_)
            .add_source(Environment::with_prefix("BUILDER").separator("_"))
            // Configurações padrão
            .set_default("server.host", "0.0.0.0")?
            .set_default("server.port", 8082)?
            .set_default("server.environment", "development")?
            .set_default("database.url", "postgres://pagemagic:password@localhost:5432/pagemagic")?
            .set_default("database.max_connections", 10)?
            .set_default("redis.url", "redis://localhost:6379")?
            .set_default("redis.pool_size", 10)?
            .set_default("services.auth_svc_url", "http://localhost:8080")?
            .set_default("services.prompt_svc_url", "http://localhost:3001")?
            .set_default("services.build_svc_url", "http://localhost:8083")?
            .set_default("build.output_dir", "/tmp/pagemagic/builds")?
            .set_default("build.template_dir", "/app/templates")?
            .set_default("build.assets_dir", "/app/assets")?
            .set_default("build.max_concurrent_builds", 5)?
            .set_default("build.build_timeout_seconds", 300)?
            .build()?;

        s.try_deserialize()
    }
}
