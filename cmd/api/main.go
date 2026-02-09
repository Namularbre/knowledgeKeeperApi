package main

import (
	"context"
	"log"
	"time"

	"knowledgeKeeperApi/internal/config"
	"knowledgeKeeperApi/internal/infra/db"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	maria, err := db.NewMariaDB(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.User,
		cfg.DB.Password,
	)
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer func() {
		_ = maria.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := maria.Ping(ctx); err != nil {
		log.Fatalf("db ping error: %v", err)
	}

	log.Println("DB connection OK")
	// Étape suivante : démarrer serveur HTTP + injecter des repositories/use-cases.
}
