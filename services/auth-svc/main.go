package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Models
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:'user'"` // user, admin
	Verified  bool      `json:"verified" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type MagicLink struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"not null"`
	Token     string    `json:"token" gorm:"unique;not null"`
	Used      bool      `json:"used" gorm:"default:false"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Request/Response models
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type MagicLinkRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// AuthService handles authentication
type AuthService struct {
	db    *gorm.DB
	redis *redis.Client
	jwt   *JWTService
}

type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

func (j *JWTService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return &claims, nil
}

func NewAuthService() *AuthService {
	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "pagemagic"),
		getEnv("DB_PASSWORD", "pagemagic_dev"),
		getEnv("DB_NAME", "pagemagic"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate
	err = db.AutoMigrate(&User{}, &Session{}, &MagicLink{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// JWT service
	jwtSecret := getEnv("JWT_SECRET", "default_secret")
	jwt := NewJWTService(jwtSecret)

	return &AuthService{
		db:    db,
		redis: rdb,
		jwt:   jwt,
	}
}

func (as *AuthService) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var existingUser User
	if err := as.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := User{
		ID:        generateID(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := as.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate token
	token, err := as.jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save session
	session := Session{
		ID:        generateID(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
	as.db.Create(&session)

	c.JSON(http.StatusCreated, AuthResponse{
		User:  user,
		Token: token,
	})
}

func (as *AuthService) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user User
	if err := as.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := as.jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save session
	session := Session{
		ID:        generateID(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
	as.db.Create(&session)

	c.JSON(http.StatusOK, AuthResponse{
		User:  user,
		Token: token,
	})
}

func (as *AuthService) SendMagicLink(c *gin.Context) {
	var req MagicLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate magic link token
	token := generateToken()

	magicLink := MagicLink{
		ID:        generateID(),
		Email:     req.Email,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := as.db.Create(&magicLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create magic link"})
		return
	}

	// TODO: Send email with magic link
	log.Printf("Magic link for %s: http://localhost:3000/auth/magic?token=%s", req.Email, token)

	c.JSON(http.StatusOK, gin.H{"message": "Magic link sent"})
}

func (as *AuthService) VerifyMagicLink(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	// Find magic link
	var magicLink MagicLink
	if err := as.db.Where("token = ? AND used = false AND expires_at > ?", token, time.Now()).First(&magicLink).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Mark as used
	as.db.Model(&magicLink).Update("used", true)

	// Find or create user
	var user User
	if err := as.db.Where("email = ?", magicLink.Email).First(&user).Error; err != nil {
		// Create new user
		user = User{
			ID:        generateID(),
			Email:     magicLink.Email,
			Name:      strings.Split(magicLink.Email, "@")[0],
			Password:  generateToken(), // Random password
			Verified:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		as.db.Create(&user)
	}

	// Generate JWT token
	jwtToken, err := as.jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		User:  user,
		Token: jwtToken,
	})
}

func (as *AuthService) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Remove "Bearer " prefix
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	claims, err := as.jwt.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := (*claims)["user_id"].(string)

	// Get user details
	var user User
	if err := as.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user":  user,
	})
}

func (as *AuthService) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" && strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
		// Invalidate session
		as.db.Where("token = ?", token).Delete(&Session{})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// Health check
func (as *AuthService) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "auth-svc",
	})
}

func generateID() string {
	return fmt.Sprintf("usr_%d", time.Now().UnixNano())
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func authMiddleware(as *AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}

		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		claims, err := as.jwt.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userID := (*claims)["user_id"].(string)
		c.Header("X-User-ID", userID)
		c.Next()
	}
}

func main() {
	// Initialize service
	authService := NewAuthService()

	// Setup Gin
	r := gin.Default()

	// Middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	v1 := r.Group("/v1")
	{
		v1.GET("/health", authService.Health)
		v1.POST("/register", authService.Register)
		v1.POST("/login", authService.Login)
		v1.POST("/magic-link", authService.SendMagicLink)
		v1.GET("/magic-link/verify", authService.VerifyMagicLink)
		v1.POST("/validate", authService.ValidateToken)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(authMiddleware(authService))
		{
			protected.POST("/logout", authService.Logout)
		}
	}

	port := getEnv("PORT", "8080")
	log.Printf("Auth service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
