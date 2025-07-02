package models

import (
	"time"
)

// Site representa um site hospedado
type Site struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Domain      string    `json:"domain" db:"domain"`
	Subdomain   string    `json:"subdomain" db:"subdomain"`
	CustomDomain string   `json:"custom_domain" db:"custom_domain"`
	Status      string    `json:"status" db:"status"` // "active", "suspended", "building", "error"
	BuildID     string    `json:"build_id" db:"build_id"`
	Version     int       `json:"version" db:"version"`
	Config      SiteConfig `json:"config" db:"config"`
	SSLEnabled  bool      `json:"ssl_enabled" db:"ssl_enabled"`
	SSLCertPath string    `json:"ssl_cert_path" db:"ssl_cert_path"`
	CDNEnabled  bool      `json:"cdn_enabled" db:"cdn_enabled"`
	CDNUrl      string    `json:"cdn_url" db:"cdn_url"`
	Analytics   bool      `json:"analytics" db:"analytics"`
	LastDeploy  *time.Time `json:"last_deploy" db:"last_deploy"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SiteConfig contém configurações específicas do site
type SiteConfig struct {
	Theme           string            `json:"theme"`
	CustomCSS       string            `json:"custom_css"`
	CustomJS        string            `json:"custom_js"`
	SEOTitle        string            `json:"seo_title"`
	SEODescription  string            `json:"seo_description"`
	SEOKeywords     []string          `json:"seo_keywords"`
	RedirectRules   []RedirectRule    `json:"redirect_rules"`
	Headers         map[string]string `json:"headers"`
	ErrorPages      map[string]string `json:"error_pages"`
	MaintenanceMode bool              `json:"maintenance_mode"`
	Password        string            `json:"password,omitempty"`
	GeoBlocking     []string          `json:"geo_blocking"`
}

// RedirectRule define uma regra de redirecionamento
type RedirectRule struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Type   int    `json:"type"` // 301, 302, etc.
	Active bool   `json:"active"`
}

// Deployment representa um deployment de site
type Deployment struct {
	ID        string             `json:"id" db:"id"`
	SiteID    string             `json:"site_id" db:"site_id"`
	BuildID   string             `json:"build_id" db:"build_id"`
	Version   int                `json:"version" db:"version"`
	Status    string             `json:"status" db:"status"` // "pending", "deploying", "deployed", "failed"
	Files     []DeploymentFile   `json:"files" db:"files"`
	Config    SiteConfig         `json:"config" db:"config"`
	Stats     DeploymentStats    `json:"stats" db:"stats"`
	Error     string             `json:"error" db:"error"`
	StartedAt time.Time          `json:"started_at" db:"started_at"`
	EndedAt   *time.Time         `json:"ended_at" db:"ended_at"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
}

// DeploymentFile representa um arquivo no deployment
type DeploymentFile struct {
	Path        string `json:"path"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Hash        string `json:"hash"`
	URL         string `json:"url"`
}

// DeploymentStats contém estatísticas do deployment
type DeploymentStats struct {
	TotalFiles   int           `json:"total_files"`
	TotalSize    int64         `json:"total_size"`
	Duration     time.Duration `json:"duration"`
	CDNPushTime  time.Duration `json:"cdn_push_time"`
	CacheCleared bool          `json:"cache_cleared"`
}

// DomainConfig representa configuração de domínio
type DomainConfig struct {
	ID           string    `json:"id" db:"id"`
	SiteID       string    `json:"site_id" db:"site_id"`
	Domain       string    `json:"domain" db:"domain"`
	Type         string    `json:"type" db:"type"` // "subdomain", "custom"
	Status       string    `json:"status" db:"status"` // "pending", "active", "failed"
	DNSRecords   []DNSRecord `json:"dns_records" db:"dns_records"`
	SSLStatus    string    `json:"ssl_status" db:"ssl_status"` // "pending", "active", "failed"
	SSLProvider  string    `json:"ssl_provider" db:"ssl_provider"`
	VerifiedAt   *time.Time `json:"verified_at" db:"verified_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// DNSRecord representa um registro DNS
type DNSRecord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority,omitempty"`
}

// TrafficStats representa estatísticas de tráfego
type TrafficStats struct {
	SiteID       string            `json:"site_id"`
	Date         time.Time         `json:"date"`
	PageViews    int64             `json:"page_views"`
	UniqueVisits int64             `json:"unique_visits"`
	Bandwidth    int64             `json:"bandwidth"`
	Requests     int64             `json:"requests"`
	Countries    map[string]int64  `json:"countries"`
	Referrers    map[string]int64  `json:"referrers"`
	Pages        map[string]int64  `json:"pages"`
	StatusCodes  map[string]int64  `json:"status_codes"`
}

// CacheStatus representa status do cache
type CacheStatus struct {
	SiteID      string    `json:"site_id"`
	HitRate     float64   `json:"hit_rate"`
	TotalHits   int64     `json:"total_hits"`
	TotalMisses int64     `json:"total_misses"`
	PurgedAt    *time.Time `json:"purged_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Request de criação de site
type CreateSiteRequest struct {
	Domain       string     `json:"domain" validate:"required"`
	CustomDomain string     `json:"custom_domain,omitempty"`
	Config       SiteConfig `json:"config"`
	SSLEnabled   bool       `json:"ssl_enabled"`
	CDNEnabled   bool       `json:"cdn_enabled"`
}

// Request de atualização de site
type UpdateSiteRequest struct {
	Config       *SiteConfig `json:"config,omitempty"`
	CustomDomain *string     `json:"custom_domain,omitempty"`
	SSLEnabled   *bool       `json:"ssl_enabled,omitempty"`
	CDNEnabled   *bool       `json:"cdn_enabled,omitempty"`
}

// Response de estatísticas
type StatsResponse struct {
	Site    *Site         `json:"site"`
	Traffic *TrafficStats `json:"traffic"`
	Cache   *CacheStatus  `json:"cache"`
}
