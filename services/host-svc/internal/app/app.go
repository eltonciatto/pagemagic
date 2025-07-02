package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"pagemagic/host-svc/internal/config"
	"pagemagic/host-svc/internal/handlers"
	"pagemagic/host-svc/internal/services"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type App struct {
	config     *config.Config
	server     *http.Server
	db         *sql.DB
	hostSvc    *services.HostingService
	deploySvc  *services.DeploymentService
	domainSvc  *services.DomainService
	cacheSvc   *services.CacheService
	metricsSvc *services.MetricsService
}

func New(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Run() error {
	// Conectar ao banco de dados
	if err := a.setupDatabase(); err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Inicializar serviços
	if err := a.setupServices(); err != nil {
		return fmt.Errorf("failed to setup services: %w", err)
	}

	// Configurar rotas
	router := a.setupRoutes()

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Configurar servidor HTTP
	a.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port),
		Handler:      c.Handler(router),
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
		IdleTimeout:  a.config.Server.IdleTimeout,
	}

	fmt.Printf("Host service starting on %s\n", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			return err
		}
	}

	if a.db != nil {
		return a.db.Close()
	}

	return nil
}

func (a *App) setupDatabase() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		a.config.Database.Host,
		a.config.Database.Port,
		a.config.Database.User,
		a.config.Database.Password,
		a.config.Database.Name,
		a.config.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	// Configurações de pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	a.db = db
	return nil
}

func (a *App) setupServices() error {
	// Serviço de hosting
	a.hostSvc = services.NewHostingService(a.db, a.config)

	// Serviço de deployment
	a.deploySvc = services.NewDeploymentService(a.db, a.config)

	// Serviço de domínios
	a.domainSvc = services.NewDomainService(a.db, a.config)

	// Serviço de cache
	a.cacheSvc = services.NewCacheService(a.config)

	// Serviço de métricas
	a.metricsSvc = services.NewMetricsService(a.config)

	return nil
}

func (a *App) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Handlers
	hostHandler := handlers.NewHostingHandler(a.hostSvc, a.deploySvc, a.domainSvc)
	deployHandler := handlers.NewDeploymentHandler(a.deploySvc, a.cacheSvc)
	domainHandler := handlers.NewDomainHandler(a.domainSvc)
	metricsHandler := handlers.NewMetricsHandler(a.metricsSvc)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Sites
	api.HandleFunc("/sites", hostHandler.CreateSite).Methods("POST")
	api.HandleFunc("/sites", hostHandler.ListSites).Methods("GET")
	api.HandleFunc("/sites/{id}", hostHandler.GetSite).Methods("GET")
	api.HandleFunc("/sites/{id}", hostHandler.UpdateSite).Methods("PUT")
	api.HandleFunc("/sites/{id}", hostHandler.DeleteSite).Methods("DELETE")
	api.HandleFunc("/sites/{id}/status", hostHandler.GetSiteStatus).Methods("GET")

	// Deployments
	api.HandleFunc("/sites/{id}/deployments", deployHandler.CreateDeployment).Methods("POST")
	api.HandleFunc("/sites/{id}/deployments", deployHandler.ListDeployments).Methods("GET")
	api.HandleFunc("/sites/{id}/deployments/{deployId}", deployHandler.GetDeployment).Methods("GET")
	api.HandleFunc("/sites/{id}/deployments/{deployId}/rollback", deployHandler.RollbackDeployment).Methods("POST")

	// Domínios
	api.HandleFunc("/sites/{id}/domains", domainHandler.AddDomain).Methods("POST")
	api.HandleFunc("/sites/{id}/domains", domainHandler.ListDomains).Methods("GET")
	api.HandleFunc("/sites/{id}/domains/{domainId}", domainHandler.UpdateDomain).Methods("PUT")
	api.HandleFunc("/sites/{id}/domains/{domainId}", domainHandler.DeleteDomain).Methods("DELETE")
	api.HandleFunc("/sites/{id}/domains/{domainId}/verify", domainHandler.VerifyDomain).Methods("POST")

	// Cache
	api.HandleFunc("/sites/{id}/cache/purge", deployHandler.PurgeCache).Methods("POST")
	api.HandleFunc("/sites/{id}/cache/status", deployHandler.GetCacheStatus).Methods("GET")

	// Analytics e métricas
	api.HandleFunc("/sites/{id}/analytics", metricsHandler.GetAnalytics).Methods("GET")
	api.HandleFunc("/sites/{id}/stats", metricsHandler.GetStats).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Prometheus metrics
	if a.config.Metrics.Enabled {
		router.Handle(a.config.Metrics.Path, metricsHandler.PrometheusHandler())
	}

	// Servidor de arquivos estáticos (sites hospedados)
	router.PathPrefix("/").Handler(hostHandler.ServeStaticFiles())

	return router
}
