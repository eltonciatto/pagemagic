use crate::models::*;
use crate::config::Config;
use tokio_postgres::{Client, NoTls};
use reqwest::Client as HttpClient;
use anyhow::{Result, anyhow};
use tracing::{info, error, warn};
use std::collections::HashMap;
use chrono::{DateTime, Utc, Duration};
use uuid::Uuid;

pub struct MeterService {
    db_client: Client,
    http_client: HttpClient,
    config: Config,
    meter_configs: Vec<MeterConfig>,
}

impl MeterService {
    pub async fn new(config: Config) -> Result<Self> {
        // Connect to database
        let (client, connection) = tokio_postgres::connect(&config.database.url, NoTls).await?;
        
        // Spawn connection task
        tokio::spawn(async move {
            if let Err(e) = connection.await {
                error!("Database connection error: {}", e);
            }
        });
        
        // Initialize HTTP client
        let http_client = HttpClient::new();
        
        // Load meter configurations
        let meter_configs = Self::load_meter_configs().await?;
        
        Ok(Self {
            db_client: client,
            http_client,
            config,
            meter_configs,
        })
    }
    
    async fn load_meter_configs() -> Result<Vec<MeterConfig>> {
        // In production, load from database or config file
        Ok(vec![
            MeterConfig {
                name: "page_generate".to_string(),
                stripe_meter_id: "mtr_page_generate".to_string(),
                aggregation_type: AggregationType::Count,
                event_filter: EventFilter {
                    event_types: vec!["page_generate".to_string()],
                    metadata_filters: HashMap::new(),
                },
            },
            MeterConfig {
                name: "ai_token".to_string(),
                stripe_meter_id: "mtr_ai_token".to_string(),
                aggregation_type: AggregationType::Sum,
                event_filter: EventFilter {
                    event_types: vec!["ai_token_usage".to_string()],
                    metadata_filters: HashMap::new(),
                },
            },
            MeterConfig {
                name: "container_hours".to_string(),
                stripe_meter_id: "mtr_container_hours".to_string(),
                aggregation_type: AggregationType::Sum,
                event_filter: EventFilter {
                    event_types: vec!["container_time".to_string()],
                    metadata_filters: HashMap::new(),
                },
            },
            MeterConfig {
                name: "storage_gb".to_string(),
                stripe_meter_id: "mtr_storage_gb".to_string(),
                aggregation_type: AggregationType::Max,
                event_filter: EventFilter {
                    event_types: vec!["storage_usage".to_string()],
                    metadata_filters: HashMap::new(),
                },
            },
        ])
    }
    
    pub async fn process_event(&self, event: &UsageEvent) -> Result<()> {
        // Store event in database
        let event_id = Uuid::new_v4();
        let query = r#"
            INSERT INTO usage_events (id, event_type, user_id, project_id, site_id, timestamp, metadata)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        "#;
        
        self.db_client.execute(
            query,
            &[
                &event_id,
                &event.event_type,
                &event.user_id,
                &event.project_id,
                &event.site_id,
                &event.timestamp,
                &serde_json::to_value(&event.metadata)?,
            ],
        ).await?;
        
        // Process for each matching meter
        for meter_config in &self.meter_configs {
            if self.event_matches_filter(event, &meter_config.event_filter) {
                self.update_meter_aggregation(meter_config, event).await?;
            }
        }
        
        Ok(())
    }
    
    pub async fn process_batch(&self, batch: &EventBatch) -> Result<usize> {
        let mut processed = 0;
        
        for event in &batch.events {
            match self.process_event(event).await {
                Ok(_) => processed += 1,
                Err(e) => {
                    warn!("Failed to process event in batch: {}", e);
                }
            }
        }
        
        Ok(processed)
    }
    
    fn event_matches_filter(&self, event: &UsageEvent, filter: &EventFilter) -> bool {
        // Check event type
        if !filter.event_types.contains(&event.event_type) {
            return false;
        }
        
        // Check metadata filters
        for (key, value) in &filter.metadata_filters {
            if let Some(event_value) = event.metadata.get(key) {
                if event_value.as_str() != Some(value) {
                    return false;
                }
            } else {
                return false;
            }
        }
        
        true
    }
    
    async fn update_meter_aggregation(&self, meter_config: &MeterConfig, event: &UsageEvent) -> Result<()> {
        let value = match meter_config.aggregation_type {
            AggregationType::Count => 1.0,
            AggregationType::Sum => {
                // Extract numeric value from metadata
                event.metadata.get("value")
                    .and_then(|v| v.as_f64())
                    .unwrap_or(1.0)
            },
            _ => 1.0, // Simplified for other types
        };
        
        // Update aggregation table
        let query = r#"
            INSERT INTO meter_aggregations (meter_name, user_id, period_start, period_end, value, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6)
            ON CONFLICT (meter_name, user_id, period_start, period_end)
            DO UPDATE SET 
                value = meter_aggregations.value + $5,
                updated_at = $6
        "#;
        
        let period_start = self.get_period_start(event.timestamp);
        let period_end = period_start + Duration::hours(1);
        
        self.db_client.execute(
            query,
            &[
                &meter_config.name,
                &event.user_id,
                &period_start,
                &period_end,
                &value,
                &Utc::now(),
            ],
        ).await?;
        
        Ok(())
    }
    
    fn get_period_start(&self, timestamp: DateTime<Utc>) -> DateTime<Utc> {
        // Round down to the hour
        timestamp.date_naive()
            .and_hms_opt(timestamp.hour(), 0, 0)
            .unwrap()
            .and_utc()
    }
    
    pub async fn sync_to_stripe(&self) -> Result<usize> {
        // Get pending aggregations to sync
        let query = r#"
            SELECT meter_name, user_id, value, period_start, period_end
            FROM meter_aggregations
            WHERE synced_to_stripe = false
            AND period_end < NOW() - INTERVAL '1 hour'
        "#;
        
        let rows = self.db_client.query(query, &[]).await?;
        let mut synced = 0;
        
        for row in rows {
            let meter_name: String = row.get(0);
            let user_id: String = row.get(1);
            let value: f64 = row.get(2);
            let period_start: DateTime<Utc> = row.get(3);
            
            // Find meter config
            if let Some(meter_config) = self.meter_configs.iter().find(|m| m.name == meter_name) {
                match self.send_to_stripe_meter(&meter_config.stripe_meter_id, &user_id, value, period_start).await {
                    Ok(_) => {
                        // Mark as synced
                        let update_query = r#"
                            UPDATE meter_aggregations
                            SET synced_to_stripe = true, synced_at = NOW()
                            WHERE meter_name = $1 AND user_id = $2 AND period_start = $3
                        "#;
                        
                        self.db_client.execute(update_query, &[&meter_name, &user_id, &period_start]).await?;
                        synced += 1;
                    },
                    Err(e) => {
                        error!("Failed to sync {} for user {}: {}", meter_name, user_id, e);
                    }
                }
            }
        }
        
        info!("Synced {} records to Stripe", synced);
        Ok(synced)
    }
    
    async fn send_to_stripe_meter(&self, meter_id: &str, user_id: &str, value: f64, timestamp: DateTime<Utc>) -> Result<()> {
        // This is a simplified implementation
        // Real implementation would use Stripe's API to send meter events
        
        let url = format!("{}/v1/billing/meter_events", self.config.stripe.api_url);
        let payload = serde_json::json!({
            "event_name": meter_id,
            "payload": {
                "value": value.to_string(),
                "stripe_customer_id": user_id  // In practice, map user_id to Stripe customer ID
            },
            "timestamp": timestamp.timestamp()
        });
        
        let response = self.http_client
            .post(&url)
            .header("Authorization", format!("Bearer {}", self.config.stripe.secret_key))
            .header("Content-Type", "application/json")
            .json(&payload)
            .send()
            .await?;
        
        if !response.status().is_success() {
            return Err(anyhow!("Stripe API error: {}", response.status()));
        }
        
        Ok(())
    }
    
    pub async fn get_active_meters(&self) -> Result<Vec<String>> {
        Ok(self.meter_configs.iter().map(|m| m.name.clone()).collect())
    }
    
    pub async fn get_user_usage(&self, user_id: &str) -> Result<Vec<AggregationResult>> {
        let query = r#"
            SELECT meter_name, user_id, value, period_start, period_end
            FROM meter_aggregations
            WHERE user_id = $1
            ORDER BY period_start DESC
            LIMIT 100
        "#;
        
        let rows = self.db_client.query(query, &[&user_id]).await?;
        let mut results = Vec::new();
        
        for row in rows {
            results.push(AggregationResult {
                meter_name: row.get(0),
                user_id: row.get(1),
                value: row.get(2),
                period_start: row.get(3),
                period_end: row.get(4),
            });
        }
        
        Ok(results)
    }
}
