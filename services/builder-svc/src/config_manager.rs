use anyhow::{Context, Result};
use config::{Config, ConfigError, Environment, File};
use serde::{Deserialize, Serialize};
use std::path::Path;

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct AppConfig {
    pub server: ServerConfig,
    pub database: DatabaseConfig,
    pub redis: RedisConfig,
    pub observability: ObservabilityConfig,
    pub auth: AuthConfig,
    pub builder: BuilderConfig,
    pub storage: StorageConfig,
    pub external_services: ExternalServicesConfig,
    pub security: SecurityConfig,
    pub features: FeatureFlags,
    pub rate_limiting: RateLimitingConfig,
    pub websocket: WebSocketConfig,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct ServerConfig {
    pub host: String,
    pub port: u16,
    pub workers: usize,
    pub max_connections: usize,
    pub keep_alive: u64,
    pub timeout: u64,
    pub graceful_shutdown_timeout: u64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct DatabaseConfig {
    pub url: String,
    pub max_connections: u32,
    pub min_connections: u32,
    pub acquire_timeout: u64,
    pub idle_timeout: u64,
    pub max_lifetime: u64,
    pub migration_auto: bool,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct RedisConfig {
    pub url: String,
    pub pool_size: u32,
    pub connection_timeout: u64,
    pub command_timeout: u64,
    pub retry_attempts: u32,
    pub retry_delay: u64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct ObservabilityConfig {
    pub metrics: MetricsConfig,
    pub tracing: TracingConfig,
    pub jaeger: JaegerConfig,
    pub logging: LoggingConfig,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct MetricsConfig {
    pub enabled: bool,
    pub port: u16,
    pub path: String,
    pub collection_interval: u64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct TracingConfig {
    pub enabled: bool,
    pub level: String,
    pub format: String,
    pub sample_rate: f64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct JaegerConfig {
    pub enabled: bool,
    pub endpoint: String,
    pub service_name: String,
    pub sample_rate: f64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct LoggingConfig {
    pub level: String,
    pub format: String,
    pub output: String,
    pub file_rotation: bool,
    pub max_file_size: u64,
    pub max_files: u32,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct AuthConfig {
    pub jwt_secret: String,
    pub jwt_expiry: u64,
    pub refresh_token_expiry: u64,
    pub bcrypt_cost: u32,
    pub max_login_attempts: u32,
    pub lockout_duration: u64,
    pub session_timeout: u64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct BuilderConfig {
    pub max_project_size: u64,
    pub max_assets_per_project: u32,
    pub max_concurrent_builds: u32,
    pub build_timeout: u64,
    pub cache_enabled: bool,
    pub cache_ttl: u64,
    pub optimization_enabled: bool,
    pub real_time_enabled: bool,
    pub collaboration_enabled: bool,
    pub version_control_enabled: bool,
    pub ai_assistance_enabled: bool,
    pub theme_customization_enabled: bool,
    pub responsive_design_enabled: bool,
    pub accessibility_checks_enabled: bool,
    pub performance_optimization_enabled: bool,
    pub seo_optimization_enabled: bool,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct StorageConfig {
    pub provider: String, // "local", "s3", "gcs", "azure"
    pub base_path: String,
    pub max_file_size: u64,
    pub allowed_extensions: Vec<String>,
    pub compression_enabled: bool,
    pub cdn_enabled: bool,
    pub cdn_url: String,
    pub backup_enabled: bool,
    pub backup_retention: u32,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct ExternalServicesConfig {
    pub openai: OpenAIConfig,
    pub stripe: StripeConfig,
    pub sendgrid: SendGridConfig,
    pub github: GitHubConfig,
    pub docker: DockerConfig,
    pub kubernetes: KubernetesConfig,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct OpenAIConfig {
    pub api_key: String,
    pub model: String,
    pub max_tokens: u32,
    pub temperature: f32,
    pub timeout: u64,
    pub rate_limit: u32,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct StripeConfig {
    pub secret_key: String,
    pub publishable_key: String,
    pub webhook_secret: String,
    pub api_version: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct SendGridConfig {
    pub api_key: String,
    pub from_email: String,
    pub from_name: String,
    pub template_id: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct GitHubConfig {
    pub token: String,
    pub webhook_secret: String,
    pub app_id: String,
    pub private_key: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct DockerConfig {
    pub registry_url: String,
    pub username: String,
    pub password: String,
    pub build_timeout: u64,
    pub max_build_memory: u64,
    pub max_build_cpu: f32,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct KubernetesConfig {
    pub config_path: String,
    pub namespace: String,
    pub cluster_name: String,
    pub registry_secret: String,
    pub resource_limits: ResourceLimits,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct ResourceLimits {
    pub cpu: String,
    pub memory: String,
    pub storage: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct SecurityConfig {
    pub cors_origins: Vec<String>,
    pub csrf_protection: bool,
    pub security_headers: bool,
    pub content_security_policy: String,
    pub rate_limiting_enabled: bool,
    pub ip_whitelist: Vec<String>,
    pub encryption_key: String,
    pub audit_logging: bool,
    pub vulnerability_scanning: bool,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct FeatureFlags {
    pub ai_generation: bool,
    pub real_time_collaboration: bool,
    pub version_control: bool,
    pub analytics: bool,
    pub a_b_testing: bool,
    pub multi_language: bool,
    pub themes: bool,
    pub plugins: bool,
    pub marketplace: bool,
    pub white_label: bool,
    pub enterprise_features: bool,
    pub mobile_app: bool,
    pub offline_mode: bool,
    pub pwa_support: bool,
    pub advanced_seo: bool,
    pub e_commerce: bool,
    pub cms_integration: bool,
    pub third_party_integrations: bool,
    pub custom_domains: bool,
    pub ssl_certificates: bool,
    pub backup_restore: bool,
    pub monitoring_dashboard: bool,
    pub user_management: bool,
    pub team_collaboration: bool,
    pub project_templates: bool,
    pub asset_library: bool,
    pub form_builder: bool,
    pub database_integration: bool,
    pub api_builder: bool,
    pub serverless_functions: bool,
    pub headless_cms: bool,
    pub multi_site_management: bool,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct RateLimitingConfig {
    pub enabled: bool,
    pub requests_per_minute: u32,
    pub burst_size: u32,
    pub window_size: u64,
    pub cleanup_interval: u64,
    pub redis_enabled: bool,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct WebSocketConfig {
    pub enabled: bool,
    pub max_connections: u32,
    pub heartbeat_interval: u64,
    pub max_message_size: u32,
    pub buffer_size: u32,
    pub compression_enabled: bool,
}

impl AppConfig {
    pub fn new(config_path: &str, environment: &str) -> Result<Self> {
        let mut builder = Config::builder()
            // Start with default configuration
            .add_source(File::with_name("configs/default").required(false))
            // Add environment-specific configuration
            .add_source(File::with_name(&format!("configs/{}", environment)).required(false))
            // Add local configuration file if provided
            .add_source(File::with_name(config_path).required(false))
            // Add environment variables with prefix
            .add_source(Environment::with_prefix("PAGEMAGIC").separator("__"));

        let config = builder
            .build()
            .context("Failed to build configuration")?;

        let app_config: AppConfig = config
            .try_deserialize()
            .context("Failed to deserialize configuration")?;

        // Validate configuration
        app_config.validate()?;

        Ok(app_config)
    }

    pub fn validate(&self) -> Result<()> {
        // Validate server configuration
        if self.server.port == 0 || self.server.port > 65535 {
            anyhow::bail!("Invalid server port: {}", self.server.port);
        }

        if self.server.workers == 0 {
            anyhow::bail!("Server workers must be greater than 0");
        }

        // Validate database configuration
        if self.database.url.is_empty() {
            anyhow::bail!("Database URL cannot be empty");
        }

        if self.database.max_connections == 0 {
            anyhow::bail!("Database max_connections must be greater than 0");
        }

        // Validate Redis configuration
        if self.redis.url.is_empty() {
            anyhow::bail!("Redis URL cannot be empty");
        }

        // Validate authentication configuration
        if self.auth.jwt_secret.len() < 32 {
            anyhow::bail!("JWT secret must be at least 32 characters long");
        }

        if self.auth.bcrypt_cost < 4 || self.auth.bcrypt_cost > 31 {
            anyhow::bail!("BCrypt cost must be between 4 and 31");
        }

        // Validate builder configuration
        if self.builder.max_project_size == 0 {
            anyhow::bail!("Max project size must be greater than 0");
        }

        if self.builder.max_concurrent_builds == 0 {
            anyhow::bail!("Max concurrent builds must be greater than 0");
        }

        // Validate storage configuration
        let valid_providers = ["local", "s3", "gcs", "azure"];
        if !valid_providers.contains(&self.storage.provider.as_str()) {
            anyhow::bail!("Invalid storage provider: {}", self.storage.provider);
        }

        // Validate external services if features are enabled
        if self.features.ai_generation && self.external_services.openai.api_key.is_empty() {
            anyhow::bail!("OpenAI API key is required when AI generation is enabled");
        }

        // Validate observability configuration
        if self.observability.metrics.enabled && self.observability.metrics.port == 0 {
            anyhow::bail!("Metrics port must be specified when metrics are enabled");
        }

        Ok(())
    }

    pub fn is_production(&self) -> bool {
        std::env::var("ENVIRONMENT")
            .unwrap_or_else(|_| "development".to_string())
            .to_lowercase() == "production"
    }

    pub fn is_development(&self) -> bool {
        !self.is_production()
    }

    pub fn get_environment(&self) -> String {
        std::env::var("ENVIRONMENT")
            .unwrap_or_else(|_| "development".to_string())
    }
}

impl Default for AppConfig {
    fn default() -> Self {
        Self {
            server: ServerConfig {
                host: "0.0.0.0".to_string(),
                port: 8080,
                workers: num_cpus::get(),
                max_connections: 1000,
                keep_alive: 60,
                timeout: 30,
                graceful_shutdown_timeout: 30,
            },
            database: DatabaseConfig {
                url: "postgresql://pagemagic:pagemagic@localhost:5432/pagemagic".to_string(),
                max_connections: 10,
                min_connections: 1,
                acquire_timeout: 30,
                idle_timeout: 600,
                max_lifetime: 3600,
                migration_auto: true,
            },
            redis: RedisConfig {
                url: "redis://localhost:6379".to_string(),
                pool_size: 10,
                connection_timeout: 5,
                command_timeout: 5,
                retry_attempts: 3,
                retry_delay: 1000,
            },
            observability: ObservabilityConfig {
                metrics: MetricsConfig {
                    enabled: true,
                    port: 9090,
                    path: "/metrics".to_string(),
                    collection_interval: 15,
                },
                tracing: TracingConfig {
                    enabled: true,
                    level: "info".to_string(),
                    format: "json".to_string(),
                    sample_rate: 0.1,
                },
                jaeger: JaegerConfig {
                    enabled: false,
                    endpoint: "http://localhost:14268/api/traces".to_string(),
                    service_name: "builder-svc".to_string(),
                    sample_rate: 0.1,
                },
                logging: LoggingConfig {
                    level: "info".to_string(),
                    format: "json".to_string(),
                    output: "stdout".to_string(),
                    file_rotation: true,
                    max_file_size: 100 * 1024 * 1024, // 100MB
                    max_files: 10,
                },
            },
            auth: AuthConfig {
                jwt_secret: "super_secret_jwt_key_change_in_production_32_chars_minimum".to_string(),
                jwt_expiry: 3600, // 1 hour
                refresh_token_expiry: 86400 * 30, // 30 days
                bcrypt_cost: 12,
                max_login_attempts: 5,
                lockout_duration: 900, // 15 minutes
                session_timeout: 3600, // 1 hour
            },
            builder: BuilderConfig {
                max_project_size: 100 * 1024 * 1024, // 100MB
                max_assets_per_project: 1000,
                max_concurrent_builds: 10,
                build_timeout: 300, // 5 minutes
                cache_enabled: true,
                cache_ttl: 3600, // 1 hour
                optimization_enabled: true,
                real_time_enabled: true,
                collaboration_enabled: true,
                version_control_enabled: true,
                ai_assistance_enabled: true,
                theme_customization_enabled: true,
                responsive_design_enabled: true,
                accessibility_checks_enabled: true,
                performance_optimization_enabled: true,
                seo_optimization_enabled: true,
            },
            storage: StorageConfig {
                provider: "local".to_string(),
                base_path: "./storage".to_string(),
                max_file_size: 10 * 1024 * 1024, // 10MB
                allowed_extensions: vec![
                    "jpg".to_string(), "jpeg".to_string(), "png".to_string(), "gif".to_string(),
                    "svg".to_string(), "webp".to_string(), "ico".to_string(),
                    "css".to_string(), "js".to_string(), "json".to_string(), "xml".to_string(),
                    "woff".to_string(), "woff2".to_string(), "ttf".to_string(), "eot".to_string(),
                ],
                compression_enabled: true,
                cdn_enabled: false,
                cdn_url: "".to_string(),
                backup_enabled: true,
                backup_retention: 30, // 30 days
            },
            external_services: ExternalServicesConfig {
                openai: OpenAIConfig {
                    api_key: "".to_string(),
                    model: "gpt-4".to_string(),
                    max_tokens: 4000,
                    temperature: 0.7,
                    timeout: 30,
                    rate_limit: 100,
                },
                stripe: StripeConfig {
                    secret_key: "".to_string(),
                    publishable_key: "".to_string(),
                    webhook_secret: "".to_string(),
                    api_version: "2023-10-16".to_string(),
                },
                sendgrid: SendGridConfig {
                    api_key: "".to_string(),
                    from_email: "noreply@pagemagic.com".to_string(),
                    from_name: "Page Magic".to_string(),
                    template_id: "".to_string(),
                },
                github: GitHubConfig {
                    token: "".to_string(),
                    webhook_secret: "".to_string(),
                    app_id: "".to_string(),
                    private_key: "".to_string(),
                },
                docker: DockerConfig {
                    registry_url: "docker.io".to_string(),
                    username: "".to_string(),
                    password: "".to_string(),
                    build_timeout: 600, // 10 minutes
                    max_build_memory: 2 * 1024 * 1024 * 1024, // 2GB
                    max_build_cpu: 2.0,
                },
                kubernetes: KubernetesConfig {
                    config_path: "~/.kube/config".to_string(),
                    namespace: "pagemagic".to_string(),
                    cluster_name: "pagemagic-cluster".to_string(),
                    registry_secret: "registry-secret".to_string(),
                    resource_limits: ResourceLimits {
                        cpu: "1000m".to_string(),
                        memory: "1Gi".to_string(),
                        storage: "10Gi".to_string(),
                    },
                },
            },
            security: SecurityConfig {
                cors_origins: vec!["*".to_string()],
                csrf_protection: true,
                security_headers: true,
                content_security_policy: "default-src 'self'".to_string(),
                rate_limiting_enabled: true,
                ip_whitelist: vec![],
                encryption_key: "super_secret_encryption_key_32_chars".to_string(),
                audit_logging: true,
                vulnerability_scanning: false,
            },
            features: FeatureFlags {
                ai_generation: true,
                real_time_collaboration: true,
                version_control: true,
                analytics: true,
                a_b_testing: true,
                multi_language: true,
                themes: true,
                plugins: true,
                marketplace: true,
                white_label: true,
                enterprise_features: true,
                mobile_app: true,
                offline_mode: true,
                pwa_support: true,
                advanced_seo: true,
                e_commerce: true,
                cms_integration: true,
                third_party_integrations: true,
                custom_domains: true,
                ssl_certificates: true,
                backup_restore: true,
                monitoring_dashboard: true,
                user_management: true,
                team_collaboration: true,
                project_templates: true,
                asset_library: true,
                form_builder: true,
                database_integration: true,
                api_builder: true,
                serverless_functions: true,
                headless_cms: true,
                multi_site_management: true,
            },
            rate_limiting: RateLimitingConfig {
                enabled: true,
                requests_per_minute: 100,
                burst_size: 10,
                window_size: 60,
                cleanup_interval: 300,
                redis_enabled: true,
            },
            websocket: WebSocketConfig {
                enabled: true,
                max_connections: 1000,
                heartbeat_interval: 30,
                max_message_size: 1024 * 1024, // 1MB
                buffer_size: 1024,
                compression_enabled: true,
            },
        }
    }
}
