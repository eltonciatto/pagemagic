package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// New cria uma nova instância do repositório
func New(databaseURL string) (*Repository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Inicializar repositórios específicos
	userRepo := NewPostgresUserRepository(db)
	magicLinkRepo := NewPostgresMagicLinkRepository(db)
	refreshTokenRepo := NewPostgresRefreshTokenRepository(db)
	authProviderRepo := NewPostgresAuthProviderRepository(db)

	return &Repository{
		User:         userRepo,
		MagicLink:    magicLinkRepo,
		RefreshToken: refreshTokenRepo,
		AuthProvider: authProviderRepo,
	}, nil
}

// Close fecha a conexão com o banco de dados
func (r *Repository) Close() error {
	// Implementar close para cada repositório
	return nil
}
