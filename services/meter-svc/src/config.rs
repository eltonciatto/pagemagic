use serde::Deserialize;

#[derive(Debug, Clone, Deserialize)]
pub struct Config {
    pub server: ServerConfig,
    pub database: DatabaseConfig,
    pub stripe: StripeConfig,
    pub nats: NatsConfig,
    pub metrics: MetricsConfig,
}

#[derive(Debug, Clone, Deserialize)]
pub struct ServerConfig {
    pub host: String,
    pub port: u16,
}

#[derive(Debug, Clone, Deserialize)]
pub struct DatabaseConfig {
    pub url: String,
    pub max_connections: u32,
}

#[derive(Debug, Clone, Deserialize)]
pub struct StripeConfig {
    pub secret_key: String,
    pub webhook_secret: String,
    pub api_url: String,
}

#[derive(Debug, Clone, Deserialize)]
pub struct NatsConfig {
    pub url: String,
    pub subjects: NatsSubjects,
}

#[derive(Debug, Clone, Deserialize)]
pub struct NatsSubjects {
    pub usage_events: String,
    pub meter_updates: String,
}

#[derive(Debug, Clone, Deserialize)]
pub struct MetricsConfig {
    pub enabled: bool,
    pub bind_address: String,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            server: ServerConfig {
                host: "0.0.0.0".to_string(),
                port: 8080,
            },
            database: DatabaseConfig {
                url: "postgresql://user:pass@localhost:5432/pagemagic".to_string(),
                max_connections: 10,
            },
            stripe: StripeConfig {
                secret_key: "sk_test_...".to_string(),
                webhook_secret: "whsec_...".to_string(),
                api_url: "https://api.stripe.com".to_string(),
            },
            nats: NatsConfig {
                url: "nats://localhost:4222".to_string(),
                subjects: NatsSubjects {
                    usage_events: "pagemagic.usage.events".to_string(),
                    meter_updates: "pagemagic.meter.updates".to_string(),
                },
            },
            metrics: MetricsConfig {
                enabled: true,
                bind_address: "0.0.0.0:9090".to_string(),
            },
        }
    }
}

pub fn load_config() -> anyhow::Result<Config> {
    let mut settings = config::Config::default();
    
    // Load from environment variables
    settings.merge(config::Environment::with_prefix("METER"))?;
    
    // Load from config file if it exists
    if let Ok(config_file) = std::env::var("CONFIG_FILE") {
        settings.merge(config::File::with_name(&config_file))?;
    }
    
    Ok(settings.try_into()?)
}
