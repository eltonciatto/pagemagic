package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Docker   DockerConfig
	Build    BuildConfig
	Storage  StorageConfig
	NATS     NATSConfig
	Redis    RedisConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	URL         string
	MaxConns    int
	MinConns    int
	MaxIdleTime time.Duration
}

type DockerConfig struct {
	Host           string
	APIVersion     string
	RegistryURL    string
	RegistryUser   string
	RegistryPass   string
	BuildxInstance string
}

type BuildConfig struct {
	WorkDir          string
	MaxConcurrent    int
	Timeout          time.Duration
	BaseImage        string
	NodeVersion      string
	NginxVersion     string
	TurbopackEnabled bool
}

type StorageConfig struct {
	Provider    string
	S3Bucket    string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	LocalPath   string
}

type NATSConfig struct {
	URL      string
	User     string
	Password string
	Timeout  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	TTL      time.Duration
}

type LoggingConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("BUILD_SVC_PORT", "8083"),
			Host:         getEnv("BUILD_SVC_HOST", "0.0.0.0"),
			Environment:  getEnv("NODE_ENV", "development"),
			ReadTimeout:  parseDuration(getEnv("SERVER_READ_TIMEOUT", "30s")),
			WriteTimeout: parseDuration(getEnv("SERVER_WRITE_TIMEOUT", "30s")),
		},
		Database: DatabaseConfig{
			URL:         getEnv("DATABASE_URL", "postgres://pagemagic:password@localhost:5432/pagemagic"),
			MaxConns:    parseInt(getEnv("DB_MAX_CONNS", "25")),
			MinConns:    parseInt(getEnv("DB_MIN_CONNS", "5")),
			MaxIdleTime: parseDuration(getEnv("DB_MAX_IDLE_TIME", "15m")),
		},
		Docker: DockerConfig{
			Host:           getEnv("DOCKER_HOST", "unix:///var/run/docker.sock"),
			APIVersion:     getEnv("DOCKER_API_VERSION", "1.41"),
			RegistryURL:    getEnv("DOCKER_REGISTRY_URL", "registry.pagemagic.io"),
			RegistryUser:   getEnv("DOCKER_REGISTRY_USER", ""),
			RegistryPass:   getEnv("DOCKER_REGISTRY_PASS", ""),
			BuildxInstance: getEnv("BUILDX_INSTANCE", "pagemagic-builder"),
		},
		Build: BuildConfig{
			WorkDir:          getEnv("BUILD_WORK_DIR", "/tmp/pagemagic/builds"),
			MaxConcurrent:    parseInt(getEnv("BUILD_MAX_CONCURRENT", "5")),
			Timeout:          parseDuration(getEnv("BUILD_TIMEOUT", "10m")),
			BaseImage:        getEnv("BUILD_BASE_IMAGE", "node:18-alpine"),
			NodeVersion:      getEnv("BUILD_NODE_VERSION", "18"),
			NginxVersion:     getEnv("BUILD_NGINX_VERSION", "1.25-alpine"),
			TurbopackEnabled: parseBool(getEnv("BUILD_TURBOPACK_ENABLED", "true")),
		},
		Storage: StorageConfig{
			Provider:    getEnv("STORAGE_PROVIDER", "s3"),
			S3Bucket:    getEnv("S3_BUCKET", "pagemagic-builds"),
			S3Region:    getEnv("S3_REGION", "us-east-1"),
			S3AccessKey: getEnv("S3_ACCESS_KEY", ""),
			S3SecretKey: getEnv("S3_SECRET_KEY", ""),
			LocalPath:   getEnv("STORAGE_LOCAL_PATH", "/app/storage"),
		},
		NATS: NATSConfig{
			URL:      getEnv("NATS_URL", "nats://localhost:4222"),
			User:     getEnv("NATS_USER", ""),
			Password: getEnv("NATS_PASSWORD", ""),
			Timeout:  parseDuration(getEnv("NATS_TIMEOUT", "30s")),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       parseInt(getEnv("REDIS_DB", "1")),
			TTL:      parseDuration(getEnv("REDIS_TTL", "24h")),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.Environment == "production" {
		if c.Docker.RegistryUser == "" || c.Docker.RegistryPass == "" {
			return fmt.Errorf("docker registry credentials required in production")
		}
		if c.Storage.Provider == "s3" && (c.Storage.S3AccessKey == "" || c.Storage.S3SecretKey == "") {
			return fmt.Errorf("S3 credentials required when using S3 storage")
		}
	}
	return nil
}

func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}

func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}
