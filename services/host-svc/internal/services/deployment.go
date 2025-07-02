package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"pagemagic/host-svc/internal/config"
	"pagemagic/host-svc/internal/models"
)

type DeploymentService struct {
	db     *sql.DB
	config *config.Config
}

func NewDeploymentService(db *sql.DB, cfg *config.Config) *DeploymentService {
	return &DeploymentService{
		db:     db,
		config: cfg,
	}
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, siteID, buildID string, files []models.DeploymentFile, siteConfig models.SiteConfig) (*models.Deployment, error) {
	now := time.Now()
	deploymentID := fmt.Sprintf("deploy_%d", now.UnixNano())

	// Obter próxima versão
	version, err := s.getNextVersion(ctx, siteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get next version: %w", err)
	}

	stats := models.DeploymentStats{
		TotalFiles: len(files),
		TotalSize:  s.calculateTotalSize(files),
	}

	query := `
		INSERT INTO deployments (id, site_id, build_id, version, status, files, config, stats, started_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, site_id, build_id, version, status, started_at, created_at
	`

	deployment := &models.Deployment{}
	err = s.db.QueryRowContext(ctx, query,
		deploymentID, siteID, buildID, version, "pending", files, siteConfig, stats, now, now,
	).Scan(
		&deployment.ID, &deployment.SiteID, &deployment.BuildID, &deployment.Version,
		&deployment.Status, &deployment.StartedAt, &deployment.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	deployment.Files = files
	deployment.Config = siteConfig
	deployment.Stats = stats

	// Processar deployment em background
	go s.processDeployment(context.Background(), deployment)

	return deployment, nil
}

func (s *DeploymentService) GetDeployment(ctx context.Context, deploymentID, siteID string) (*models.Deployment, error) {
	query := `
		SELECT id, site_id, build_id, version, status, files, config, stats, error, started_at, ended_at, created_at
		FROM deployments 
		WHERE id = $1 AND site_id = $2
	`

	deployment := &models.Deployment{}
	err := s.db.QueryRowContext(ctx, query, deploymentID, siteID).Scan(
		&deployment.ID, &deployment.SiteID, &deployment.BuildID, &deployment.Version,
		&deployment.Status, &deployment.Files, &deployment.Config, &deployment.Stats,
		&deployment.Error, &deployment.StartedAt, &deployment.EndedAt, &deployment.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deployment not found")
		}
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return deployment, nil
}

func (s *DeploymentService) ListDeployments(ctx context.Context, siteID string, limit, offset int) ([]*models.Deployment, error) {
	query := `
		SELECT id, site_id, build_id, version, status, files, config, stats, error, started_at, ended_at, created_at
		FROM deployments 
		WHERE site_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, query, siteID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	defer rows.Close()

	var deployments []*models.Deployment
	for rows.Next() {
		deployment := &models.Deployment{}
		err := rows.Scan(
			&deployment.ID, &deployment.SiteID, &deployment.BuildID, &deployment.Version,
			&deployment.Status, &deployment.Files, &deployment.Config, &deployment.Stats,
			&deployment.Error, &deployment.StartedAt, &deployment.EndedAt, &deployment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment: %w", err)
		}
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

func (s *DeploymentService) UpdateDeploymentStatus(ctx context.Context, deploymentID, status, errorMsg string) error {
	now := time.Now()
	var endedAt *time.Time

	if status == "deployed" || status == "failed" {
		endedAt = &now
	}

	query := `
		UPDATE deployments 
		SET status = $1, error = $2, ended_at = $3
		WHERE id = $4
	`

	_, err := s.db.ExecContext(ctx, query, status, errorMsg, endedAt, deploymentID)
	if err != nil {
		return fmt.Errorf("failed to update deployment status: %w", err)
	}

	return nil
}

func (s *DeploymentService) RollbackToDeployment(ctx context.Context, siteID, deploymentID string) error {
	// Obter deployment alvo
	deployment, err := s.GetDeployment(ctx, deploymentID, siteID)
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	if deployment.Status != "deployed" {
		return fmt.Errorf("cannot rollback to non-deployed version")
	}

	// Atualizar site para usar esta versão
	query := `
		UPDATE sites 
		SET build_id = $1, version = $2, config = $3, updated_at = $4
		WHERE id = $5
	`

	_, err = s.db.ExecContext(ctx, query, deployment.BuildID, deployment.Version, deployment.Config, time.Now(), siteID)
	if err != nil {
		return fmt.Errorf("failed to rollback site: %w", err)
	}

	return nil
}

func (s *DeploymentService) getNextVersion(ctx context.Context, siteID string) (int, error) {
	query := `SELECT COALESCE(MAX(version), 0) + 1 FROM deployments WHERE site_id = $1`

	var version int
	err := s.db.QueryRowContext(ctx, query, siteID).Scan(&version)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func (s *DeploymentService) calculateTotalSize(files []models.DeploymentFile) int64 {
	var total int64
	for _, file := range files {
		total += file.Size
	}
	return total
}

func (s *DeploymentService) processDeployment(ctx context.Context, deployment *models.Deployment) {
	// Atualizar status para "deploying"
	s.UpdateDeploymentStatus(ctx, deployment.ID, "deploying", "")

	// Simular processo de deployment
	// Em uma implementação real, aqui seria:
	// 1. Upload dos arquivos para storage (S3/GCS)
	// 2. Configuração do CDN
	// 3. Configuração de SSL se necessário
	// 4. Atualização de DNS
	// 5. Warm up do cache

	time.Sleep(5 * time.Second) // Simular processo

	// Por agora, marcar como deployado
	s.UpdateDeploymentStatus(ctx, deployment.ID, "deployed", "")

	// Atualizar site
	s.updateSiteAfterDeployment(ctx, deployment)
}

func (s *DeploymentService) updateSiteAfterDeployment(ctx context.Context, deployment *models.Deployment) {
	now := time.Now()
	query := `
		UPDATE sites 
		SET build_id = $1, version = $2, config = $3, status = 'active', last_deploy = $4, updated_at = $5
		WHERE id = $6
	`

	s.db.ExecContext(ctx, query, deployment.BuildID, deployment.Version, deployment.Config, now, now, deployment.SiteID)
}
