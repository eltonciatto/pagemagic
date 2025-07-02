package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"pagemagic/host-svc/internal/config"
	"pagemagic/host-svc/internal/models"
)

type HostingService struct {
	db     *sql.DB
	config *config.Config
}

func NewHostingService(db *sql.DB, cfg *config.Config) *HostingService {
	return &HostingService{
		db:     db,
		config: cfg,
	}
}

func (s *HostingService) CreateSite(ctx context.Context, userID string, req *models.CreateSiteRequest) (*models.Site, error) {
	now := time.Now()
	siteID := fmt.Sprintf("site_%d", now.UnixNano())

	query := `
		INSERT INTO sites (id, user_id, domain, subdomain, custom_domain, status, config, ssl_enabled, cdn_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, user_id, domain, subdomain, custom_domain, status, ssl_enabled, cdn_enabled, created_at, updated_at
	`

	site := &models.Site{}
	err := s.db.QueryRowContext(ctx, query,
		siteID, userID, req.Domain, req.Domain, req.CustomDomain,
		"pending", req.Config, req.SSLEnabled, req.CDNEnabled, now, now,
	).Scan(
		&site.ID, &site.UserID, &site.Domain, &site.Subdomain, &site.CustomDomain,
		&site.Status, &site.SSLEnabled, &site.CDNEnabled, &site.CreatedAt, &site.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}

	return site, nil
}

func (s *HostingService) GetSite(ctx context.Context, siteID, userID string) (*models.Site, error) {
	query := `
		SELECT id, user_id, domain, subdomain, custom_domain, status, build_id, version, 
			   config, ssl_enabled, ssl_cert_path, cdn_enabled, cdn_url, analytics, 
			   last_deploy, created_at, updated_at
		FROM sites 
		WHERE id = $1 AND user_id = $2
	`

	site := &models.Site{}
	err := s.db.QueryRowContext(ctx, query, siteID, userID).Scan(
		&site.ID, &site.UserID, &site.Domain, &site.Subdomain, &site.CustomDomain,
		&site.Status, &site.BuildID, &site.Version, &site.Config, &site.SSLEnabled,
		&site.SSLCertPath, &site.CDNEnabled, &site.CDNUrl, &site.Analytics,
		&site.LastDeploy, &site.CreatedAt, &site.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("site not found")
		}
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	return site, nil
}

func (s *HostingService) ListSites(ctx context.Context, userID string) ([]*models.Site, error) {
	query := `
		SELECT id, user_id, domain, subdomain, custom_domain, status, build_id, version,
			   config, ssl_enabled, ssl_cert_path, cdn_enabled, cdn_url, analytics,
			   last_deploy, created_at, updated_at
		FROM sites 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sites: %w", err)
	}
	defer rows.Close()

	var sites []*models.Site
	for rows.Next() {
		site := &models.Site{}
		err := rows.Scan(
			&site.ID, &site.UserID, &site.Domain, &site.Subdomain, &site.CustomDomain,
			&site.Status, &site.BuildID, &site.Version, &site.Config, &site.SSLEnabled,
			&site.SSLCertPath, &site.CDNEnabled, &site.CDNUrl, &site.Analytics,
			&site.LastDeploy, &site.CreatedAt, &site.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan site: %w", err)
		}
		sites = append(sites, site)
	}

	return sites, nil
}

func (s *HostingService) UpdateSite(ctx context.Context, siteID, userID string, req *models.UpdateSiteRequest) (*models.Site, error) {
	// Verificar se o site existe
	site, err := s.GetSite(ctx, siteID, userID)
	if err != nil {
		return nil, err
	}

	// Atualizar campos se fornecidos
	if req.Config != nil {
		site.Config = *req.Config
	}
	if req.CustomDomain != nil {
		site.CustomDomain = *req.CustomDomain
	}
	if req.SSLEnabled != nil {
		site.SSLEnabled = *req.SSLEnabled
	}
	if req.CDNEnabled != nil {
		site.CDNEnabled = *req.CDNEnabled
	}

	site.UpdatedAt = time.Now()

	query := `
		UPDATE sites 
		SET config = $1, custom_domain = $2, ssl_enabled = $3, cdn_enabled = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
	`

	_, err = s.db.ExecContext(ctx, query,
		site.Config, site.CustomDomain, site.SSLEnabled, site.CDNEnabled,
		site.UpdatedAt, siteID, userID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update site: %w", err)
	}

	return site, nil
}

func (s *HostingService) DeleteSite(ctx context.Context, siteID, userID string) error {
	query := `DELETE FROM sites WHERE id = $1 AND user_id = $2`

	result, err := s.db.ExecContext(ctx, query, siteID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("site not found")
	}

	return nil
}

func (s *HostingService) UpdateSiteStatus(ctx context.Context, siteID, status string) error {
	query := `UPDATE sites SET status = $1, updated_at = $2 WHERE id = $3`

	_, err := s.db.ExecContext(ctx, query, status, time.Now(), siteID)
	if err != nil {
		return fmt.Errorf("failed to update site status: %w", err)
	}

	return nil
}

func (s *HostingService) GetSiteByDomain(ctx context.Context, domain string) (*models.Site, error) {
	query := `
		SELECT id, user_id, domain, subdomain, custom_domain, status, build_id, version,
			   config, ssl_enabled, ssl_cert_path, cdn_enabled, cdn_url, analytics,
			   last_deploy, created_at, updated_at
		FROM sites 
		WHERE domain = $1 OR custom_domain = $1 OR subdomain = $1
		AND status = 'active'
	`

	site := &models.Site{}
	err := s.db.QueryRowContext(ctx, query, domain).Scan(
		&site.ID, &site.UserID, &site.Domain, &site.Subdomain, &site.CustomDomain,
		&site.Status, &site.BuildID, &site.Version, &site.Config, &site.SSLEnabled,
		&site.SSLCertPath, &site.CDNEnabled, &site.CDNUrl, &site.Analytics,
		&site.LastDeploy, &site.CreatedAt, &site.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("site not found for domain: %s", domain)
		}
		return nil, fmt.Errorf("failed to get site by domain: %w", err)
	}

	return site, nil
}
