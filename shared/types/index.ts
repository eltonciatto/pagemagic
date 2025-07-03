// Page Magic - Shared TypeScript Types
// Tipos compartilhados entre front-end, back-end e mobile

// ==========================================
// BASIC TYPES
// ==========================================

export type UUID = string;
export type Timestamp = string; // ISO 8601 string
export type JSON = Record<string, any>;

// ==========================================
// USER & AUTHENTICATION
// ==========================================

export type AuthProvider = 'email' | 'google' | 'github' | 'apple';
export type UserStatus = 'active' | 'inactive' | 'suspended' | 'deleted';
export type SubscriptionPlan = 'free' | 'basic' | 'pro' | 'enterprise';

export interface User {
  id: UUID;
  email: string;
  email_verified: boolean;
  first_name?: string;
  last_name?: string;
  avatar_url?: string;
  status: UserStatus;
  locale: string;
  timezone: string;
  mobile_verified?: boolean;
  mobile_number?: string;
  push_notifications_enabled?: boolean;
  last_login_at?: Timestamp;
  created_at: Timestamp;
  updated_at: Timestamp;
}

export interface UserAuthProvider {
  id: UUID;
  user_id: UUID;
  provider: AuthProvider;
  provider_user_id: string;
  provider_email?: string;
  provider_data?: JSON;
  created_at: Timestamp;
  updated_at: Timestamp;
}

export interface Subscription {
  id: UUID;
  user_id: UUID;
  stripe_subscription_id?: string;
  stripe_customer_id?: string;
  plan: SubscriptionPlan;
  status: string;
  current_period_start?: Timestamp;
  current_period_end?: Timestamp;
  cancel_at_period_end: boolean;
  trial_start?: Timestamp;
  trial_end?: Timestamp;
  metadata?: JSON;
  created_at: Timestamp;
  updated_at: Timestamp;
}

// ==========================================
// PROJECTS & SITES
// ==========================================

export type ProjectStatus = 'draft' | 'building' | 'published' | 'archived' | 'error';
export type SiteType = 'landing_page' | 'blog' | 'portfolio' | 'ecommerce' | 'store' | 'custom';

export interface Project {
  id: UUID;
  user_id: UUID;
  team_id?: UUID;
  name: string;
  description?: string;
  type: SiteType;
  status: ProjectStatus;
  domain?: string;
  subdomain?: string;
  custom_domain?: string;
  ssl_enabled: boolean;
  password_protected: boolean;
  seo_title?: string;
  seo_description?: string;
  seo_keywords?: string[];
  favicon_url?: string;
  social_image_url?: string;
  analytics_enabled: boolean;
  analytics_tracking_id?: string;
  analytics_config?: JSON;
  conversion_goals?: JSON[];
  ab_testing_enabled?: boolean;
  cache_settings?: JSON;
  cdn_settings?: JSON;
  settings: JSON;
  metadata: JSON;
  created_at: Timestamp;
  updated_at: Timestamp;
  published_at?: Timestamp;
  archived_at?: Timestamp;
}

export interface ProjectBuild {
  id: UUID;
  project_id: UUID;
  version_number: number;
  docker_image_id?: string;
  build_log?: string;
  build_status: string;
  build_started_at?: Timestamp;
  build_completed_at?: Timestamp;
  is_current: boolean;
  size_bytes?: number;
  created_at: Timestamp;
}

// ==========================================
// PAGES & CONTENT
// ==========================================

export type SectionType = 'hero' | 'about' | 'services' | 'portfolio' | 'testimonials' | 'contact' | 'footer' | 'custom';

export interface Page {
  id: UUID;
  project_id: UUID;
  name: string;
  slug: string;
  title?: string;
  description?: string;
  is_home: boolean;
  is_published: boolean;
  order_index: number;
  seo_title?: string;
  seo_description?: string;
  seo_keywords?: string[];
  created_at: Timestamp;
  updated_at: Timestamp;
}

export interface Section {
  id: UUID;
  page_id: UUID;
  type: SectionType;
  name: string;
  order_index: number;
  html_content?: string;
  css_styles?: string;
  js_scripts?: string;
  data: JSON;
  settings: JSON;
  is_visible: boolean;
  lazy_loading?: boolean;
  priority?: number;
  cache_duration_seconds?: number;
  created_at: Timestamp;
  updated_at: Timestamp;
}

// ==========================================
// AI GENERATION
// ==========================================

export type GenerationStatus = 'queued' | 'processing' | 'completed' | 'failed' | 'cancelled';
export type AIModel = 'llama-3-70b' | 'llama-3-8b' | 'gpt-4' | 'gpt-4-turbo' | 'gpt-3.5-turbo' | 'claude-3-sonnet' | 'claude-3-opus' | 'claude-3-haiku';

export interface AIGeneration {
  id: UUID;
  user_id: UUID;
  project_id?: UUID;
  prompt: string;
  model: AIModel;
  status: GenerationStatus;
  response_data?: JSON;
  tokens_used?: number;
  processing_time_ms?: number;
  error_message?: string;
  context_data?: JSON;
  prompt_template?: string;
  generation_mode?: string;
  parent_generation_id?: UUID;
  metadata?: JSON;
  created_at: Timestamp;
  completed_at?: Timestamp;
}

export interface PromptTemplate {
  id: UUID;
  name: string;
  description?: string;
  category: string;
  template: string;
  variables: string[];
  is_public: boolean;
  created_by?: UUID;
  usage_count: number;
  rating: number;
  created_at: Timestamp;
  updated_at: Timestamp;
}

// ==========================================
// MEDIA & ASSETS
// ==========================================

export type MediaType = 'image' | 'video' | 'audio' | 'document' | 'font' | 'other';

export interface MediaFile {
  id: UUID;
  user_id: UUID;
  project_id?: UUID;
  filename: string;
  original_filename: string;
  mime_type: string;
  type: MediaType;
  size_bytes: number;
  width?: number;
  height?: number;
  duration_seconds?: number;
  url: string;
  thumbnail_url?: string;
  alt_text?: string;
  description?: string;
  metadata: JSON;
  created_at: Timestamp;
}

// ==========================================
// DOMAINS & SSL
// ==========================================

export type DomainStatus = 'pending' | 'active' | 'expired' | 'suspended' | 'error';

export interface CustomDomain {
  id: UUID;
  project_id: UUID;
  domain: string;
  status: DomainStatus;
  dns_configured: boolean;
  ssl_certificate_id?: string;
  ssl_issued_at?: Timestamp;
  ssl_expires_at?: Timestamp;
  verified_at?: Timestamp;
  last_check_at?: Timestamp;
  error_message?: string;
  registrar_id?: UUID;
  auto_renew?: boolean;
  expires_at?: Timestamp;
  nameservers?: string[];
  dns_records?: JSON[];
  created_at: Timestamp;
  updated_at: Timestamp;
}

export interface DomainRegistrar {
  id: UUID;
  name: string;
  api_endpoint: string;
  supported_tlds: string[];
  is_active: boolean;
  credentials_schema?: JSON;
  created_at: Timestamp;
}

// ==========================================
// BILLING & USAGE
// ==========================================

export type MeterType = 'page_generate' | 'ai_token' | 'container_hours' | 'storage_gb';

export interface UsageRecord {
  id: UUID;
  user_id: UUID;
  subscription_id?: UUID;
  meter_type: MeterType;
  quantity: number;
  timestamp: Timestamp;
  stripe_usage_record_id?: string;
  metadata?: JSON;
  created_at: Timestamp;
}

// ==========================================
// MOBILE & PUSH NOTIFICATIONS
// ==========================================

export type MobilePlatform = 'ios' | 'android' | 'expo';

export interface MobilePushToken {
  id: UUID;
  user_id: UUID;
  token: string;
  platform: MobilePlatform;
  app_version?: string;
  device_id?: string;
  is_active: boolean;
  created_at: Timestamp;
  updated_at: Timestamp;
}

export interface PushNotification {
  id: UUID;
  user_id: UUID;
  title: string;
  body: string;
  data?: JSON;
  sent_at?: Timestamp;
  delivered_at?: Timestamp;
  clicked_at?: Timestamp;
  error_message?: string;
  created_at: Timestamp;
}

// ==========================================
// TEAMS & COLLABORATION
// ==========================================

export interface Team {
  id: UUID;
  name: string;
  slug: string;
  owner_id: UUID;
  plan: SubscriptionPlan;
  settings: JSON;
  created_at: Timestamp;
  updated_at: Timestamp;
}

export interface TeamMember {
  id: UUID;
  team_id: UUID;
  user_id: UUID;
  role: string; // 'owner' | 'admin' | 'editor' | 'viewer'
  invited_by?: UUID;
  invited_at?: Timestamp;
  joined_at?: Timestamp;
  created_at: Timestamp;
}

// ==========================================
// E-COMMERCE
// ==========================================

export interface Product {
  id: UUID;
  project_id: UUID;
  name: string;
  description?: string;
  price: number;
  currency: string;
  sku?: string;
  inventory_quantity: number;
  track_inventory: boolean;
  weight?: number;
  dimensions?: JSON;
  images: string[];
  is_active: boolean;
  seo_title?: string;
  seo_description?: string;
  metadata: JSON;
  created_at: Timestamp;
  updated_at: Timestamp;
}

export interface ProductVariant {
  id: UUID;
  product_id: UUID;
  name: string;
  price?: number;
  sku?: string;
  inventory_quantity: number;
  attributes: JSON;
  is_active: boolean;
  created_at: Timestamp;
  updated_at: Timestamp;
}

// ==========================================
// INTERNATIONALIZATION
// ==========================================

export interface Language {
  code: string;
  name: string;
  native_name: string;
  is_active: boolean;
  is_rtl: boolean;
  created_at: Timestamp;
}

export interface Translation {
  id: UUID;
  key: string;
  language_code: string;
  value: string;
  context?: string;
  project_id?: UUID;
  created_at: Timestamp;
  updated_at: Timestamp;
}

// ==========================================
// WEBHOOKS & API
// ==========================================

export interface Webhook {
  id: UUID;
  user_id: UUID;
  project_id?: UUID;
  url: string;
  events: string[];
  secret: string;
  is_active: boolean;
  last_triggered_at?: Timestamp;
  created_at: Timestamp;
  updated_at: Timestamp;
}

export interface WebhookDelivery {
  id: UUID;
  webhook_id: UUID;
  event_type: string;
  payload: JSON;
  response_status?: number;
  response_body?: string;
  delivery_duration_ms?: number;
  attempted_at: Timestamp;
  delivered_at?: Timestamp;
}

export interface APIKey {
  id: UUID;
  user_id: UUID;
  name: string;
  key_hash: string;
  key_prefix: string;
  scopes: string[];
  last_used_at?: Timestamp;
  expires_at?: Timestamp;
  is_active: boolean;
  created_at: Timestamp;
}

// ==========================================
// AUDIT & LOGGING
// ==========================================

export type AuditAction = 'create' | 'update' | 'delete' | 'login' | 'logout' | 'publish' | 'unpublish';

export interface AuditLog {
  id: UUID;
  user_id?: UUID;
  action: AuditAction;
  resource_type: string;
  resource_id?: UUID;
  old_values?: JSON;
  new_values?: JSON;
  ip_address?: string;
  user_agent?: string;
  created_at: Timestamp;
}

// ==========================================
// API REQUEST/RESPONSE TYPES
// ==========================================

export interface APIResponse<T = any> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  };
  meta?: {
    pagination?: {
      page: number;
      limit: number;
      total: number;
      totalPages: number;
    };
  };
}

export interface PaginationParams {
  page?: number;
  limit?: number;
  sort?: string;
  order?: 'asc' | 'desc';
}

export interface FilterParams {
  search?: string;
  status?: string;
  type?: string;
  created_after?: string;
  created_before?: string;
}

// ==========================================
// GENERATION REQUESTS
// ==========================================

export interface GenerationRequest {
  project_id?: UUID;
  prompt: string;
  model?: AIModel;
  locale?: string;
  context?: JSON;
  template_id?: UUID;
  generation_mode?: 'standard' | 'iterative' | 'collaborative';
  parent_generation_id?: UUID;
}

export interface GenerationResponse {
  generation_id: UUID;
  status: GenerationStatus;
  estimated_completion_time?: number;
}

// ==========================================
// BUILDER SERVICE TYPES
// ==========================================

export interface BuilderRequest {
  sections: Array<{
    id: string;
    type: SectionType;
    html: string;
    css?: string;
    js?: string;
    data?: JSON;
  }>;
  css?: string;
  js?: string;
  metadata?: JSON;
}

export interface BuilderResponse {
  ast_id: UUID;
  warnings: string[];
  optimizations?: string[];
}

// ==========================================
// ANALYTICS TYPES
// ==========================================

export interface AnalyticsEvent {
  project_id: UUID;
  event_type: string;
  event_count: number;
  event_date: string;
  metadata: JSON;
}

export interface PageView {
  project_id: UUID;
  page_path: string;
  visitor_id: string;
  session_id: string;
  ip_address?: string;
  user_agent?: string;
  referrer?: string;
  country?: string;
  city?: string;
  device_type?: 'desktop' | 'mobile' | 'tablet';
  browser?: string;
  os?: string;
  duration_seconds?: number;
  metadata?: JSON;
}

export interface PerformanceMetrics {
  project_id: UUID;
  page_path: string;
  visitor_id: string;
  lcp?: number; // Largest Contentful Paint
  fid?: number; // First Input Delay
  cls?: number; // Cumulative Layout Shift
  fcp?: number; // First Contentful Paint
  ttfb?: number; // Time to First Byte
  dom_load_time?: number;
  window_load_time?: number;
  connection_type?: string;
  device_memory?: number;
  metadata?: JSON;
}

// ==========================================
// CACHE TYPES
// ==========================================

export interface CacheInvalidation {
  id: UUID;
  project_id: UUID;
  cache_key: string;
  invalidation_type: 'path' | 'tag' | 'full';
  status: 'pending' | 'processing' | 'completed' | 'failed';
  provider: string;
  provider_request_id?: string;
  error_message?: string;
  created_at: Timestamp;
  completed_at?: Timestamp;
}

// ==========================================
// EDITOR TYPES
// ==========================================

export type ComponentType =
  | 'text'
  | 'heading'
  | 'image'
  | 'button'
  | 'container'
  | 'grid'
  | 'form'
  | 'input'
  | 'textarea'
  | 'select'
  | 'video'
  | 'iframe'
  | 'spacer'
  | 'divider'
  | 'icon'
  | 'carousel'
  | 'modal'
  | 'tabs'
  | 'accordion'
  | 'navbar'
  | 'footer'
  | 'sidebar';

export interface Component {
  id: string;
  type: ComponentType;
  props: Record<string, any>;
  styles: Record<string, any>;
  children: Component[];
  parent_id?: string;
  order_index?: number;
}

export interface ComponentProperty {
  name: string;
  type: 'text' | 'number' | 'boolean' | 'color' | 'select' | 'file' | 'array';
  label: string;
  description?: string;
  default?: any;
  options?: Array<{ label: string; value: any }>;
  required?: boolean;
  validation?: {
    min?: number;
    max?: number;
    pattern?: string;
    message?: string;
  };
}

export interface Theme {
  id: string;
  name: string;
  description?: string;
  colors: {
    primary: string;
    secondary: string;
    accent: string;
    background: string;
    surface: string;
    text: string;
    textSecondary: string;
    border: string;
    success: string;
    warning: string;
    error: string;
  };
  typography: {
    fontFamily: string;
    headingFontFamily?: string;
    fontSizes: Record<string, string>;
    fontWeights: Record<string, number>;
    lineHeights: Record<string, number>;
  };
  spacing: Record<string, string>;
  borderRadius: Record<string, string>;
  shadows: Record<string, string>;
  breakpoints: {
    sm: string;
    md: string;
    lg: string;
    xl: string;
  };
}

export interface Template {
  id: string;
  name: string;
  description?: string;
  category: string;
  thumbnail_url?: string;
  preview_url?: string;
  is_premium: boolean;
  tags: string[];
  sections: Section[];
  theme: Theme;
  metadata: JSON;
  created_at: Timestamp;
  updated_at: Timestamp;
}

// ==========================================
// EXPORT ALL TYPES
// ==========================================

export type { };
