package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"pagemagic/build-svc/internal/config"
	"pagemagic/build-svc/internal/models"
)

type TemplateService struct {
	db     *sql.DB
	config *config.Config
}

func NewTemplateService(db *sql.DB, cfg *config.Config) *TemplateService {
	return &TemplateService{
		db:     db,
		config: cfg,
	}
}

func (s *TemplateService) ListTemplates(ctx context.Context, category string, limit, offset int) ([]*models.Template, error) {
	var query string
	var args []interface{}

	if category != "" {
		query = `
			SELECT id, name, description, category, framework, preview, thumbnail, 
				   tags, version, author, license, featured, downloads, rating, created_at, updated_at
			FROM templates 
			WHERE category = $1
			ORDER BY featured DESC, downloads DESC, rating DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{category, limit, offset}
	} else {
		query = `
			SELECT id, name, description, category, framework, preview, thumbnail,
				   tags, version, author, license, featured, downloads, rating, created_at, updated_at
			FROM templates 
			ORDER BY featured DESC, downloads DESC, rating DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	defer rows.Close()

	var templates []*models.Template
	for rows.Next() {
		template := &models.Template{}
		var tagsJSON []byte

		err := rows.Scan(
			&template.ID, &template.Name, &template.Description, &template.Category,
			&template.Framework, &template.Preview, &template.Thumbnail,
			&tagsJSON, &template.Version, &template.Author, &template.License,
			&template.Featured, &template.Downloads, &template.Rating,
			&template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &template.Tags)
		}

		templates = append(templates, template)
	}

	return templates, nil
}

func (s *TemplateService) GetTemplate(ctx context.Context, templateID string) (*models.Template, error) {
	query := `
		SELECT id, name, description, category, framework, preview, thumbnail,
			   config, files, variables, tags, version, author, license, featured, downloads, rating,
			   created_at, updated_at
		FROM templates 
		WHERE id = $1
	`

	template := &models.Template{}
	var configJSON, filesJSON, variablesJSON, tagsJSON []byte

	err := s.db.QueryRowContext(ctx, query, templateID).Scan(
		&template.ID, &template.Name, &template.Description, &template.Category,
		&template.Framework, &template.Preview, &template.Thumbnail,
		&configJSON, &filesJSON, &variablesJSON, &tagsJSON,
		&template.Version, &template.Author, &template.License, &template.Featured,
		&template.Downloads, &template.Rating, &template.CreatedAt, &template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Unmarshal JSON fields
	if len(configJSON) > 0 {
		json.Unmarshal(configJSON, &template.Config)
	}
	if len(filesJSON) > 0 {
		json.Unmarshal(filesJSON, &template.Files)
	}
	if len(variablesJSON) > 0 {
		json.Unmarshal(variablesJSON, &template.Variables)
	}
	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &template.Tags)
	}

	return template, nil
}

func (s *TemplateService) GetFeaturedTemplates(ctx context.Context, limit int) ([]*models.Template, error) {
	query := `
		SELECT id, name, description, category, framework, preview, thumbnail,
			   tags, version, author, featured, downloads, rating, created_at, updated_at
		FROM templates 
		WHERE featured = true
		ORDER BY downloads DESC, rating DESC
		LIMIT $1
	`

	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get featured templates: %w", err)
	}
	defer rows.Close()

	var templates []*models.Template
	for rows.Next() {
		template := &models.Template{}
		var tagsJSON []byte

		err := rows.Scan(
			&template.ID, &template.Name, &template.Description, &template.Category,
			&template.Framework, &template.Preview, &template.Thumbnail,
			&tagsJSON, &template.Version, &template.Author, &template.Featured,
			&template.Downloads, &template.Rating, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &template.Tags)
		}

		templates = append(templates, template)
	}

	return templates, nil
}

func (s *TemplateService) SearchTemplates(ctx context.Context, query string, category, framework string, limit, offset int) ([]*models.Template, error) {
	sqlQuery := `
		SELECT id, name, description, category, framework, preview, thumbnail,
			   tags, version, author, featured, downloads, rating, created_at, updated_at
		FROM templates 
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if query != "" {
		sqlQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d OR tags::text ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+query+"%")
		argIndex++
	}

	if category != "" {
		sqlQuery += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if framework != "" {
		sqlQuery += fmt.Sprintf(" AND framework = $%d", argIndex)
		args = append(args, framework)
		argIndex++
	}

	sqlQuery += " ORDER BY featured DESC, downloads DESC, rating DESC"
	sqlQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search templates: %w", err)
	}
	defer rows.Close()

	var templates []*models.Template
	for rows.Next() {
		template := &models.Template{}
		var tagsJSON []byte

		err := rows.Scan(
			&template.ID, &template.Name, &template.Description, &template.Category,
			&template.Framework, &template.Preview, &template.Thumbnail,
			&tagsJSON, &template.Version, &template.Author, &template.Featured,
			&template.Downloads, &template.Rating, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &template.Tags)
		}

		templates = append(templates, template)
	}

	return templates, nil
}

func (s *TemplateService) IncrementDownloads(ctx context.Context, templateID string) error {
	query := `UPDATE templates SET downloads = downloads + 1 WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, templateID)
	if err != nil {
		return fmt.Errorf("failed to increment downloads: %w", err)
	}

	return nil
}

func (s *TemplateService) GetCategories(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT category FROM templates ORDER BY category`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *TemplateService) GetFrameworks(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT framework FROM templates ORDER BY framework`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get frameworks: %w", err)
	}
	defer rows.Close()

	var frameworks []string
	for rows.Next() {
		var framework string
		if err := rows.Scan(&framework); err != nil {
			return nil, fmt.Errorf("failed to scan framework: %w", err)
		}
		frameworks = append(frameworks, framework)
	}

	return frameworks, nil
}
