package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"pagemagic/auth-svc/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// UserRepository interface para repositório de usuários
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
}

// AuthProviderRepository interface para repositório de provedores de auth
type AuthProviderRepository interface {
	Create(ctx context.Context, provider *models.UserAuthProvider) error
	GetByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider models.AuthProvider) (*models.UserAuthProvider, error)
	GetByProviderAndUserID(ctx context.Context, provider models.AuthProvider, providerUserID string) (*models.UserAuthProvider, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserAuthProvider, error)
}

// MagicLinkRepository interface para repositório de magic links
type MagicLinkRepository interface {
	Create(ctx context.Context, link *models.MagicLink) error
	GetByToken(ctx context.Context, token string) (*models.MagicLink, error)
	MarkAsUsed(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	DeleteByEmail(ctx context.Context, email string) error
}

// RefreshTokenRepository interface para repositório de refresh tokens
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	MarkAsUsed(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// Repository agregador de todos os repositórios
type Repository struct {
	User         UserRepository
	AuthProvider AuthProviderRepository
	MagicLink    MagicLinkRepository
	RefreshToken RefreshTokenRepository
}

// PostgresUserRepository implementação PostgreSQL do UserRepository
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository cria uma nova instância do repositório
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create cria um novo usuário
func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			id, email, email_verified, password_hash, first_name, last_name,
			avatar_url, status, locale, timezone, mobile_verified, mobile_number,
			push_notifications_enabled, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.EmailVerified, user.PasswordHash,
		user.FirstName, user.LastName, user.AvatarURL, user.Status,
		user.Locale, user.Timezone, user.MobileVerified, user.MobileNumber,
		user.PushNotificationsEnabled, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID busca usuário por ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, email_verified, password_hash, first_name, last_name,
			   avatar_url, status, locale, timezone, mobile_verified, mobile_number,
			   push_notifications_enabled, last_login_at, created_at, updated_at
		FROM users WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.AvatarURL, &user.Status,
		&user.Locale, &user.Timezone, &user.MobileVerified, &user.MobileNumber,
		&user.PushNotificationsEnabled, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail busca usuário por email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, email_verified, password_hash, first_name, last_name,
			   avatar_url, status, locale, timezone, mobile_verified, mobile_number,
			   push_notifications_enabled, last_login_at, created_at, updated_at
		FROM users WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.AvatarURL, &user.Status,
		&user.Locale, &user.Timezone, &user.MobileVerified, &user.MobileNumber,
		&user.PushNotificationsEnabled, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update atualiza um usuário
func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET
			email = $2, email_verified = $3, password_hash = $4, first_name = $5,
			last_name = $6, avatar_url = $7, status = $8, locale = $9, timezone = $10,
			mobile_verified = $11, mobile_number = $12, push_notifications_enabled = $13,
			updated_at = $14
		WHERE id = $1`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.EmailVerified, user.PasswordHash,
		user.FirstName, user.LastName, user.AvatarURL, user.Status,
		user.Locale, user.Timezone, user.MobileVerified, user.MobileNumber,
		user.PushNotificationsEnabled, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete remove um usuário (soft delete)
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET status = 'deleted', updated_at = $2 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List lista usuários com paginação
func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, email, email_verified, password_hash, first_name, last_name,
			   avatar_url, status, locale, timezone, mobile_verified, mobile_number,
			   push_notifications_enabled, last_login_at, created_at, updated_at
		FROM users 
		WHERE status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.EmailVerified, &user.PasswordHash,
			&user.FirstName, &user.LastName, &user.AvatarURL, &user.Status,
			&user.Locale, &user.Timezone, &user.MobileVerified, &user.MobileNumber,
			&user.PushNotificationsEnabled, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, nil
}

// UpdateLastLogin atualiza o último login do usuário
func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $2, updated_at = $2 WHERE id = $1`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// PostgresMagicLinkRepository implementação PostgreSQL do MagicLinkRepository
type PostgresMagicLinkRepository struct {
	db *sql.DB
}

// NewPostgresMagicLinkRepository cria uma nova instância do repositório
func NewPostgresMagicLinkRepository(db *sql.DB) *PostgresMagicLinkRepository {
	return &PostgresMagicLinkRepository{db: db}
}

// Create cria um novo magic link
func (r *PostgresMagicLinkRepository) Create(ctx context.Context, link *models.MagicLink) error {
	query := `
		INSERT INTO magic_links (id, email, token, used, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		link.ID, link.Email, link.Token, link.Used, link.ExpiresAt, link.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create magic link: %w", err)
	}

	return nil
}

// GetByToken busca magic link por token
func (r *PostgresMagicLinkRepository) GetByToken(ctx context.Context, token string) (*models.MagicLink, error) {
	query := `
		SELECT id, email, token, used, expires_at, created_at
		FROM magic_links WHERE token = $1`

	link := &models.MagicLink{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&link.ID, &link.Email, &link.Token, &link.Used, &link.ExpiresAt, &link.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("magic link not found")
		}
		return nil, fmt.Errorf("failed to get magic link: %w", err)
	}

	return link, nil
}

// MarkAsUsed marca o magic link como usado
func (r *PostgresMagicLinkRepository) MarkAsUsed(ctx context.Context, token string) error {
	query := `UPDATE magic_links SET used = true WHERE token = $1`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to mark magic link as used: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("magic link not found")
	}

	return nil
}

// DeleteExpired remove magic links expirados
func (r *PostgresMagicLinkRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM magic_links WHERE expires_at < $1`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired magic links: %w", err)
	}

	return nil
}

// DeleteByEmail remove magic links por email
func (r *PostgresMagicLinkRepository) DeleteByEmail(ctx context.Context, email string) error {
	query := `DELETE FROM magic_links WHERE email = $1`

	_, err := r.db.ExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("failed to delete magic links by email: %w", err)
	}

	return nil
}

// NewRepository cria uma nova instância do agregador de repositórios
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:      NewPostgresUserRepository(db),
		MagicLink: NewPostgresMagicLinkRepository(db),
		// TODO: Implementar outros repositórios
	}
}
