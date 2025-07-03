use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use std::collections::HashMap;
use uuid::Uuid;

// Core project models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Project {
    pub id: Uuid,
    pub name: String,
    pub description: Option<String>,
    pub owner_id: Uuid,
    pub team_id: Option<Uuid>,
    pub template_id: Option<Uuid>,
    pub theme_id: Option<Uuid>,
    pub visibility: ProjectVisibility,
    pub status: ProjectStatus,
    pub settings: serde_json::Value, // ProjectSettings serialized
    pub metadata: serde_json::Value, // ProjectMetadata serialized
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub last_accessed_at: Option<DateTime<Utc>>,
    pub deleted_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "project_visibility", rename_all = "lowercase")]
pub enum ProjectVisibility {
    Private,
    Public,
    Team,
    Organization,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "project_status", rename_all = "lowercase")]
pub enum ProjectStatus {
    Draft,
    Active,
    Published,
    Archived,
    Deleted,
    Suspended,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProjectSettings {
    pub auto_save: bool,
    pub version_control: bool,
    pub collaboration: bool,
    pub comments: bool,
    pub analytics: bool,
    pub seo_optimization: bool,
    pub performance_optimization: bool,
    pub accessibility_checks: bool,
    pub responsive_design: bool,
    pub theme_customization: bool,
    pub custom_code: bool,
    pub third_party_integrations: bool,
    pub backup_enabled: bool,
    pub cdn_enabled: bool,
    pub ssl_enabled: bool,
    pub custom_domains: Vec<String>,
    pub environment_vars: HashMap<String, String>,
    pub build_settings: BuildSettings,
    pub deployment_settings: DeploymentSettings,
    pub notification_settings: NotificationSettings,
    pub security_settings: SecuritySettings,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BuildSettings {
    pub auto_build: bool,
    pub build_on_push: bool,
    pub minify_assets: bool,
    pub optimize_images: bool,
    pub generate_sourcemaps: bool,
    pub enable_compression: bool,
    pub target_browsers: Vec<String>,
    pub framework: Option<String>,
    pub build_command: Option<String>,
    pub output_directory: String,
    pub environment: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DeploymentSettings {
    pub auto_deploy: bool,
    pub deployment_branch: String,
    pub preview_deployments: bool,
    pub rollback_enabled: bool,
    pub health_checks: bool,
    pub zero_downtime: bool,
    pub scaling: ScalingSettings,
    pub regions: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ScalingSettings {
    pub auto_scaling: bool,
    pub min_instances: u32,
    pub max_instances: u32,
    pub cpu_threshold: f32,
    pub memory_threshold: f32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NotificationSettings {
    pub email_notifications: bool,
    pub slack_notifications: bool,
    pub webhook_notifications: bool,
    pub build_notifications: bool,
    pub deployment_notifications: bool,
    pub error_notifications: bool,
    pub comment_notifications: bool,
    pub collaboration_notifications: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SecuritySettings {
    pub two_factor_required: bool,
    pub ip_whitelist: Vec<String>,
    pub access_control: bool,
    pub audit_logging: bool,
    pub vulnerability_scanning: bool,
    pub dependency_scanning: bool,
    pub code_scanning: bool,
    pub secrets_scanning: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProjectMetadata {
    pub tags: Vec<String>,
    pub category: Option<String>,
    pub language: String,
    pub framework: Option<String>,
    pub license: Option<String>,
    pub repository_url: Option<String>,
    pub homepage_url: Option<String>,
    pub documentation_url: Option<String>,
    pub support_url: Option<String>,
    pub keywords: Vec<String>,
    pub custom_fields: HashMap<String, serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProjectStats {
    pub total_pages: u32,
    pub total_components: u32,
    pub total_assets: u32,
    pub total_collaborators: u32,
    pub total_comments: u32,
    pub total_versions: u32,
    pub total_deployments: u32,
    pub size_bytes: u64,
    pub build_time_seconds: Option<f64>,
    pub last_build_status: Option<BuildStatus>,
    pub performance_score: Option<f32>,
    pub seo_score: Option<f32>,
    pub accessibility_score: Option<f32>,
    pub uptime_percentage: Option<f32>,
    pub avg_response_time_ms: Option<f32>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "build_status", rename_all = "lowercase")]
pub enum BuildStatus {
    Pending,
    Running,
    Success,
    Failed,
    Cancelled,
}

// User and permission models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct User {
    pub id: Uuid,
    pub email: String,
    pub username: String,
    pub display_name: Option<String>,
    pub avatar_url: Option<String>,
    pub role: UserRole,
    pub status: UserStatus,
    pub email_verified: bool,
    pub two_factor_enabled: bool,
    pub last_login_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "user_role", rename_all = "lowercase")]
pub enum UserRole {
    SuperAdmin,
    Admin,
    Manager,
    Developer,
    Designer,
    Viewer,
    Guest,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "user_status", rename_all = "lowercase")]
pub enum UserStatus {
    Active,
    Inactive,
    Suspended,
    Pending,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserPermissions {
    pub can_read: bool,
    pub can_write: bool,
    pub can_delete: bool,
    pub can_admin: bool,
    pub can_deploy: bool,
    pub can_invite: bool,
    pub can_comment: bool,
    pub can_approve: bool,
    pub specific_permissions: HashMap<String, bool>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "permission_level", rename_all = "lowercase")]
pub enum PermissionLevel {
    Owner,
    Admin,
    Write,
    Read,
    Comment,
    None,
}

// Page and component models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Page {
    pub id: Uuid,
    pub project_id: Uuid,
    pub name: String,
    pub path: String,
    pub title: Option<String>,
    pub description: Option<String>,
    pub content: serde_json::Value, // AST or component tree
    pub metadata: serde_json::Value, // SEO, social, etc.
    pub status: PageStatus,
    pub visibility: PageVisibility,
    pub parent_id: Option<Uuid>,
    pub sort_order: i32,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub published_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "page_status", rename_all = "lowercase")]
pub enum PageStatus {
    Draft,
    Published,
    Archived,
    Deleted,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "page_visibility", rename_all = "lowercase")]
pub enum PageVisibility {
    Public,
    Private,
    Protected,
    Hidden,
}

#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Component {
    pub id: Uuid,
    pub project_id: Uuid,
    pub name: String,
    pub type_name: String,
    pub category: String,
    pub description: Option<String>,
    pub props: serde_json::Value,
    pub styles: serde_json::Value,
    pub content: serde_json::Value,
    pub version: String,
    pub is_reusable: bool,
    pub is_published: bool,
    pub created_by: Uuid,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

// Asset and media models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Asset {
    pub id: Uuid,
    pub project_id: Uuid,
    pub name: String,
    pub file_name: String,
    pub file_path: String,
    pub file_size: i64,
    pub mime_type: String,
    pub asset_type: AssetType,
    pub metadata: serde_json::Value,
    pub alt_text: Option<String>,
    pub tags: Vec<String>,
    pub is_optimized: bool,
    pub cdn_url: Option<String>,
    pub uploaded_by: Uuid,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "asset_type", rename_all = "lowercase")]
pub enum AssetType {
    Image,
    Video,
    Audio,
    Document,
    Font,
    Icon,
    Style,
    Script,
    Other,
}

// Template and theme models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Template {
    pub id: Uuid,
    pub name: String,
    pub description: Option<String>,
    pub category: String,
    pub tags: Vec<String>,
    pub preview_url: Option<String>,
    pub content: serde_json::Value, // Template structure
    pub config: serde_json::Value,  // Template configuration
    pub is_premium: bool,
    pub is_featured: bool,
    pub price: Option<f64>,
    pub usage_count: i32,
    pub rating: Option<f32>,
    pub created_by: Uuid,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Theme {
    pub id: Uuid,
    pub name: String,
    pub description: Option<String>,
    pub colors: serde_json::Value,
    pub typography: serde_json::Value,
    pub spacing: serde_json::Value,
    pub borders: serde_json::Value,
    pub shadows: serde_json::Value,
    pub animations: serde_json::Value,
    pub custom_css: Option<String>,
    pub is_dark_mode: bool,
    pub is_premium: bool,
    pub created_by: Uuid,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

// Analytics and metrics models
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnalyticsMetrics {
    pub page_views: u64,
    pub unique_visitors: u64,
    pub bounce_rate: f32,
    pub avg_session_duration: f32,
    pub conversion_rate: f32,
    pub goal_completions: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PerformanceMetrics {
    pub lighthouse_score: f32,
    pub first_contentful_paint: f32,
    pub largest_contentful_paint: f32,
    pub cumulative_layout_shift: f32,
    pub first_input_delay: f32,
    pub total_blocking_time: f32,
    pub time_to_interactive: f32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EngagementMetrics {
    pub clicks: u64,
    pub scroll_depth: f32,
    pub time_on_page: f32,
    pub interactions: u64,
    pub form_submissions: u64,
    pub downloads: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ConversionMetrics {
    pub goal_completions: u64,
    pub conversion_rate: f32,
    pub revenue: f64,
    pub transactions: u64,
    pub average_order_value: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GeographicData {
    pub countries: HashMap<String, u64>,
    pub regions: HashMap<String, u64>,
    pub cities: HashMap<String, u64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DeviceData {
    pub desktop: u64,
    pub mobile: u64,
    pub tablet: u64,
    pub browsers: HashMap<String, u64>,
    pub operating_systems: HashMap<String, u64>,
    pub screen_resolutions: HashMap<String, u64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TimeRange {
    pub start: DateTime<Utc>,
    pub end: DateTime<Utc>,
    pub granularity: String, // "hour", "day", "week", "month"
}

// Pagination and filtering
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PaginationInfo {
    pub page: u32,
    pub limit: u32,
    pub total: u64,
    pub pages: u32,
    pub has_next: bool,
    pub has_prev: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProjectFilters {
    pub status: Option<Vec<ProjectStatus>>,
    pub visibility: Option<Vec<ProjectVisibility>>,
    pub tags: Option<Vec<String>>,
    pub categories: Option<Vec<String>>,
    pub date_range: Option<TimeRange>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "sort_order", rename_all = "lowercase")]
pub enum SortOrder {
    Asc,
    Desc,
}

// Export and import models
#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "export_format", rename_all = "lowercase")]
pub enum ExportFormat {
    Json,
    Zip,
    Tar,
    Git,
    Docker,
    Static,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "compression_type", rename_all = "lowercase")]
pub enum CompressionType {
    None,
    Gzip,
    Brotli,
    Lz4,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExportDestination {
    pub type_name: String, // "local", "s3", "github", "gitlab"
    pub config: HashMap<String, serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ImportSource {
    pub type_name: String, // "file", "url", "github", "gitlab"
    pub config: HashMap<String, serde_json::Value>,
}

// Collaboration models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Comment {
    pub id: Uuid,
    pub project_id: Uuid,
    pub page_id: Option<Uuid>,
    pub component_id: Option<Uuid>,
    pub user_id: Uuid,
    pub content: String,
    pub position: Option<serde_json::Value>, // x, y coordinates
    pub status: CommentStatus,
    pub parent_id: Option<Uuid>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub resolved_at: Option<DateTime<Utc>>,
    pub resolved_by: Option<Uuid>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "comment_status", rename_all = "lowercase")]
pub enum CommentStatus {
    Open,
    Resolved,
    Archived,
}

// Version control models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Version {
    pub id: Uuid,
    pub project_id: Uuid,
    pub version_number: String,
    pub description: Option<String>,
    pub changes: serde_json::Value,
    pub created_by: Uuid,
    pub created_at: DateTime<Utc>,
    pub is_major: bool,
    pub is_published: bool,
}

// Deployment models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Deployment {
    pub id: Uuid,
    pub project_id: Uuid,
    pub version_id: Option<Uuid>,
    pub environment: String,
    pub status: DeploymentStatus,
    pub url: Option<String>,
    pub build_logs: Option<String>,
    pub deploy_logs: Option<String>,
    pub started_at: DateTime<Utc>,
    pub completed_at: Option<DateTime<Utc>>,
    pub deployed_by: Uuid,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "deployment_status", rename_all = "lowercase")]
pub enum DeploymentStatus {
    Pending,
    Building,
    Deploying,
    Success,
    Failed,
    Cancelled,
    Rollback,
}

// Integration models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Integration {
    pub id: Uuid,
    pub project_id: Uuid,
    pub type_name: String,
    pub name: String,
    pub config: serde_json::Value,
    pub status: IntegrationStatus,
    pub last_sync_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "integration_status", rename_all = "lowercase")]
pub enum IntegrationStatus {
    Active,
    Inactive,
    Error,
    Pending,
}

// Domain and SSL models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Domain {
    pub id: Uuid,
    pub project_id: Uuid,
    pub domain_name: String,
    pub is_primary: bool,
    pub is_verified: bool,
    pub ssl_enabled: bool,
    pub ssl_certificate_id: Option<Uuid>,
    pub dns_records: serde_json::Value,
    pub created_at: DateTime<Utc>,
    pub verified_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct SslCertificate {
    pub id: Uuid,
    pub domain_id: Uuid,
    pub certificate: String,
    pub private_key: String,
    pub chain: Option<String>,
    pub expires_at: DateTime<Utc>,
    pub auto_renew: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

// Form and data models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Form {
    pub id: Uuid,
    pub project_id: Uuid,
    pub name: String,
    pub description: Option<String>,
    pub fields: serde_json::Value,
    pub settings: serde_json::Value,
    pub submissions_count: i32,
    pub is_active: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct FormSubmission {
    pub id: Uuid,
    pub form_id: Uuid,
    pub data: serde_json::Value,
    pub ip_address: Option<String>,
    pub user_agent: Option<String>,
    pub submitted_at: DateTime<Utc>,
}

// API and function models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct ApiEndpoint {
    pub id: Uuid,
    pub project_id: Uuid,
    pub name: String,
    pub path: String,
    pub method: String,
    pub description: Option<String>,
    pub request_schema: Option<serde_json::Value>,
    pub response_schema: Option<serde_json::Value>,
    pub handler_code: String,
    pub is_active: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct ServerlessFunction {
    pub id: Uuid,
    pub project_id: Uuid,
    pub name: String,
    pub description: Option<String>,
    pub runtime: String,
    pub code: String,
    pub environment_vars: serde_json::Value,
    pub timeout: i32,
    pub memory: i32,
    pub is_deployed: bool,
    pub deployment_url: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub last_deployed_at: Option<DateTime<Utc>>,
}

// Database connection models
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct DatabaseConnection {
    pub id: Uuid,
    pub project_id: Uuid,
    pub name: String,
    pub database_type: String,
    pub connection_string: String, // Encrypted
    pub schema_config: Option<serde_json::Value>,
    pub is_active: bool,
    pub last_tested_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}
