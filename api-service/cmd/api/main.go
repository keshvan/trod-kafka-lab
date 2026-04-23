package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/keshvan/trod-kafka-lab/api-service/internal/app"
	"github.com/keshvan/trod-kafka-lab/api-service/internal/config"
)

func main() {
	cfg := config.MustLoad()

	apiApp, err := app.New(cfg)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	apiApp.HTTPServer.Start()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiApp.HTTPServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown http server: %v", err)
	}

	if err := apiApp.Shutdown(); err != nil {
		log.Printf("shutdown dependencies: %v", err)
	}
}
