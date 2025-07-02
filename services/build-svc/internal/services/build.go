package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"pagemagic/build-svc/internal/config"
	"pagemagic/build-svc/internal/models"
)

type BuildService struct {
	db     *sql.DB
	config *config.Config
}

func NewBuildService(db *sql.DB, cfg *config.Config) *BuildService {
	return &BuildService{
		db:     db,
		config: cfg,
	}
}

func (s *BuildService) CreateBuildJob(ctx context.Context, userID string, req *models.CreateBuildJobRequest) (*models.BuildJob, error) {
	now := time.Now()
	jobID := fmt.Sprintf("build_%d", now.UnixNano())

	job := &models.BuildJob{
		ID:          jobID,
		SiteID:      req.SiteID,
		UserID:      userID,
		Type:        req.Type,
		Status:      "pending",
		SourceType:  req.SourceType,
		SourceData:  req.SourceData,
		BuildConfig: req.BuildConfig,
		Logs:        []models.BuildLog{},
		StartedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `
		INSERT INTO build_jobs (id, site_id, user_id, type, status, source_type, source_data, build_config, logs, started_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, site_id, user_id, type, status, source_type, started_at, created_at, updated_at
	`

	sourceDataJSON, _ := json.Marshal(job.SourceData)
	buildConfigJSON, _ := json.Marshal(job.BuildConfig)
	logsJSON, _ := json.Marshal(job.Logs)

	err := s.db.QueryRowContext(ctx, query,
		job.ID, job.SiteID, job.UserID, job.Type, job.Status, job.SourceType,
		sourceDataJSON, buildConfigJSON, logsJSON, job.StartedAt, job.CreatedAt, job.UpdatedAt,
	).Scan(
		&job.ID, &job.SiteID, &job.UserID, &job.Type, &job.Status, &job.SourceType,
		&job.StartedAt, &job.CreatedAt, &job.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create build job: %w", err)
	}

	// Iniciar processo de build em background
	go s.processBuildJob(context.Background(), job)

	return job, nil
}

func (s *BuildService) GetBuildJob(ctx context.Context, jobID, userID string) (*models.BuildJob, error) {
	query := `
		SELECT id, site_id, user_id, type, status, source_type, source_data, build_config, 
			   output, logs, error, started_at, ended_at, created_at, updated_at
		FROM build_jobs 
		WHERE id = $1 AND user_id = $2
	`

	job := &models.BuildJob{}
	var sourceDataJSON, buildConfigJSON, outputJSON, logsJSON []byte
	var endedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, query, jobID, userID).Scan(
		&job.ID, &job.SiteID, &job.UserID, &job.Type, &job.Status, &job.SourceType,
		&sourceDataJSON, &buildConfigJSON, &outputJSON, &logsJSON, &job.Error,
		&job.StartedAt, &endedAt, &job.CreatedAt, &job.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("build job not found")
		}
		return nil, fmt.Errorf("failed to get build job: %w", err)
	}

	if endedAt.Valid {
		job.EndedAt = &endedAt.Time
	}

	// Unmarshal JSON fields
	json.Unmarshal(sourceDataJSON, &job.SourceData)
	json.Unmarshal(buildConfigJSON, &job.BuildConfig)
	json.Unmarshal(outputJSON, &job.Output)
	json.Unmarshal(logsJSON, &job.Logs)

	return job, nil
}

func (s *BuildService) ListBuildJobs(ctx context.Context, siteID, userID string, limit, offset int) ([]*models.BuildJob, error) {
	query := `
		SELECT id, site_id, user_id, type, status, source_type, started_at, ended_at, created_at, updated_at
		FROM build_jobs 
		WHERE site_id = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := s.db.QueryContext(ctx, query, siteID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list build jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*models.BuildJob
	for rows.Next() {
		job := &models.BuildJob{}
		var endedAt sql.NullTime

		err := rows.Scan(
			&job.ID, &job.SiteID, &job.UserID, &job.Type, &job.Status, &job.SourceType,
			&job.StartedAt, &endedAt, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan build job: %w", err)
		}

		if endedAt.Valid {
			job.EndedAt = &endedAt.Time
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (s *BuildService) UpdateBuildJobStatus(ctx context.Context, jobID, status, errorMsg string) error {
	now := time.Now()
	var endedAt *time.Time

	if status == "completed" || status == "failed" {
		endedAt = &now
	}

	query := `
		UPDATE build_jobs 
		SET status = $1, error = $2, ended_at = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := s.db.ExecContext(ctx, query, status, errorMsg, endedAt, now, jobID)
	if err != nil {
		return fmt.Errorf("failed to update build job status: %w", err)
	}

	return nil
}

func (s *BuildService) AddBuildLog(ctx context.Context, jobID, level, message string, data interface{}) error {
	log := models.BuildLog{
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Data:      data,
	}

	// Obter logs atuais
	var logsJSON []byte
	query := `SELECT logs FROM build_jobs WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, jobID).Scan(&logsJSON)
	if err != nil {
		return fmt.Errorf("failed to get current logs: %w", err)
	}

	var logs []models.BuildLog
	if len(logsJSON) > 0 {
		json.Unmarshal(logsJSON, &logs)
	}

	logs = append(logs, log)
	newLogsJSON, _ := json.Marshal(logs)

	// Atualizar logs
	updateQuery := `UPDATE build_jobs SET logs = $1, updated_at = $2 WHERE id = $3`
	_, err = s.db.ExecContext(ctx, updateQuery, newLogsJSON, time.Now(), jobID)
	if err != nil {
		return fmt.Errorf("failed to update logs: %w", err)
	}

	return nil
}

func (s *BuildService) processBuildJob(ctx context.Context, job *models.BuildJob) {
	// Log início do processo
	s.AddBuildLog(ctx, job.ID, "info", "Starting build process", nil)

	// Atualizar status para "building"
	s.UpdateBuildJobStatus(ctx, job.ID, "building", "")

	// Simular processo de build baseado no tipo de fonte
	var output *models.BuildOutput
	var err error

	switch job.SourceType {
	case "visual_editor":
		output, err = s.buildFromVisualEditor(ctx, job)
	case "code":
		output, err = s.buildFromCode(ctx, job)
	case "template":
		output, err = s.buildFromTemplate(ctx, job)
	default:
		err = fmt.Errorf("unsupported source type: %s", job.SourceType)
	}

	if err != nil {
		s.AddBuildLog(ctx, job.ID, "error", "Build failed", map[string]string{"error": err.Error()})
		s.UpdateBuildJobStatus(ctx, job.ID, "failed", err.Error())
		return
	}

	// Salvar output
	s.saveBuildOutput(ctx, job.ID, output)

	// Log conclusão
	s.AddBuildLog(ctx, job.ID, "info", "Build completed successfully", map[string]interface{}{
		"files_count": len(output.Files),
		"total_size":  output.Stats.TotalSize,
		"duration":    output.Stats.Duration,
	})

	s.UpdateBuildJobStatus(ctx, job.ID, "completed", "")
}

func (s *BuildService) buildFromVisualEditor(ctx context.Context, job *models.BuildJob) (*models.BuildOutput, error) {
	s.AddBuildLog(ctx, job.ID, "info", "Building from visual editor data", nil)
	
	// Simular processo de build
	time.Sleep(2 * time.Second)

	// Gerar HTML a partir dos dados do editor visual
	htmlContent := s.generateHTMLFromVisualData(job.SourceData)
	cssContent := s.generateCSSFromVisualData(job.SourceData)

	files := []models.BuildFile{
		{
			Path:        "/index.html",
			Content:     htmlContent,
			ContentType: "text/html",
			Size:        int64(len(htmlContent)),
			Hash:        fmt.Sprintf("%x", time.Now().UnixNano()),
		},
		{
			Path:        "/style.css",
			Content:     cssContent,
			ContentType: "text/css",
			Size:        int64(len(cssContent)),
			Hash:        fmt.Sprintf("%x", time.Now().UnixNano()),
		},
	}

	stats := models.BuildStats{
		TotalFiles: len(files),
		TotalSize:  int64(len(htmlContent) + len(cssContent)),
		Duration:   2 * time.Second,
	}

	return &models.BuildOutput{
		Files: files,
		Stats: stats,
	}, nil
}

func (s *BuildService) buildFromCode(ctx context.Context, job *models.BuildJob) (*models.BuildOutput, error) {
	s.AddBuildLog(ctx, job.ID, "info", "Building from code", nil)
	
	// Simular processo de build mais complexo
	time.Sleep(5 * time.Second)

	// Em implementação real, aqui seria:
	// 1. Fazer download do código fonte
	// 2. Instalar dependências (npm install, etc.)
	// 3. Executar build (webpack, vite, etc.)
	// 4. Otimizar assets
	// 5. Gerar manifests

	return &models.BuildOutput{
		Files: []models.BuildFile{},
		Stats: models.BuildStats{
			Duration: 5 * time.Second,
		},
	}, nil
}

func (s *BuildService) buildFromTemplate(ctx context.Context, job *models.BuildJob) (*models.BuildOutput, error) {
	s.AddBuildLog(ctx, job.ID, "info", "Building from template", nil)
	
	// Simular processo de build de template
	time.Sleep(3 * time.Second)

	return &models.BuildOutput{
		Files: []models.BuildFile{},
		Stats: models.BuildStats{
			Duration: 3 * time.Second,
		},
	}, nil
}

func (s *BuildService) generateHTMLFromVisualData(data map[string]interface{}) string {
	// Simular geração de HTML a partir dos dados do editor visual
	return `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Site PageMagic</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <header>
        <h1>Bem-vindo ao PageMagic</h1>
    </header>
    <main>
        <section>
            <h2>Criado com PageMagic</h2>
            <p>Este site foi criado usando o editor visual do PageMagic.</p>
        </section>
    </main>
    <footer>
        <p>&copy; 2024 PageMagic</p>
    </footer>
</body>
</html>`
}

func (s *BuildService) generateCSSFromVisualData(data map[string]interface{}) string {
	// Simular geração de CSS a partir dos dados do editor visual
	return `body {
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 0;
    line-height: 1.6;
}

header {
    background: #333;
    color: white;
    padding: 1rem;
    text-align: center;
}

main {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
}

footer {
    background: #f4f4f4;
    padding: 1rem;
    text-align: center;
    margin-top: 2rem;
}`
}

func (s *BuildService) saveBuildOutput(ctx context.Context, jobID string, output *models.BuildOutput) error {
	outputJSON, _ := json.Marshal(output)

	query := `UPDATE build_jobs SET output = $1, updated_at = $2 WHERE id = $3`
	_, err := s.db.ExecContext(ctx, query, outputJSON, time.Now(), jobID)
	if err != nil {
		return fmt.Errorf("failed to save build output: %w", err)
	}

	return nil
}
