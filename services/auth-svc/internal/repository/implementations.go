package repository

import (
	"context"
	"database/sql"
	"fmt"

	"pagemagic/auth-svc/internal/models"

	"github.com/google/uuid"
)

// PostgresMagicLinkRepository implementação PostgreSQL do MagicLinkRepository
type PostgresMagicLinkRepository struct {
	db *sql.DB
}

func NewPostgresMagicLinkRepository(db *sql.DB) *PostgresMagicLinkRepository {
	return &PostgresMagicLinkRepository{db: db}
}

func (r *PostgresMagicLinkRepository) Create(ctx context.Context, link *models.MagicLink) error {
	query := `
		INSERT INTO magic_links (id, email, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		link.ID, link.Email, link.Token, link.ExpiresAt, link.Used, link.CreatedAt,
	)
	return err
}

func (r *PostgresMagicLinkRepository) GetByToken(ctx context.Context, token string) (*models.MagicLink, error) {
	query := `
		SELECT id, email, token, expires_at, used, created_at
		FROM magic_links WHERE token = $1`

	link := &models.MagicLink{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&link.ID, &link.Email, &link.Token, &link.ExpiresAt, &link.Used, &link.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("magic link not found")
		}
		return nil, fmt.Errorf("failed to get magic link: %w", err)
	}

	return link, nil
}

func (r *PostgresMagicLinkRepository) MarkAsUsed(ctx context.Context, token string) error {
	query := `UPDATE magic_links SET used = true WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *PostgresMagicLinkRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM magic_links WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresMagicLinkRepository) DeleteByEmail(ctx context.Context, email string) error {
	query := `DELETE FROM magic_links WHERE email = $1`
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}

// PostgresRefreshTokenRepository implementação PostgreSQL do RefreshTokenRepository
type PostgresRefreshTokenRepository struct {
	db *sql.DB
}

func NewPostgresRefreshTokenRepository(db *sql.DB) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

func (r *PostgresRefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.Token, token.ExpiresAt, token.Used, token.CreatedAt,
	)
	return err
}

func (r *PostgresRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM refresh_tokens WHERE token = $1`

	refreshToken := &models.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token,
		&refreshToken.ExpiresAt, &refreshToken.Used, &refreshToken.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return refreshToken, nil
}

func (r *PostgresRefreshTokenRepository) MarkAsUsed(ctx context.Context, token string) error {
	query := `UPDATE refresh_tokens SET used = true WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *PostgresRefreshTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *PostgresRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

// PostgresAuthProviderRepository implementação PostgreSQL do AuthProviderRepository
type PostgresAuthProviderRepository struct {
	db *sql.DB
}

func NewPostgresAuthProviderRepository(db *sql.DB) *PostgresAuthProviderRepository {
	return &PostgresAuthProviderRepository{db: db}
}

func (r *PostgresAuthProviderRepository) Create(ctx context.Context, provider *models.UserAuthProvider) error {
	query := `
		INSERT INTO user_auth_providers (id, user_id, provider, provider_user_id, provider_data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		provider.ID, provider.UserID, provider.Provider,
		provider.ProviderUserID, provider.ProviderData, provider.CreatedAt,
	)
	return err
}

func (r *PostgresAuthProviderRepository) GetByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider models.AuthProvider) (*models.UserAuthProvider, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_data, created_at
		FROM user_auth_providers WHERE user_id = $1 AND provider = $2`

	authProvider := &models.UserAuthProvider{}
	err := r.db.QueryRowContext(ctx, query, userID, provider).Scan(
		&authProvider.ID, &authProvider.UserID, &authProvider.Provider,
		&authProvider.ProviderUserID, &authProvider.ProviderData, &authProvider.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("auth provider not found")
		}
		return nil, fmt.Errorf("failed to get auth provider: %w", err)
	}

	return authProvider, nil
}

func (r *PostgresAuthProviderRepository) GetByProviderAndUserID(ctx context.Context, provider models.AuthProvider, providerUserID string) (*models.UserAuthProvider, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_data, created_at
		FROM user_auth_providers WHERE provider = $1 AND provider_user_id = $2`

	authProvider := &models.UserAuthProvider{}
	err := r.db.QueryRowContext(ctx, query, provider, providerUserID).Scan(
		&authProvider.ID, &authProvider.UserID, &authProvider.Provider,
		&authProvider.ProviderUserID, &authProvider.ProviderData, &authProvider.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("auth provider not found")
		}
		return nil, fmt.Errorf("failed to get auth provider: %w", err)
	}

	return authProvider, nil
}

func (r *PostgresAuthProviderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_auth_providers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresAuthProviderRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserAuthProvider, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_data, created_at
		FROM user_auth_providers WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list auth providers: %w", err)
	}
	defer rows.Close()

	var providers []*models.UserAuthProvider
	for rows.Next() {
		provider := &models.UserAuthProvider{}
		err := rows.Scan(
			&provider.ID, &provider.UserID, &provider.Provider,
			&provider.ProviderUserID, &provider.ProviderData, &provider.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan auth provider: %w", err)
		}
		providers = append(providers, provider)
	}

	return providers, nil
}
