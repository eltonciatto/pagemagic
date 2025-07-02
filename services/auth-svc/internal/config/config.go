package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config contém todas as configurações da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	NATS     NATSConfig
	Email    EmailConfig
	OAuth    OAuthConfig
	Logging  LoggingConfig
}

// ServerConfig configurações do servidor HTTP
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Environment  string
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	MaxConns int
	MinConns int
}

// JWTConfig configurações JWT
type JWTConfig struct {
	Secret             string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	RefreshTokenSecret string
}

// RedisConfig configurações Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	TTL      time.Duration
}

// NATSConfig configurações NATS
type NATSConfig struct {
	URL      string
	User     string
	Password string
	Timeout  time.Duration
}

// EmailConfig configurações de email
type EmailConfig struct {
	Provider    string
	SendGridKey string
	FromEmail   string
	FromName    string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
}

// OAuthConfig configurações OAuth
type OAuthConfig struct {
	Google GoogleOAuthConfig
	GitHub GitHubOAuthConfig
	Apple  AppleOAuthConfig
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type GitHubOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type AppleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	TeamID       string
	KeyID        string
	PrivateKey   string
}

// LoggingConfig configurações de logging
type LoggingConfig struct {
	Level  string
	Format string
}

// Load carrega as configurações das variáveis de ambiente
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "0.0.0.0"),
			ReadTimeout:  parseDuration(getEnv("SERVER_READ_TIMEOUT", "10s")),
			WriteTimeout: parseDuration(getEnv("SERVER_WRITE_TIMEOUT", "10s")),
			IdleTimeout:  parseDuration(getEnv("SERVER_IDLE_TIMEOUT", "60s")),
			Environment:  getEnv("NODE_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			Name:     getEnv("POSTGRES_DB", "pagemagic"),
			User:     getEnv("POSTGRES_USER", "pagemagic"),
			Password: getEnv("POSTGRES_PASSWORD", "password"),
			SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
			MaxConns: parseInt(getEnv("POSTGRES_MAX_CONNS", "25")),
			MinConns: parseInt(getEnv("POSTGRES_MIN_CONNS", "5")),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenTTL:     parseDuration(getEnv("JWT_ACCESS_TTL", "15m")),
			RefreshTokenTTL:    parseDuration(getEnv("JWT_REFRESH_TTL", "7d")),
			RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "your-refresh-secret"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       parseInt(getEnv("REDIS_DB", "0")),
			TTL:      parseDuration(getEnv("REDIS_TTL", "24h")),
		},
		NATS: NATSConfig{
			URL:      getEnv("NATS_URL", "nats://localhost:4222"),
			User:     getEnv("NATS_USER", ""),
			Password: getEnv("NATS_PASSWORD", ""),
			Timeout:  parseDuration(getEnv("NATS_TIMEOUT", "30s")),
		},
		Email: EmailConfig{
			Provider:    getEnv("EMAIL_PROVIDER", "sendgrid"),
			SendGridKey: getEnv("SENDGRID_API_KEY", ""),
			FromEmail:   getEnv("FROM_EMAIL", "noreply@pagemagic.io"),
			FromName:    getEnv("FROM_NAME", "Page Magic"),
			SMTPHost:    getEnv("SMTP_HOST", ""),
			SMTPPort:    getEnv("SMTP_PORT", "587"),
			SMTPUser:    getEnv("SMTP_USER", ""),
			SMTPPass:    getEnv("SMTP_PASSWORD", ""),
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
			},
			GitHub: GitHubOAuthConfig{
				ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
				ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GITHUB_REDIRECT_URL", ""),
			},
			Apple: AppleOAuthConfig{
				ClientID:     getEnv("APPLE_CLIENT_ID", ""),
				ClientSecret: getEnv("APPLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("APPLE_REDIRECT_URL", ""),
				TeamID:       getEnv("APPLE_TEAM_ID", ""),
				KeyID:        getEnv("APPLE_KEY_ID", ""),
				PrivateKey:   getEnv("APPLE_PRIVATE_KEY", ""),
			},
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// Validar configurações críticas
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate valida se as configurações obrigatórias estão presentes
func (c *Config) Validate() error {
	if c.JWT.Secret == "your-secret-key" && c.Server.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}

	if c.Database.Password == "password" && c.Server.Environment == "production" {
		return fmt.Errorf("POSTGRES_PASSWORD must be set in production")
	}

	return nil
}

// DatabaseURL retorna a URL de conexão do banco de dados
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// RedisAddr retorna o endereço Redis
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// IsProduction verifica se está em ambiente de produção
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// IsDevelopment verifica se está em ambiente de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseInt converte string para int
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// parseDuration converte string para time.Duration
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}
