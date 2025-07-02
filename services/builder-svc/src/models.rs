use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};
use std::collections::HashMap;

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
    pub assets: Vec<Asset>,
    pub metadata: SiteMetadata,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Section {
    pub id: String,
    pub component_type: ComponentType,
    pub content: serde_json::Value,
    pub styles: serde_json::Value,
    pub order: i32,
    pub responsive: ResponsiveConfig,
    pub animations: Vec<Animation>,
    pub interactions: Vec<Interaction>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum ComponentType {
    Header,
    Navigation,
    Hero,
    TextBlock,
    ImageGallery,
    ContactForm,
    Testimonials,
    Pricing,
    Features,
    Footer,
    Custom,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ResponsiveConfig {
    pub mobile: Option<serde_json::Value>,
    pub tablet: Option<serde_json::Value>,
    pub desktop: Option<serde_json::Value>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Animation {
    pub name: String,
    pub trigger: AnimationTrigger,
    pub duration: f32,
    pub delay: f32,
    pub easing: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum AnimationTrigger {
    OnLoad,
    OnScroll,
    OnHover,
    OnClick,
    OnVisible,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Interaction {
    pub event: String,
    pub action: InteractionAction,
    pub target: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum InteractionAction {
    Navigate(String),
    ShowModal(String),
    ScrollTo(String),
    ToggleClass(String),
    Custom(String),
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Theme {
    pub name: String,
    pub colors: ColorPalette,
    pub fonts: FontConfig,
    pub spacing: SpacingConfig,
    pub borders: BorderConfig,
    pub shadows: ShadowConfig,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ColorPalette {
    pub primary: String,
    pub secondary: String,
    pub accent: String,
    pub background: String,
    pub surface: String,
    pub text: String,
    pub text_secondary: String,
    pub success: String,
    pub warning: String,
    pub error: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct FontConfig {
    pub primary: String,
    pub secondary: String,
    pub headings: String,
    pub body: String,
    pub sizes: FontSizes,
    pub weights: FontWeights,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct FontSizes {
    pub xs: String,
    pub sm: String,
    pub base: String,
    pub lg: String,
    pub xl: String,
    pub xl2: String,
    pub xl3: String,
    pub xl4: String,
    pub xl5: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct FontWeights {
    pub thin: u16,
    pub light: u16,
    pub normal: u16,
    pub medium: u16,
    pub semibold: u16,
    pub bold: u16,
    pub extrabold: u16,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SpacingConfig {
    pub xs: String,
    pub sm: String,
    pub md: String,
    pub lg: String,
    pub xl: String,
    pub xl2: String,
    pub xl3: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BorderConfig {
    pub radius: BorderRadius,
    pub width: BorderWidth,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BorderRadius {
    pub sm: String,
    pub md: String,
    pub lg: String,
    pub xl: String,
    pub full: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BorderWidth {
    pub thin: String,
    pub base: String,
    pub thick: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ShadowConfig {
    pub sm: String,
    pub md: String,
    pub lg: String,
    pub xl: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SiteSettings {
    pub language: String,
    pub favicon: Option<String>,
    pub meta_description: Option<String>,
    pub meta_keywords: Vec<String>,
    pub analytics: AnalyticsConfig,
    pub seo: SeoConfig,
    pub performance: PerformanceConfig,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AnalyticsConfig {
    pub google_analytics_id: Option<String>,
    pub facebook_pixel_id: Option<String>,
    pub custom_scripts: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SeoConfig {
    pub title_template: String,
    pub open_graph: OpenGraphConfig,
    pub twitter_card: TwitterCardConfig,
    pub canonical_url: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct OpenGraphConfig {
    pub title: Option<String>,
    pub description: Option<String>,
    pub image: Option<String>,
    pub url: Option<String>,
    pub site_name: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct TwitterCardConfig {
    pub card_type: String,
    pub site: Option<String>,
    pub creator: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct PerformanceConfig {
    pub lazy_loading: bool,
    pub image_optimization: bool,
    pub code_splitting: bool,
    pub compression: bool,
    pub caching_strategy: CachingStrategy,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum CachingStrategy {
    None,
    Static,
    Dynamic,
    Hybrid,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Asset {
    pub id: String,
    pub name: String,
    pub asset_type: AssetType,
    pub url: String,
    pub size: u64,
    pub mime_type: String,
    pub optimized_versions: HashMap<String, String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum AssetType {
    Image,
    Video,
    Audio,
    Font,
    Icon,
    Document,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SiteMetadata {
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub version: String,
    pub framework: String,
    pub build_tool: String,
    pub dependencies: Vec<Dependency>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Dependency {
    pub name: String,
    pub version: String,
    pub dependency_type: DependencyType,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum DependencyType {
    Framework,
    Library,
    Plugin,
    Tool,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum BuildType {
    Development,
    Staging,
    Production,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildOptions {
    pub minify: bool,
    pub optimize_images: bool,
    pub generate_sitemap: bool,
    pub enable_pwa: bool,
    pub target_framework: TargetFramework,
    pub deployment_target: DeploymentTarget,
    pub custom_domain: Option<String>,
    pub environment_vars: HashMap<String, String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum TargetFramework {
    React,
    Vue,
    Angular,
    Svelte,
    VanillaJS,
    NextJS,
    Nuxt,
    Gatsby,
    Static,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum DeploymentTarget {
    Vercel,
    Netlify,
    AWS,
    Azure,
    GCP,
    Custom(String),
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildResult {
    pub build_id: Uuid,
    pub project_id: Uuid,
    pub status: BuildStatus,
    pub output_url: Option<String>,
    pub preview_url: Option<String>,
    pub logs: Vec<BuildLog>,
    pub artifacts: Vec<BuildArtifact>,
    pub metrics: BuildMetrics,
    pub created_at: DateTime<Utc>,
    pub completed_at: Option<DateTime<Utc>>,
    pub error: Option<BuildError>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum BuildStatus {
    Queued,
    Building,
    Completed,
    Failed,
    Cancelled,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildLog {
    pub timestamp: DateTime<Utc>,
    pub level: LogLevel,
    pub message: String,
    pub phase: BuildPhase,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum LogLevel {
    Debug,
    Info,
    Warning,
    Error,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum BuildPhase {
    Initialization,
    DependencyInstall,
    CodeGeneration,
    Compilation,
    Optimization,
    AssetProcessing,
    Deployment,
    Cleanup,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildArtifact {
    pub name: String,
    pub artifact_type: ArtifactType,
    pub url: String,
    pub size: u64,
    pub checksum: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum ArtifactType {
    Bundle,
    SourceMap,
    Assets,
    Manifest,
    ServiceWorker,
    Sitemap,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildMetrics {
    pub build_time: u64,
    pub bundle_size: u64,
    pub asset_count: u32,
    pub lighthouse_score: Option<LighthouseScore>,
    pub performance_metrics: PerformanceMetrics,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct LighthouseScore {
    pub performance: f32,
    pub accessibility: f32,
    pub best_practices: f32,
    pub seo: f32,
    pub pwa: Option<f32>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct PerformanceMetrics {
    pub first_contentful_paint: Option<u32>,
    pub largest_contentful_paint: Option<u32>,
    pub first_input_delay: Option<u32>,
    pub cumulative_layout_shift: Option<f32>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildError {
    pub error_type: ErrorType,
    pub message: String,
    pub stack_trace: Option<String>,
    pub file: Option<String>,
    pub line: Option<u32>,
    pub column: Option<u32>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum ErrorType {
    SyntaxError,
    DependencyError,
    ConfigurationError,
    NetworkError,
    PermissionError,
    ResourceError,
    TimeoutError,
    Unknown,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ComponentTemplate {
    pub id: String,
    pub name: String,
    pub description: String,
    pub category: ComponentCategory,
    pub framework: TargetFramework,
    pub template: String,
    pub styles: String,
    pub props: Vec<ComponentProp>,
    pub preview_image: Option<String>,
    pub is_premium: bool,
    pub tags: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum ComponentCategory {
    Layout,
    Content,
    Navigation,
    Form,
    Media,
    Social,
    Ecommerce,
    Marketing,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ComponentProp {
    pub name: String,
    pub prop_type: PropType,
    pub default_value: Option<serde_json::Value>,
    pub required: bool,
    pub description: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum PropType {
    String,
    Number,
    Boolean,
    Array,
    Object,
    Color,
    Image,
    Url,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildJob {
    pub id: Uuid,
    pub project_id: Uuid,
    pub user_id: Uuid,
    pub build_request: BuildRequest,
    pub status: BuildStatus,
    pub priority: u8,
    pub scheduled_at: DateTime<Utc>,
    pub started_at: Option<DateTime<Utc>>,
    pub completed_at: Option<DateTime<Utc>>,
    pub worker_id: Option<String>,
    pub retry_count: u8,
    pub max_retries: u8,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BuildQueue {
    pub pending: Vec<BuildJob>,
    pub running: Vec<BuildJob>,
    pub completed: Vec<BuildJob>,
    pub failed: Vec<BuildJob>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct WebsiteOutput {
    pub html_files: HashMap<String, String>,
    pub css_files: HashMap<String, String>,
    pub js_files: HashMap<String, String>,
    pub assets: HashMap<String, Vec<u8>>,
    pub manifest: SiteManifest,
    pub sitemap: Option<String>,
    pub robots_txt: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SiteManifest {
    pub name: String,
    pub short_name: String,
    pub description: String,
    pub start_url: String,
    pub display: String,
    pub theme_color: String,
    pub background_color: String,
    pub icons: Vec<ManifestIcon>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ManifestIcon {
    pub src: String,
    pub sizes: String,
    pub icon_type: String,
    pub purpose: Option<String>,
}
