package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"pagemagic/host-svc/internal/app"
	"pagemagic/host-svc/internal/config"

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

	// Inicializar aplicação
	application := app.New(cfg)

	// Canal para capturar sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Executar aplicação em goroutine
	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("Failed to run application: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Shutting down host service...")

	// Graceful shutdown
	if err := application.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown application: %v", err)
	}

	log.Println("Host service exited")
}
