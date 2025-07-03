package models

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

type User struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	Email                    string     `json:"email" db:"email"`
	EmailVerified            bool       `json:"email_verified" db:"email_verified"`
	FirstName                *string    `json:"first_name,omitempty" db:"first_name"`
	LastName                 *string    `json:"last_name,omitempty" db:"last_name"`
	AvatarURL                *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Status                   UserStatus `json:"status" db:"status"`
	Locale                   string     `json:"locale" db:"locale"`
	Timezone                 string     `json:"timezone" db:"timezone"`
	MobileVerified           *bool      `json:"mobile_verified,omitempty" db:"mobile_verified"`
	MobileNumber             *string    `json:"mobile_number,omitempty" db:"mobile_number"`
	PushNotificationsEnabled *bool      `json:"push_notifications_enabled,omitempty" db:"push_notifications_enabled"`
	LastLoginAt              *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
}

type MagicLink struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	Used      bool       `json:"used" db:"used"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
}

type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	Revoked   bool       `json:"revoked" db:"revoked"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

type Session struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	IPAddress *string   `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent *string   `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Table definitions for database
const (
	UsersTable         = "users"
	MagicLinksTable    = "magic_links"
	RefreshTokensTable = "refresh_tokens"
	SessionsTable      = "sessions"
)
