package models

import (
	"time"

	"github.com/google/uuid"
)

// AuthProvider tipos de provedores de autenticação
type AuthProvider string

const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
	AuthProviderGitHub AuthProvider = "github"
	AuthProviderApple  AuthProvider = "apple"
)

// UserStatus status do usuário
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// User modelo do usuário
type User struct {
	ID                        uuid.UUID  `json:"id" db:"id"`
	Email                     string     `json:"email" db:"email"`
	EmailVerified             bool       `json:"email_verified" db:"email_verified"`
	PasswordHash              *string    `json:"-" db:"password_hash"`
	FirstName                 *string    `json:"first_name" db:"first_name"`
	LastName                  *string    `json:"last_name" db:"last_name"`
	AvatarURL                 *string    `json:"avatar_url" db:"avatar_url"`
	Status                    UserStatus `json:"status" db:"status"`
	Locale                    string     `json:"locale" db:"locale"`
	Timezone                  string     `json:"timezone" db:"timezone"`
	MobileVerified            bool       `json:"mobile_verified" db:"mobile_verified"`
	MobileNumber              *string    `json:"mobile_number" db:"mobile_number"`
	PushNotificationsEnabled  bool       `json:"push_notifications_enabled" db:"push_notifications_enabled"`
	LastLoginAt               *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt                 time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at" db:"updated_at"`
}

// UserAuthProvider modelo de provedor de autenticação do usuário
type UserAuthProvider struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	UserID         uuid.UUID    `json:"user_id" db:"user_id"`
	Provider       AuthProvider `json:"provider" db:"provider"`
	ProviderUserID string       `json:"provider_user_id" db:"provider_user_id"`
	ProviderEmail  *string      `json:"provider_email" db:"provider_email"`
	ProviderData   *string      `json:"provider_data" db:"provider_data"` // JSON
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
}

// MagicLink modelo de magic link
type MagicLink struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Token     string    `json:"token" db:"token"`
	Used      bool      `json:"used" db:"used"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// RefreshToken modelo de refresh token
type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	Used      bool       `json:"used" db:"used"`
	UsedAt    *time.Time `json:"used_at" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// CreateUserRequest request para criação de usuário
type CreateUserRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=8"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Locale    string  `json:"locale"`
	Timezone  string  `json:"timezone"`
}

// LoginRequest request de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// MagicLinkRequest request de magic link
type MagicLinkRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyMagicLinkRequest request de verificação de magic link
type VerifyMagicLinkRequest struct {
	Token string `json:"token" validate:"required"`
}

// RefreshTokenRequest request de refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateUserRequest request de atualização de usuário
type UpdateUserRequest struct {
	FirstName                *string `json:"first_name"`
	LastName                 *string `json:"last_name"`
	AvatarURL                *string `json:"avatar_url"`
	Locale                   *string `json:"locale"`
	Timezone                 *string `json:"timezone"`
	MobileNumber             *string `json:"mobile_number"`
	PushNotificationsEnabled *bool   `json:"push_notifications_enabled"`
}

// ChangePasswordRequest request de mudança de senha
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// AuthResponse resposta de autenticação
type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// JWTClaims claims do JWT
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Exp    int64     `json:"exp"`
	Iat    int64     `json:"iat"`
	Type   string    `json:"type"` // "access" ou "refresh"
}

// OAuthUserInfo informações do usuário OAuth
type OAuthUserInfo struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	EmailVerified bool    `json:"email_verified"`
	FirstName     *string `json:"first_name"`
	LastName      *string `json:"last_name"`
	AvatarURL     *string `json:"avatar_url"`
	Provider      string  `json:"provider"`
}

// APIResponse resposta padrão da API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *APIMeta    `json:"meta,omitempty"`
}

// APIError erro da API
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// APIMeta metadados da API
type APIMeta struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// FullName retorna o nome completo do usuário
func (u *User) FullName() string {
	if u.FirstName != nil && u.LastName != nil {
		return *u.FirstName + " " + *u.LastName
	}
	if u.FirstName != nil {
		return *u.FirstName
	}
	if u.LastName != nil {
		return *u.LastName
	}
	return u.Email
}

// IsActive verifica se o usuário está ativo
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsExpired verifica se o magic link expirou
func (m *MagicLink) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}

// IsValid verifica se o magic link é válido (não usado e não expirado)
func (m *MagicLink) IsValid() bool {
	return !m.Used && !m.IsExpired()
}

// IsExpired verifica se o refresh token expirou
func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsValid verifica se o refresh token é válido (não usado e não expirado)
func (r *RefreshToken) IsValid() bool {
	return !r.Used && !r.IsExpired()
}
