package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"pagemagic/auth-svc/internal/config"
	"pagemagic/auth-svc/internal/models"
	"pagemagic/auth-svc/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    repository.UserRepository
	magicRepo   repository.MagicLinkRepository
	refreshRepo repository.RefreshTokenRepository
	config      *config.Config
}

func NewAuthService(repo *repository.Repository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo:    repo.User,
		magicRepo:   repo.MagicLink,
		refreshRepo: repo.RefreshToken,
		config:      config,
	}
}

func (s *AuthService) SendMagicLink(ctx context.Context, email string) (*models.MagicLink, error) {
	// Gerar token único
	token, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Criar magic link
	magicLink := &models.MagicLink{
		ID:        uuid.New(),
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute), // 15 minutos
		Used:      false,
		CreatedAt: time.Now(),
	}

	// Salvar no banco
	if err := s.magicRepo.Create(ctx, magicLink); err != nil {
		return nil, fmt.Errorf("failed to create magic link: %w", err)
	}

	// TODO: Enviar email com o link
	// Por enquanto apenas retorna o magic link
	return magicLink, nil
}

func (s *AuthService) VerifyMagicLink(ctx context.Context, token string) (*models.User, string, string, error) {
	// Buscar magic link
	magicLink, err := s.magicRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, "", "", fmt.Errorf("invalid token: %w", err)
	}

	// Verificar se não expirou
	if time.Now().After(magicLink.ExpiresAt) {
		return nil, "", "", fmt.Errorf("token expired")
	}

	// Verificar se não foi usado
	if magicLink.Used {
		return nil, "", "", fmt.Errorf("token already used")
	}

	// Marcar como usado
	if err := s.magicRepo.MarkAsUsed(ctx, token); err != nil {
		return nil, "", "", fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Buscar ou criar usuário
	user, err := s.userRepo.GetByEmail(ctx, magicLink.Email)
	if err != nil {
		// Usuário não existe, criar novo
		user = &models.User{
			ID:        uuid.New(),
			Email:     magicLink.Email,
			Status:    models.UserStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, "", "", fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Gerar tokens JWT
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Verificar refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	// Verificar se é refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", fmt.Errorf("not a refresh token")
	}

	// Buscar usuário
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Gerar novo access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verificar se é access token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, fmt.Errorf("not an access token")
	}

	// Buscar usuário
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (s *AuthService) generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"type":    "access",
		"exp":     time.Now().Add(s.config.JWT.AccessTokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *AuthService) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"type":    "refresh",
		"exp":     time.Now().Add(s.config.JWT.RefreshTokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.RefreshTokenSecret))
}
