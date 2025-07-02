package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	CDN      CDNConfig
	Auth     AuthConfig
	SSL      SSLConfig
	Cache    CacheConfig
	Metrics  MetricsConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type StorageConfig struct {
	Type      string // "s3", "gcs", "local"
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
	LocalPath string
}

type CDNConfig struct {
	Provider  string // "cloudflare", "aws", "gcp"
	Zone      string
	APIKey    string
	APISecret string
	BaseURL   string
}

type AuthConfig struct {
	JWTSecret   string
	ServiceName string
}

type SSLConfig struct {
	CertPath  string
	KeyPath   string
	AutoSSL   bool
	ACMEEmail string
	ACMEDir   string
}

type CacheConfig struct {
	RedisURL  string
	TTL       time.Duration
	MaxMemory string
}

type MetricsConfig struct {
	Enabled bool
	Port    string
	Path    string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("HOST_PORT", "8080"),
			Host:         getEnv("HOST_BIND", "0.0.0.0"),
			ReadTimeout:  parseDuration(getEnv("HOST_READ_TIMEOUT", "30s")),
			WriteTimeout: parseDuration(getEnv("HOST_WRITE_TIMEOUT", "30s")),
			IdleTimeout:  parseDuration(getEnv("HOST_IDLE_TIMEOUT", "120s")),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "pagemagic"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Storage: StorageConfig{
			Type:      getEnv("STORAGE_TYPE", "s3"),
			Bucket:    getEnv("STORAGE_BUCKET", "pagemagic-sites"),
			Region:    getEnv("STORAGE_REGION", "us-east-1"),
			AccessKey: getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			LocalPath: getEnv("STORAGE_LOCAL_PATH", "/data/sites"),
		},
		CDN: CDNConfig{
			Provider:  getEnv("CDN_PROVIDER", "cloudflare"),
			Zone:      getEnv("CDN_ZONE", ""),
			APIKey:    getEnv("CDN_API_KEY", ""),
			APISecret: getEnv("CDN_API_SECRET", ""),
			BaseURL:   getEnv("CDN_BASE_URL", ""),
		},
		Auth: AuthConfig{
			JWTSecret:   getEnv("JWT_SECRET", ""),
			ServiceName: "host-svc",
		},
		SSL: SSLConfig{
			CertPath:  getEnv("SSL_CERT_PATH", "/etc/ssl/certs"),
			KeyPath:   getEnv("SSL_KEY_PATH", "/etc/ssl/private"),
			AutoSSL:   getBool(getEnv("SSL_AUTO", "true")),
			ACMEEmail: getEnv("ACME_EMAIL", ""),
			ACMEDir:   getEnv("ACME_DIR", "/etc/ssl/acme"),
		},
		Cache: CacheConfig{
			RedisURL:  getEnv("REDIS_URL", "redis://localhost:6379"),
			TTL:       parseDuration(getEnv("CACHE_TTL", "24h")),
			MaxMemory: getEnv("CACHE_MAX_MEMORY", "256mb"),
		},
		Metrics: MetricsConfig{
			Enabled: getBool(getEnv("METRICS_ENABLED", "true")),
			Port:    getEnv("METRICS_PORT", "9090"),
			Path:    getEnv("METRICS_PATH", "/metrics"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBool(value string) bool {
	result, _ := strconv.ParseBool(value)
	return result
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 30 * time.Second
	}
	return duration
}
