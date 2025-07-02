use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};
use uuid::Uuid;
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UsageEvent {
    pub id: Option<Uuid>,
    pub event_type: String,
    pub user_id: String,
    pub project_id: Option<String>,
    pub site_id: Option<String>,
    pub timestamp: DateTime<Utc>,
    pub metadata: HashMap<String, serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EventBatch {
    pub events: Vec<UsageEvent>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MeterRecord {
    pub meter_name: String,
    pub value: f64,
    pub timestamp: DateTime<Utc>,
    pub dimensions: HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StripeUsageRecord {
    pub subscription_item: String,
    pub quantity: u64,
    pub timestamp: i64,
    pub action: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AggregationResult {
    pub meter_name: String,
    pub user_id: String,
    pub value: f64,
    pub period_start: DateTime<Utc>,
    pub period_end: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MeterConfig {
    pub name: String,
    pub stripe_meter_id: String,
    pub aggregation_type: AggregationType,
    pub event_filter: EventFilter,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AggregationType {
    Count,
    Sum,
    Max,
    Min,
    Average,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EventFilter {
    pub event_types: Vec<String>,
    pub metadata_filters: HashMap<String, String>,
}
