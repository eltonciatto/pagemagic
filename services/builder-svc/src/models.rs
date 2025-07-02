use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildRequest {
    pub project_id: Uuid,
    pub user_id: Uuid,
    pub site_data: SiteData,
    pub build_type: BuildType,
    pub options: BuildOptions,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SiteData {
    pub name: String,
    pub description: String,
    pub sections: Vec<Section>,
    pub theme: Theme,
    pub settings: SiteSettings,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Section {
    pub id: String,
    pub component_type: String,
    pub content: serde_json::Value,
    pub styles: serde_json::Value,
    pub order: i32,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Theme {
    pub name: String,
    pub colors: ColorPalette,
    pub fonts: FontConfig,
    pub spacing: SpacingConfig,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ColorPalette {
    pub primary: String,
    pub secondary: String,
    pub accent: String,
    pub background: String,
    pub text: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct FontConfig {
    pub heading: String,
    pub body: String,
    pub mono: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SpacingConfig {
    pub small: String,
    pub medium: String,
    pub large: String,
    pub xl: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SiteSettings {
    pub seo: SeoConfig,
    pub analytics: AnalyticsConfig,
    pub performance: PerformanceConfig,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SeoConfig {
    pub title: String,
    pub description: String,
    pub keywords: Vec<String>,
    pub og_image: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AnalyticsConfig {
    pub google_analytics: Option<String>,
    pub facebook_pixel: Option<String>,
    pub custom_scripts: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct PerformanceConfig {
    pub lazy_loading: bool,
    pub image_optimization: bool,
    pub minify_css: bool,
    pub minify_js: bool,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum BuildType {
    StaticSite,
    ReactApp,
    NextjsApp,
    ComponentLibrary,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildOptions {
    pub target_platform: TargetPlatform,
    pub optimization_level: OptimizationLevel,
    pub include_source_maps: bool,
    pub custom_css: Option<String>,
    pub custom_js: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum TargetPlatform {
    Web,
    Mobile,
    Desktop,
    Universal,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum OptimizationLevel {
    Development,
    Production,
    Maximum,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildResponse {
    pub build_id: Uuid,
    pub status: BuildStatus,
    pub message: String,
    pub artifacts: Option<BuildArtifacts>,
    pub created_at: DateTime<Utc>,
    pub completed_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum BuildStatus {
    Queued,
    InProgress,
    Completed,
    Failed,
    Cancelled,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildArtifacts {
    pub html_files: Vec<String>,
    pub css_files: Vec<String>,
    pub js_files: Vec<String>,
    pub asset_files: Vec<String>,
    pub deployment_url: Option<String>,
    pub bundle_size: u64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ComponentBuildRequest {
    pub component_type: String,
    pub props: serde_json::Value,
    pub children: Option<Vec<ComponentBuildRequest>>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ComponentBuildResponse {
    pub html: String,
    pub css: String,
    pub js: Option<String>,
    pub dependencies: Vec<String>,
}
