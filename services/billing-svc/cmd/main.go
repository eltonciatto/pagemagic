package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"pagemagic/billing-svc/internal/config"
	"pagemagic/billing-svc/internal/handlers"
	"pagemagic/billing-svc/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize services
	billingService := services.NewBillingService(db)

	// Initialize handlers
	billingHandler := handlers.NewBillingHandler(billingService, cfg.StripeWebhookSecret)

	// Setup router
	router := setupRouter(billingHandler)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting billing service on %s", addr)
	log.Fatal(router.Run(addr))
}

func setupRouter(billingHandler *handlers.BillingHandler) *gin.Engine {
	// Set gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", billingHandler.HealthCheck)

	// API routes
	api := router.Group("/v1")
	{
		api.POST("/customers", billingHandler.CreateCustomer)
		api.POST("/subscriptions", billingHandler.CreateSubscription)
		api.GET("/users/:user_id/subscriptions", billingHandler.GetUserSubscriptions)
	}

	// Webhook routes
	webhooks := router.Group("/webhooks")
	{
		webhooks.POST("/stripe", billingHandler.HandleWebhook)
	}

	return router
}
