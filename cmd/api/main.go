package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Namularbre/knowledgeKeeperApi/internal/config"
	"github.com/Namularbre/knowledgeKeeperApi/internal/infra/db"

	httpserver "github.com/Namularbre/knowledgeKeeperApi/internal/infra/http"
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

	server := httpserver.NewServer(cfg.Port)
	server.RegisterRoutes()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
