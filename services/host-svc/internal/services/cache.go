package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pagemagic/host-svc/internal/config"
	"pagemagic/host-svc/internal/models"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	config *config.Config
	client *redis.Client
}

func NewCacheService(cfg *config.Config) *CacheService {
	// Parse Redis URL
	opt, err := redis.ParseURL(cfg.Cache.RedisURL)
	if err != nil {
		// Fallback para configuração padrão
		opt = &redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		}
	}

	client := redis.NewClient(opt)

	return &CacheService{
		config: cfg,
		client: client,
	}
}

func (s *CacheService) PurgeCache(ctx context.Context, siteID string) error {
	// Purgar cache do site específico
	pattern := fmt.Sprintf("site:%s:*", siteID)
	
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get cache keys: %w", err)
	}

	if len(keys) > 0 {
		err = s.client.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete cache keys: %w", err)
		}
	}

	// Atualizar status de purge
	statusKey := fmt.Sprintf("site:%s:cache_status", siteID)
	status := models.CacheStatus{
		SiteID:    siteID,
		PurgedAt:  &[]time.Time{time.Now()}[0],
		UpdatedAt: time.Now(),
	}

	statusData, _ := json.Marshal(status)
	s.client.Set(ctx, statusKey, statusData, s.config.Cache.TTL)

	return nil
}

func (s *CacheService) GetCacheStatus(ctx context.Context, siteID string) (*models.CacheStatus, error) {
	statusKey := fmt.Sprintf("site:%s:cache_status", siteID)
	
	data, err := s.client.Get(ctx, statusKey).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache status não existe, criar um novo
			return &models.CacheStatus{
				SiteID:      siteID,
				HitRate:     0.0,
				TotalHits:   0,
				TotalMisses: 0,
				UpdatedAt:   time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("failed to get cache status: %w", err)
	}

	var status models.CacheStatus
	err = json.Unmarshal([]byte(data), &status)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache status: %w", err)
	}

	return &status, nil
}

func (s *CacheService) CacheFile(ctx context.Context, siteID, path string, content []byte, contentType string) error {
	key := fmt.Sprintf("site:%s:file:%s", siteID, path)
	
	fileData := map[string]interface{}{
		"content":      content,
		"content_type": contentType,
		"cached_at":    time.Now(),
	}

	data, err := json.Marshal(fileData)
	if err != nil {
		return fmt.Errorf("failed to marshal file data: %w", err)
	}

	return s.client.Set(ctx, key, data, s.config.Cache.TTL).Err()
}

func (s *CacheService) GetCachedFile(ctx context.Context, siteID, path string) ([]byte, string, error) {
	key := fmt.Sprintf("site:%s:file:%s", siteID, path)
	
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, "", fmt.Errorf("file not in cache")
		}
		return nil, "", fmt.Errorf("failed to get cached file: %w", err)
	}

	var fileData map[string]interface{}
	err = json.Unmarshal([]byte(data), &fileData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal file data: %w", err)
	}

	content, ok := fileData["content"].([]byte)
	if !ok {
		return nil, "", fmt.Errorf("invalid content format in cache")
	}

	contentType, ok := fileData["content_type"].(string)
	if !ok {
		contentType = "application/octet-stream"
	}

	// Atualizar estatísticas de hit
	s.updateCacheStats(ctx, siteID, true)

	return content, contentType, nil
}

func (s *CacheService) updateCacheStats(ctx context.Context, siteID string, hit bool) {
	statusKey := fmt.Sprintf("site:%s:cache_status", siteID)
	
	// Usar pipeline para operações atômicas
	pipe := s.client.Pipeline()
	
	if hit {
		pipe.HIncrBy(ctx, statusKey, "total_hits", 1)
	} else {
		pipe.HIncrBy(ctx, statusKey, "total_misses", 1)
	}
	
	pipe.HSet(ctx, statusKey, "updated_at", time.Now())
	pipe.Expire(ctx, statusKey, s.config.Cache.TTL)
	
	pipe.Exec(ctx)
}

func (s *CacheService) WarmUpCache(ctx context.Context, siteID string, files []models.DeploymentFile) error {
	// Implementar warm-up do cache com arquivos do deployment
	for _, file := range files {
		// Simular carregamento e cache dos arquivos
		// Em implementação real, baixaria do storage e cachear
		cacheKey := fmt.Sprintf("site:%s:file:%s", siteID, file.Path)
		
		fileData := map[string]interface{}{
			"path":         file.Path,
			"content_type": file.ContentType,
			"size":         file.Size,
			"hash":         file.Hash,
			"url":          file.URL,
			"cached_at":    time.Now(),
		}

		data, _ := json.Marshal(fileData)
		s.client.Set(ctx, cacheKey, data, s.config.Cache.TTL)
	}

	return nil
}

func (s *CacheService) Close() error {
	return s.client.Close()
}
