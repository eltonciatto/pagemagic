package services

import (
	"context"
	"time"

	"pagemagic/host-svc/internal/config"
	"pagemagic/host-svc/internal/models"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsService struct {
	config *config.Config

	// Prometheus metrics
	siteRequests *prometheus.CounterVec
	responseTime *prometheus.HistogramVec
	activeSites  prometheus.Gauge
	deployments  *prometheus.CounterVec
	cacheHitRate *prometheus.GaugeVec
}

func NewMetricsService(cfg *config.Config) *MetricsService {
	s := &MetricsService{
		config: cfg,
	}

	if cfg.Metrics.Enabled {
		s.initPrometheusMetrics()
	}

	return s
}

func (s *MetricsService) initPrometheusMetrics() {
	s.siteRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pagemagic_site_requests_total",
			Help: "Total number of requests to hosted sites",
		},
		[]string{"site_id", "domain", "status_code"},
	)

	s.responseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pagemagic_site_response_time_seconds",
			Help:    "Response time for hosted sites",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"site_id", "domain"},
	)

	s.activeSites = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "pagemagic_active_sites_total",
			Help: "Total number of active hosted sites",
		},
	)

	s.deployments = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pagemagic_deployments_total",
			Help: "Total number of deployments",
		},
		[]string{"site_id", "status"},
	)

	s.cacheHitRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pagemagic_cache_hit_rate",
			Help: "Cache hit rate for hosted sites",
		},
		[]string{"site_id"},
	)

	// Registrar métricas
	prometheus.MustRegister(s.siteRequests)
	prometheus.MustRegister(s.responseTime)
	prometheus.MustRegister(s.activeSites)
	prometheus.MustRegister(s.deployments)
	prometheus.MustRegister(s.cacheHitRate)
}

func (s *MetricsService) RecordSiteRequest(siteID, domain, statusCode string, duration time.Duration) {
	if s.config.Metrics.Enabled {
		s.siteRequests.WithLabelValues(siteID, domain, statusCode).Inc()
		s.responseTime.WithLabelValues(siteID, domain).Observe(duration.Seconds())
	}
}

func (s *MetricsService) RecordDeployment(siteID, status string) {
	if s.config.Metrics.Enabled {
		s.deployments.WithLabelValues(siteID, status).Inc()
	}
}

func (s *MetricsService) UpdateActiveSites(count float64) {
	if s.config.Metrics.Enabled {
		s.activeSites.Set(count)
	}
}

func (s *MetricsService) UpdateCacheHitRate(siteID string, hitRate float64) {
	if s.config.Metrics.Enabled {
		s.cacheHitRate.WithLabelValues(siteID).Set(hitRate)
	}
}

func (s *MetricsService) GetAnalytics(ctx context.Context, siteID string, startDate, endDate time.Time) (*models.TrafficStats, error) {
	// Implementar coleta de analytics
	// Por enquanto, retornar dados mock
	return &models.TrafficStats{
		SiteID:       siteID,
		Date:         startDate,
		PageViews:    1000,
		UniqueVisits: 750,
		Bandwidth:    1024 * 1024 * 100, // 100MB
		Requests:     1500,
		Countries: map[string]int64{
			"BR": 400,
			"US": 300,
			"CA": 200,
			"MX": 100,
		},
		Referrers: map[string]int64{
			"google.com":   500,
			"facebook.com": 200,
			"direct":       300,
		},
		Pages: map[string]int64{
			"/":        600,
			"/about":   200,
			"/contact": 100,
			"/blog":    100,
		},
		StatusCodes: map[string]int64{
			"200": 1400,
			"404": 80,
			"500": 20,
		},
	}, nil
}

func (s *MetricsService) GetStats(ctx context.Context, siteID string) (*models.StatsResponse, error) {
	// Implementar coleta de estatísticas gerais
	// Combinar dados de analytics, cache, etc.

	analytics, err := s.GetAnalytics(ctx, siteID, time.Now().AddDate(0, 0, -7), time.Now())
	if err != nil {
		return nil, err
	}

	cacheStatus := &models.CacheStatus{
		SiteID:      siteID,
		HitRate:     0.85,
		TotalHits:   1200,
		TotalMisses: 200,
		UpdatedAt:   time.Now(),
	}

	return &models.StatsResponse{
		Traffic: analytics,
		Cache:   cacheStatus,
	}, nil
}

func (s *MetricsService) PrometheusHandler() promhttp.Handler {
	return promhttp.Handler()
}
