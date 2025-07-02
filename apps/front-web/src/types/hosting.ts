export interface Site {
  id: string;
  user_id: string;
  domain: string;
  subdomain: string;
  custom_domain?: string;
  status: 'active' | 'pending' | 'building' | 'error' | 'suspended';
  build_id?: string;
  version: number;
  config: SiteConfig;
  ssl_enabled: boolean;
  ssl_cert_path?: string;
  cdn_enabled: boolean;
  cdn_url?: string;
  analytics: boolean;
  last_deploy?: string;
  created_at: string;
  updated_at: string;
}

export interface SiteConfig {
  theme: string;
  custom_css: string;
  custom_js: string;
  seo_title: string;
  seo_description: string;
  seo_keywords: string[];
  redirect_rules: RedirectRule[];
  headers: Record<string, string>;
  error_pages: Record<string, string>;
  maintenance_mode: boolean;
  password?: string;
  geo_blocking: string[];
}

export interface RedirectRule {
  from: string;
  to: string;
  type: number;
  active: boolean;
}

export interface CreateSiteRequest {
  domain: string;
  custom_domain?: string;
  config: SiteConfig;
  ssl_enabled: boolean;
  cdn_enabled: boolean;
}

export interface UpdateSiteRequest {
  config?: SiteConfig;
  custom_domain?: string;
  ssl_enabled?: boolean;
  cdn_enabled?: boolean;
}

export interface Deployment {
  id: string;
  site_id: string;
  build_id: string;
  version: number;
  status: 'pending' | 'deploying' | 'deployed' | 'failed';
  files: DeploymentFile[];
  config: SiteConfig;
  stats: DeploymentStats;
  error?: string;
  started_at: string;
  ended_at?: string;
  created_at: string;
}

export interface DeploymentFile {
  path: string;
  content_type: string;
  size: number;
  hash: string;
  url: string;
}

export interface DeploymentStats {
  total_files: number;
  total_size: number;
  duration: number;
  cdn_push_time: number;
  cache_cleared: boolean;
}

export interface DomainConfig {
  id: string;
  site_id: string;
  domain: string;
  type: 'subdomain' | 'custom';
  status: 'pending' | 'active' | 'failed';
  dns_records: DNSRecord[];
  ssl_status: 'pending' | 'active' | 'failed';
  ssl_provider?: string;
  verified_at?: string;
  created_at: string;
  updated_at: string;
}

export interface DNSRecord {
  type: string;
  name: string;
  value: string;
  ttl: number;
  priority?: number;
}

export interface TrafficStats {
  site_id: string;
  date: string;
  page_views: number;
  unique_visits: number;
  bandwidth: number;
  requests: number;
  countries: Record<string, number>;
  referrers: Record<string, number>;
  pages: Record<string, number>;
  status_codes: Record<string, number>;
}

export interface CacheStatus {
  site_id: string;
  hit_rate: number;
  total_hits: number;
  total_misses: number;
  purged_at?: string;
  updated_at: string;
}
