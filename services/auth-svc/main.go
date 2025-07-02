package main

import (
	"log"

	"pagemagic/auth-svc/internal/app"
	"pagemagic/auth-svc/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Carregar configuração
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Inicializar e executar a aplicação
	application := app.New(cfg)
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
