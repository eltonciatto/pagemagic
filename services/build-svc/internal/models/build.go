package models

import (
	"time"
)

// BuildJob representa um job de build
type BuildJob struct {
	ID          string                 `json:"id" db:"id"`
	SiteID      string                 `json:"site_id" db:"site_id"`
	UserID      string                 `json:"user_id" db:"user_id"`
	Type        string                 `json:"type" db:"type"`               // "full", "incremental"
	Status      string                 `json:"status" db:"status"`           // "pending", "building", "completed", "failed"
	SourceType  string                 `json:"source_type" db:"source_type"` // "visual_editor", "code", "template"
	SourceData  map[string]interface{} `json:"source_data" db:"source_data"`
	BuildConfig BuildConfig            `json:"build_config" db:"build_config"`
	Output      BuildOutput            `json:"output" db:"output"`
	Logs        []BuildLog             `json:"logs" db:"logs"`
	Error       string                 `json:"error" db:"error"`
	StartedAt   time.Time              `json:"started_at" db:"started_at"`
	EndedAt     *time.Time             `json:"ended_at" db:"ended_at"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// BuildConfig contém configurações de build
type BuildConfig struct {
	Framework    string            `json:"framework"` // "static", "react", "vue", "angular", "next"
	Theme        string            `json:"theme"`
	Plugins      []string          `json:"plugins"`
	Environment  map[string]string `json:"environment"`
	Optimization BuildOptimization `json:"optimization"`
	SEO          SEOConfig         `json:"seo"`
	PWA          PWAConfig         `json:"pwa"`
	Performance  PerformanceConfig `json:"performance"`
}

// BuildOptimization configurações de otimização
type BuildOptimization struct {
	MinifyHTML     bool `json:"minify_html"`
	MinifyCSS      bool `json:"minify_css"`
	MinifyJS       bool `json:"minify_js"`
	OptimizeImages bool `json:"optimize_images"`
	GenerateWebP   bool `json:"generate_webp"`
	LazyLoading    bool `json:"lazy_loading"`
	TreeShaking    bool `json:"tree_shaking"`
	CodeSplitting  bool `json:"code_splitting"`
}

// SEOConfig configurações de SEO
type SEOConfig struct {
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Keywords       []string          `json:"keywords"`
	Author         string            `json:"author"`
	Language       string            `json:"language"`
	Canonical      string            `json:"canonical"`
	MetaTags       map[string]string `json:"meta_tags"`
	OpenGraph      OpenGraphConfig   `json:"open_graph"`
	TwitterCard    TwitterCardConfig `json:"twitter_card"`
	StructuredData interface{}       `json:"structured_data"`
	Sitemap        bool              `json:"sitemap"`
	RobotsTxt      string            `json:"robots_txt"`
}

// OpenGraphConfig configurações do Open Graph
type OpenGraphConfig struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	SiteName    string `json:"site_name"`
}

// TwitterCardConfig configurações do Twitter Card
type TwitterCardConfig struct {
	Card        string `json:"card"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Creator     string `json:"creator"`
}

// PWAConfig configurações de PWA
type PWAConfig struct {
	Enabled         bool      `json:"enabled"`
	Name            string    `json:"name"`
	ShortName       string    `json:"short_name"`
	Description     string    `json:"description"`
	ThemeColor      string    `json:"theme_color"`
	BackgroundColor string    `json:"background_color"`
	Icons           []PWAIcon `json:"icons"`
	StartURL        string    `json:"start_url"`
	Display         string    `json:"display"`
	Orientation     string    `json:"orientation"`
	ServiceWorker   bool      `json:"service_worker"`
	OfflineSupport  bool      `json:"offline_support"`
}

// PWAIcon ícone do PWA
type PWAIcon struct {
	Src     string `json:"src"`
	Sizes   string `json:"sizes"`
	Type    string `json:"type"`
	Purpose string `json:"purpose"`
}

// PerformanceConfig configurações de performance
type PerformanceConfig struct {
	CriticalCSS    bool `json:"critical_css"`
	PreloadFonts   bool `json:"preload_fonts"`
	PreloadImages  bool `json:"preload_images"`
	ResourceHints  bool `json:"resource_hints"`
	Compression    bool `json:"compression"`
	HTTP2Push      bool `json:"http2_push"`
	ServiceWorker  bool `json:"service_worker"`
	InlineSmallCSS bool `json:"inline_small_css"`
	InlineSmallJS  bool `json:"inline_small_js"`
}

// BuildOutput resultado do build
type BuildOutput struct {
	Files         []BuildFile `json:"files"`
	Assets        []BuildFile `json:"assets"`
	Pages         []PageInfo  `json:"pages"`
	Stats         BuildStats  `json:"stats"`
	Manifest      interface{} `json:"manifest"`
	Sitemap       string      `json:"sitemap"`
	RobotsTxt     string      `json:"robots_txt"`
	ServiceWorker string      `json:"service_worker"`
}

// BuildFile arquivo gerado no build
type BuildFile struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Hash        string `json:"hash"`
	Compressed  bool   `json:"compressed"`
	Optimized   bool   `json:"optimized"`
}

// PageInfo informações da página
type PageInfo struct {
	Path        string            `json:"path"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Keywords    []string          `json:"keywords"`
	MetaTags    map[string]string `json:"meta_tags"`
	Language    string            `json:"language"`
	LastMod     time.Time         `json:"last_mod"`
}

// BuildStats estatísticas do build
type BuildStats struct {
	TotalFiles     int                `json:"total_files"`
	TotalSize      int64              `json:"total_size"`
	CompressedSize int64              `json:"compressed_size"`
	OptimizedFiles int                `json:"optimized_files"`
	Duration       time.Duration      `json:"duration"`
	Performance    PerformanceMetrics `json:"performance"`
}

// PerformanceMetrics métricas de performance
type PerformanceMetrics struct {
	LighthouseScore   int     `json:"lighthouse_score"`
	PageSpeedScore    int     `json:"pagespeed_score"`
	FirstPaint        float64 `json:"first_paint"`
	FirstContentful   float64 `json:"first_contentful"`
	LargestContentful float64 `json:"largest_contentful"`
	CumulativeLayout  float64 `json:"cumulative_layout"`
	TimeToInteractive float64 `json:"time_to_interactive"`
}

// BuildLog log do processo de build
type BuildLog struct {
	Level     string      `json:"level"` // "info", "warn", "error", "debug"
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

// Template configurações de template
type Template struct {
	ID          string             `json:"id" db:"id"`
	Name        string             `json:"name" db:"name"`
	Description string             `json:"description" db:"description"`
	Category    string             `json:"category" db:"category"`
	Framework   string             `json:"framework" db:"framework"`
	Preview     string             `json:"preview" db:"preview"`
	Thumbnail   string             `json:"thumbnail" db:"thumbnail"`
	Config      BuildConfig        `json:"config" db:"config"`
	Files       []TemplateFile     `json:"files" db:"files"`
	Variables   []TemplateVariable `json:"variables" db:"variables"`
	Tags        []string           `json:"tags" db:"tags"`
	Version     string             `json:"version" db:"version"`
	Author      string             `json:"author" db:"author"`
	License     string             `json:"license" db:"license"`
	Featured    bool               `json:"featured" db:"featured"`
	Downloads   int                `json:"downloads" db:"downloads"`
	Rating      float64            `json:"rating" db:"rating"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// TemplateFile arquivo do template
type TemplateFile struct {
	Path        string `json:"path"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Replaceable bool   `json:"replaceable"`
	Required    bool   `json:"required"`
}

// TemplateVariable variável do template
type TemplateVariable struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"` // "string", "number", "boolean", "color", "image", "select"
	Label       string        `json:"label"`
	Description string        `json:"description"`
	Default     interface{}   `json:"default"`
	Required    bool          `json:"required"`
	Options     []interface{} `json:"options,omitempty"`
	Validation  interface{}   `json:"validation,omitempty"`
}

// Requests
type CreateBuildJobRequest struct {
	SiteID      string                 `json:"site_id" validate:"required"`
	Type        string                 `json:"type" validate:"required"`
	SourceType  string                 `json:"source_type" validate:"required"`
	SourceData  map[string]interface{} `json:"source_data" validate:"required"`
	BuildConfig BuildConfig            `json:"build_config"`
}

type UpdateBuildJobRequest struct {
	Status      *string      `json:"status,omitempty"`
	BuildConfig *BuildConfig `json:"build_config,omitempty"`
	Error       *string      `json:"error,omitempty"`
}

// Responses
type BuildJobResponse struct {
	Job    *BuildJob    `json:"job"`
	Output *BuildOutput `json:"output,omitempty"`
}

type TemplateListResponse struct {
	Templates []Template `json:"templates"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}
