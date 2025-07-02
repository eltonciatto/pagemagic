package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"pagemagic/build-svc/internal/config"
	"pagemagic/build-svc/internal/handlers"
	"pagemagic/build-svc/internal/services"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type App struct {
	config      *config.Config
	server      *http.Server
	db          *sql.DB
	buildSvc    *services.BuildService
	templateSvc *services.TemplateService
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
	}

	fmt.Printf("Build service starting on %s\n", a.server.Addr)
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
	// Serviço de build
	a.buildSvc = services.NewBuildService(a.db, a.config)

	// Serviço de templates
	a.templateSvc = services.NewTemplateService(a.db, a.config)

	return nil
}

func (a *App) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Handlers
	buildHandler := handlers.NewBuildHandler(a.buildSvc)
	templateHandler := handlers.NewTemplateHandler(a.templateSvc)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Build jobs
	api.HandleFunc("/builds", buildHandler.CreateBuildJob).Methods("POST")
	api.HandleFunc("/builds/{jobId}", buildHandler.GetBuildJob).Methods("GET")
	api.HandleFunc("/builds/{jobId}/cancel", buildHandler.CancelBuildJob).Methods("POST")
	api.HandleFunc("/builds/{jobId}/retry", buildHandler.RetryBuildJob).Methods("POST")
	api.HandleFunc("/builds/{jobId}/logs", buildHandler.GetBuildLogs).Methods("GET")
	api.HandleFunc("/sites/{siteId}/builds", buildHandler.ListBuildJobs).Methods("GET")

	// Templates
	api.HandleFunc("/templates", templateHandler.ListTemplates).Methods("GET")
	api.HandleFunc("/templates/featured", templateHandler.GetFeaturedTemplates).Methods("GET")
	api.HandleFunc("/templates/search", templateHandler.SearchTemplates).Methods("GET")
	api.HandleFunc("/templates/categories", templateHandler.GetCategories).Methods("GET")
	api.HandleFunc("/templates/frameworks", templateHandler.GetFrameworks).Methods("GET")
	api.HandleFunc("/templates/{id}", templateHandler.GetTemplate).Methods("GET")
	api.HandleFunc("/templates/{id}/download", templateHandler.DownloadTemplate).Methods("POST")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
