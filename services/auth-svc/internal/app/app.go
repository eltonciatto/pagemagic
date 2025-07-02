package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pagemagic/auth-svc/internal/config"
	"pagemagic/auth-svc/internal/handlers"
	"pagemagic/auth-svc/internal/repository"
	"pagemagic/auth-svc/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	config *config.Config
	server *http.Server
}

func New(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Run() error {
	// Configurar modo do Gin
	if a.config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar repositório
	repo, err := repository.New(a.config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}
	defer repo.Close()

	// Inicializar serviços
	authService := services.NewAuthService(repo, a.config)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Configurar rotas
	router := a.setupRoutes(authHandler)

	// Configurar servidor HTTP
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Port),
		Handler: router,
	}

	// Canal para capturar sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Auth service starting on port %s", a.config.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
	return nil
}

func (a *App) setupRoutes(authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.Default()

	// Middleware de CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "auth-svc"})
	})

	// Métricas do Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Grupo de rotas da API
	api := router.Group("/api/v1")
	{
		// Rotas de autenticação
		auth := api.Group("/auth")
		{
			auth.POST("/magic-link", authHandler.SendMagicLink)
			auth.POST("/verify", authHandler.VerifyMagicLink)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Rotas protegidas
		protected := api.Group("/")
		protected.Use(authHandler.AuthMiddleware())
		{
			protected.GET("/profile", authHandler.GetProfile)
			protected.PUT("/profile", authHandler.UpdateProfile)
		}
	}

	return router
}
