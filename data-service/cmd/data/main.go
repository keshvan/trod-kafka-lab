package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/keshvan/trod-kafka-lab/data-service/internal/app"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/config"
)

func main() {
	cfg := config.MustLoad()

	dataApp, err := app.New(cfg)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	consumerErrCh := make(chan error, 1)

	dataApp.HTTPServer.Start()

	go func() {
		consumerErrCh <- dataApp.Consumer.Start(ctx)
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-consumerErrCh:
		if err != nil {
			log.Printf("kafka consumer stopped with error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dataApp.HTTPServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown http server: %v", err)
	}

	if err := dataApp.Shutdown(); err != nil {
		log.Printf("shutdown dependencies: %v", err)
	}
}
