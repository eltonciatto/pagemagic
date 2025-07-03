package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Models
type User struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email            string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash     string     `json:"-" gorm:"not null"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Avatar           string     `json:"avatar"`
	EmailVerified    bool       `json:"email_verified" gorm:"default:false"`
	TwoFactorSecret  string     `json:"-"`
	TwoFactorEnabled bool       `json:"two_factor_enabled" gorm:"default:false"`
	Status           string     `json:"status" gorm:"default:active"` // active, suspended, deleted
	LastLoginAt      *time.Time `json:"last_login_at"`
	LastLoginIP      string     `json:"last_login_ip"`
	FailedAttempts   int        `json:"-" gorm:"default:0"`
	LockedUntil      *time.Time `json:"-"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" gorm:"index"`

	// Relationships
	Roles     []Role     `json:"roles" gorm:"many2many:user_roles;"`
	Sessions  []Session  `json:"-"`
	AuditLogs []AuditLog `json:"-"`
	Devices   []Device   `json:"-"`
}

type Role struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Permission struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	Resource    string    `json:"resource" gorm:"not null"`
	Action      string    `json:"action" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Session struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	DeviceID     uuid.UUID  `json:"device_id" gorm:"type:uuid;index"`
	AccessToken  string     `json:"-" gorm:"uniqueIndex;not null"`
	RefreshToken string     `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt    time.Time  `json:"expires_at"`
	IPAddress    string     `json:"ip_address"`
	UserAgent    string     `json:"user_agent"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastUsedAt   time.Time  `json:"last_used_at"`
	CreatedAt    time.Time  `json:"created_at"`
	RevokedAt    *time.Time `json:"revoked_at"`
}

type Device struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID           uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`     // web, mobile, desktop
	Platform         string    `json:"platform"` // ios, android, windows, macos, linux
	Fingerprint      string    `json:"fingerprint" gorm:"uniqueIndex;not null"`
	IsTrusted        bool      `json:"is_trusted" gorm:"default:false"`
	LastSeenAt       time.Time `json:"last_seen_at"`
	LastSeenIP       string    `json:"last_seen_ip"`
	LastSeenLocation string    `json:"last_seen_location"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	SessionID  *uuid.UUID `json:"session_id" gorm:"type:uuid;index"`
	Action     string     `json:"action" gorm:"not null;index"`
	Resource   string     `json:"resource" gorm:"index"`
	ResourceID *uuid.UUID `json:"resource_id" gorm:"type:uuid;index"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	Metadata   string     `json:"metadata" gorm:"type:jsonb"`
	Result     string     `json:"result"` // success, failure, error
	CreatedAt  time.Time  `json:"created_at" gorm:"index"`
}

type OAuthProvider struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name         string    `json:"name" gorm:"uniqueIndex;not null"`
	ClientID     string    `json:"client_id" gorm:"not null"`
	ClientSecret string    `json:"-"`
	AuthURL      string    `json:"auth_url"`
	TokenURL     string    `json:"token_url"`
	UserInfoURL  string    `json:"user_info_url"`
	Scopes       string    `json:"scopes"`
	IsEnabled    bool      `json:"is_enabled" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Service
type AuthService struct {
	db           *gorm.DB
	redis        *redis.Client
	enforcer     *casbin.Enforcer
	jwtSecret    []byte
	oauthConfigs map[string]*oauth2.Config
}

// Metrics
var (
	authRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pagemagic_auth_requests_total",
			Help: "Total number of authentication requests",
		},
		[]string{"method", "result"},
	)

	authDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "pagemagic_auth_duration_seconds",
			Help: "Authentication request duration",
		},
		[]string{"method"},
	)

	activeSessions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "pagemagic_active_sessions",
			Help: "Number of active user sessions",
		},
	)
)

func init() {
	prometheus.MustRegister(authRequests)
	prometheus.MustRegister(authDuration)
	prometheus.MustRegister(activeSessions)
}

// Initialize service
func NewAuthService() *AuthService {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=pagemagic_auth port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&User{}, &Role{}, &Permission{}, &Session{}, &Device{}, &AuditLog{}, &OAuthProvider{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Redis connection
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: strings.TrimPrefix(redisURL, "redis://"),
	})

	// Casbin enforcer
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatal("Failed to create casbin adapter:", err)
	}

	enforcer, err := casbin.NewEnforcer("rbac_model.conf", adapter)
	if err != nil {
		log.Fatal("Failed to create casbin enforcer:", err)
	}

	// JWT secret
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("your-super-secret-jwt-key-change-in-production")
	}

	// OAuth configs
	oauthConfigs := map[string]*oauth2.Config{
		"google": {
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("BASE_URL") + "/api/v1/auth/oauth/google/callback",
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		},
		"github": {
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("BASE_URL") + "/api/v1/auth/oauth/github/callback",
			Scopes:       []string{"user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		},
	}

	service := &AuthService{
		db:           db,
		redis:        rdb,
		enforcer:     enforcer,
		jwtSecret:    jwtSecret,
		oauthConfigs: oauthConfigs,
	}

	// Initialize default roles and permissions
	service.initializeDefaultData()

	return service
}

func (s *AuthService) initializeDefaultData() {
	// Create default permissions
	permissions := []Permission{
		{Name: "sites.create", Resource: "sites", Action: "create", Description: "Create new sites"},
		{Name: "sites.read", Resource: "sites", Action: "read", Description: "View sites"},
		{Name: "sites.update", Resource: "sites", Action: "update", Description: "Edit sites"},
		{Name: "sites.delete", Resource: "sites", Action: "delete", Description: "Delete sites"},
		{Name: "billing.read", Resource: "billing", Action: "read", Description: "View billing information"},
		{Name: "billing.manage", Resource: "billing", Action: "manage", Description: "Manage billing and subscriptions"},
		{Name: "admin.users", Resource: "admin", Action: "users", Description: "Manage users"},
		{Name: "admin.system", Resource: "admin", Action: "system", Description: "System administration"},
	}

	for _, perm := range permissions {
		s.db.FirstOrCreate(&perm, Permission{Name: perm.Name})
	}

	// Create default roles
	roles := []struct {
		Name        string
		Description string
		Permissions []string
	}{
		{
			Name:        "user",
			Description: "Regular user with basic permissions",
			Permissions: []string{"sites.create", "sites.read", "sites.update", "sites.delete", "billing.read"},
		},
		{
			Name:        "pro",
			Description: "Pro user with enhanced features",
			Permissions: []string{"sites.create", "sites.read", "sites.update", "sites.delete", "billing.read", "billing.manage"},
		},
		{
			Name:        "admin",
			Description: "Administrator with full access",
			Permissions: []string{"sites.create", "sites.read", "sites.update", "sites.delete", "billing.read", "billing.manage", "admin.users", "admin.system"},
		},
	}

	for _, roleData := range roles {
		var role Role
		s.db.FirstOrCreate(&role, Role{Name: roleData.Name})
		role.Description = roleData.Description

		// Add permissions to role
		var permissions []Permission
		s.db.Where("name IN ?", roleData.Permissions).Find(&permissions)
		s.db.Model(&role).Association("Permissions").Replace(permissions)
		s.db.Save(&role)
	}
}

// JWT Claims
type JWTClaims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	SessionID string   `json:"session_id"`
	jwt.RegisteredClaims
}

// Authentication handlers
func (s *AuthService) Register(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		s.logAudit(nil, nil, "auth.register", "user", nil, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{"error": err.Error()}, "failure")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		authRequests.WithLabelValues("register", "error").Inc()
		return
	}

	timer := prometheus.NewTimer(authDuration.WithLabelValues("register"))
	defer timer.ObserveDuration()

	// Check if user already exists
	var existingUser User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		s.logAudit(nil, nil, "auth.register", "user", nil, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{"email": req.Email, "error": "user_exists"}, "failure")
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		authRequests.WithLabelValues("register", "failure").Inc()
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		authRequests.WithLabelValues("register", "error").Inc()
		return
	}

	// Create user
	user := User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Status:       "active",
	}

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		authRequests.WithLabelValues("register", "error").Inc()
		return
	}

	// Assign default role
	var defaultRole Role
	s.db.Where("name = ?", "user").First(&defaultRole)
	s.db.Model(&user).Association("Roles").Append(&defaultRole)

	// Log audit
	s.logAudit(&user.ID, nil, "auth.register", "user", &user.ID, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{"email": req.Email}, "success")

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(&user, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		authRequests.WithLabelValues("register", "error").Inc()
		return
	}

	authRequests.WithLabelValues("register", "success").Inc()
	activeSessions.Inc()

	c.JSON(http.StatusCreated, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    3600, // 1 hour
	})
}

func (s *AuthService) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		TOTPCode string `json:"totp_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		authRequests.WithLabelValues("login", "error").Inc()
		return
	}

	timer := prometheus.NewTimer(authDuration.WithLabelValues("login"))
	defer timer.ObserveDuration()

	// Find user
	var user User
	if err := s.db.Preload("Roles").Where("email = ?", req.Email).First(&user).Error; err != nil {
		s.logAudit(nil, nil, "auth.login", "user", nil, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{"email": req.Email, "error": "user_not_found"}, "failure")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		authRequests.WithLabelValues("login", "failure").Inc()
		return
	}

	// Check if account is locked
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		s.logAudit(&user.ID, nil, "auth.login", "user", &user.ID, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{"error": "account_locked"}, "failure")
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Account temporarily locked"})
		authRequests.WithLabelValues("login", "failure").Inc()
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Increment failed attempts
		user.FailedAttempts++
		if user.FailedAttempts >= 5 {
			lockUntil := time.Now().Add(15 * time.Minute)
			user.LockedUntil = &lockUntil
		}
		s.db.Save(&user)

		s.logAudit(&user.ID, nil, "auth.login", "user", &user.ID, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{"error": "invalid_password", "failed_attempts": user.FailedAttempts}, "failure")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		authRequests.WithLabelValues("login", "failure").Inc()
		return
	}

	// Check 2FA if enabled
	if user.TwoFactorEnabled {
		if req.TOTPCode == "" {
			c.JSON(http.StatusPreconditionRequired, gin.H{
				"error":        "2FA required",
				"requires_2fa": true,
			})
			return
		}

		if !totp.Validate(req.TOTPCode, user.TwoFactorSecret) {
			s.logAudit(&user.ID, nil, "auth.login", "user", &user.ID, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{"error": "invalid_2fa"}, "failure")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
			authRequests.WithLabelValues("login", "failure").Inc()
			return
		}
	}

	// Reset failed attempts
	user.FailedAttempts = 0
	user.LockedUntil = nil
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = c.ClientIP()
	s.db.Save(&user)

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(&user, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		authRequests.WithLabelValues("login", "error").Inc()
		return
	}

	s.logAudit(&user.ID, nil, "auth.login", "user", &user.ID, c.ClientIP(), c.GetHeader("User-Agent"), map[string]interface{}{"email": req.Email}, "success")
	authRequests.WithLabelValues("login", "success").Inc()
	activeSessions.Inc()

	c.JSON(http.StatusOK, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    3600,
	})
}

func (s *AuthService) generateTokens(user *User, ipAddress, userAgent string) (string, string, error) {
	// Create session
	session := Session{
		UserID:     user.ID,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		IsActive:   true,
		LastUsedAt: time.Now(),
	}

	// Generate access token
	accessTokenBytes := make([]byte, 32)
	rand.Read(accessTokenBytes)
	session.AccessToken = base64.URLEncoding.EncodeToString(accessTokenBytes)

	// Generate refresh token
	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	session.RefreshToken = base64.URLEncoding.EncodeToString(refreshTokenBytes)

	if err := s.db.Create(&session).Error; err != nil {
		return "", "", err
	}

	// Get user roles
	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}

	// Create JWT
	claims := JWTClaims{
		UserID:    user.ID.String(),
		Email:     user.Email,
		Roles:     roleNames,
		SessionID: session.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "pagemagic-auth",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, session.RefreshToken, nil
}

func (s *AuthService) logAudit(userID, sessionID *uuid.UUID, action, resource string, resourceID *uuid.UUID, ipAddress, userAgent string, metadata map[string]interface{}, result string) {
	metadataJSON, _ := json.Marshal(metadata)

	audit := AuditLog{
		UserID:     userID,
		SessionID:  sessionID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Metadata:   string(metadataJSON),
		Result:     result,
	}

	s.db.Create(&audit)
}

func main() {
	authService := NewAuthService()

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "auth-svc"})
	})

	// Auth routes
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authService.Register)
			auth.POST("/login", authService.Login)
			// TODO: Add more endpoints (OAuth, 2FA, etc.)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Auth service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
