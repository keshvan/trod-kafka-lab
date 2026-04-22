package main

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/config"
	"github.com/keshvan/trod-kafka-lab/data-service/internal/db"
)

func main() {
	cfg := config.MustLoad()

	if err := db.RunMigrations(cfg.DatabaseURL); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("run migrations: %v", err)
	}

	log.Println("migrations completed")
}
