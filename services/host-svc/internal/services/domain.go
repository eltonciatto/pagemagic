package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"pagemagic/host-svc/internal/config"
	"pagemagic/host-svc/internal/models"
)

type DomainService struct {
	db     *sql.DB
	config *config.Config
}

func NewDomainService(db *sql.DB, cfg *config.Config) *DomainService {
	return &DomainService{
		db:     db,
		config: cfg,
	}
}

func (s *DomainService) AddDomain(ctx context.Context, siteID, domain, domainType string) (*models.DomainConfig, error) {
	now := time.Now()
	domainID := fmt.Sprintf("domain_%d", now.UnixNano())

	// Gerar registros DNS necessários
	dnsRecords := s.generateDNSRecords(domain, domainType)

	query := `
		INSERT INTO domain_configs (id, site_id, domain, type, status, dns_records, ssl_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, site_id, domain, type, status, dns_records, ssl_status, created_at, updated_at
	`

	domainConfig := &models.DomainConfig{}
	err := s.db.QueryRowContext(ctx, query,
		domainID, siteID, domain, domainType, "pending", dnsRecords, "pending", now, now,
	).Scan(
		&domainConfig.ID, &domainConfig.SiteID, &domainConfig.Domain, &domainConfig.Type,
		&domainConfig.Status, &domainConfig.DNSRecords, &domainConfig.SSLStatus,
		&domainConfig.CreatedAt, &domainConfig.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to add domain: %w", err)
	}

	// Iniciar processo de verificação em background
	go s.processDomainVerification(context.Background(), domainConfig)

	return domainConfig, nil
}

func (s *DomainService) ListDomains(ctx context.Context, siteID string) ([]*models.DomainConfig, error) {
	query := `
		SELECT id, site_id, domain, type, status, dns_records, ssl_status, ssl_provider, verified_at, created_at, updated_at
		FROM domain_configs 
		WHERE site_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, siteID)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}
	defer rows.Close()

	var domains []*models.DomainConfig
	for rows.Next() {
		domain := &models.DomainConfig{}
		err := rows.Scan(
			&domain.ID, &domain.SiteID, &domain.Domain, &domain.Type, &domain.Status,
			&domain.DNSRecords, &domain.SSLStatus, &domain.SSLProvider, &domain.VerifiedAt,
			&domain.CreatedAt, &domain.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan domain: %w", err)
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

func (s *DomainService) UpdateDomain(ctx context.Context, domainID, siteID string, updates map[string]interface{}) (*models.DomainConfig, error) {
	// Implementar lógica de atualização de domínio
	query := `UPDATE domain_configs SET updated_at = $1 WHERE id = $2 AND site_id = $3`

	_, err := s.db.ExecContext(ctx, query, time.Now(), domainID, siteID)
	if err != nil {
		return nil, fmt.Errorf("failed to update domain: %w", err)
	}

	return s.GetDomain(ctx, domainID, siteID)
}

func (s *DomainService) DeleteDomain(ctx context.Context, domainID, siteID string) error {
	query := `DELETE FROM domain_configs WHERE id = $1 AND site_id = $2`

	result, err := s.db.ExecContext(ctx, query, domainID, siteID)
	if err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("domain not found")
	}

	return nil
}

func (s *DomainService) GetDomain(ctx context.Context, domainID, siteID string) (*models.DomainConfig, error) {
	query := `
		SELECT id, site_id, domain, type, status, dns_records, ssl_status, ssl_provider, verified_at, created_at, updated_at
		FROM domain_configs 
		WHERE id = $1 AND site_id = $2
	`

	domain := &models.DomainConfig{}
	err := s.db.QueryRowContext(ctx, query, domainID, siteID).Scan(
		&domain.ID, &domain.SiteID, &domain.Domain, &domain.Type, &domain.Status,
		&domain.DNSRecords, &domain.SSLStatus, &domain.SSLProvider, &domain.VerifiedAt,
		&domain.CreatedAt, &domain.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("domain not found")
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return domain, nil
}

func (s *DomainService) VerifyDomain(ctx context.Context, domainID, siteID string) error {
	domain, err := s.GetDomain(ctx, domainID, siteID)
	if err != nil {
		return err
	}

	// Implementar verificação real de DNS
	verified := s.verifyDNSRecords(domain.Domain, domain.DNSRecords)

	status := "failed"
	var verifiedAt *time.Time
	if verified {
		status = "active"
		now := time.Now()
		verifiedAt = &now
	}

	query := `
		UPDATE domain_configs 
		SET status = $1, verified_at = $2, updated_at = $3
		WHERE id = $4 AND site_id = $5
	`

	_, err = s.db.ExecContext(ctx, query, status, verifiedAt, time.Now(), domainID, siteID)
	if err != nil {
		return fmt.Errorf("failed to update domain verification: %w", err)
	}

	return nil
}

func (s *DomainService) generateDNSRecords(domain, domainType string) []models.DNSRecord {
	records := []models.DNSRecord{}

	if domainType == "custom" {
		// A record apontando para o IP do servidor
		records = append(records, models.DNSRecord{
			Type:  "A",
			Name:  "@",
			Value: s.config.Server.Host,
			TTL:   300,
		})

		// CNAME para www
		records = append(records, models.DNSRecord{
			Type:  "CNAME",
			Name:  "www",
			Value: domain,
			TTL:   300,
		})
	} else {
		// Para subdomínio, apenas CNAME
		records = append(records, models.DNSRecord{
			Type:  "CNAME",
			Name:  domain,
			Value: "pagemagic.app",
			TTL:   300,
		})
	}

	return records
}

func (s *DomainService) verifyDNSRecords(domain string, records []models.DNSRecord) bool {
	// Implementar verificação real de DNS usando bibliotecas como "net" ou "miekg/dns"
	// Por enquanto, simular verificação
	return true
}

func (s *DomainService) processDomainVerification(ctx context.Context, domain *models.DomainConfig) {
	// Aguardar um tempo para propagação DNS
	time.Sleep(30 * time.Second)

	// Verificar registros DNS
	verified := s.verifyDNSRecords(domain.Domain, domain.DNSRecords)

	status := "failed"
	var verifiedAt *time.Time
	if verified {
		status = "active"
		now := time.Now()
		verifiedAt = &now

		// Se verificado, iniciar processo de SSL
		go s.processSSLSetup(ctx, domain)
	}

	query := `
		UPDATE domain_configs 
		SET status = $1, verified_at = $2, updated_at = $3
		WHERE id = $4
	`

	s.db.ExecContext(ctx, query, status, verifiedAt, time.Now(), domain.ID)
}

func (s *DomainService) processSSLSetup(ctx context.Context, domain *models.DomainConfig) {
	// Implementar configuração automática de SSL usando ACME/Let's Encrypt
	// Por enquanto, simular processo
	time.Sleep(60 * time.Second)

	query := `
		UPDATE domain_configs 
		SET ssl_status = 'active', ssl_provider = 'letsencrypt', updated_at = $1
		WHERE id = $2
	`

	s.db.ExecContext(ctx, query, time.Now(), domain.ID)
}
