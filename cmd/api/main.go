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

	_ "github.com/Namularbre/knowledgeKeeperApi/docs"
	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/app"
	authinfra "github.com/Namularbre/knowledgeKeeperApi/internal/auth/infra"
	authhttp "github.com/Namularbre/knowledgeKeeperApi/internal/auth/infra/http"
	authsql "github.com/Namularbre/knowledgeKeeperApi/internal/auth/infra/sql"
	"github.com/Namularbre/knowledgeKeeperApi/internal/config"
	"github.com/Namularbre/knowledgeKeeperApi/internal/infra/db"
	httpserver "github.com/Namularbre/knowledgeKeeperApi/internal/infra/http"
)

// @title           knowledgeKeeperApi
// @version         1.0.0
// @description     Personal knowledge keeper API. Auth (register/login/refresh), and metadata endpoints.
// @BasePath        /
//
// @securityDefinitions.apikey BearerAuth
// @in   header
// @name Authorization
// @description Type "Bearer {token}" where {token} is the access_token returned by /auth/login.
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

	bootCtx, bootCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer bootCancel()

	if err := maria.Ping(bootCtx); err != nil {
		log.Fatalf("db ping error: %v", err)
	}
	log.Println("DB connection OK")

	if err := maria.ApplySchema(bootCtx, authsql.Schema); err != nil {
		log.Fatalf("schema apply error: %v", err)
	}
	log.Println("Auth schema applied")

	users := authinfra.NewMySQLUserRepository(maria.DB())
	refreshes := authinfra.NewMySQLRefreshTokenRepository(maria.DB())
	hasher := authinfra.NewBcryptHasher(0)
	issuer := authinfra.NewJWTIssuer(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	handlers := authhttp.Handlers{
		Register: authhttp.RegisterHandler{UC: app.RegisterUser{Users: users, Hasher: hasher}},
		Login: authhttp.LoginHandler{UC: app.LoginUser{
			Users:         users,
			RefreshTokens: refreshes,
			Hasher:        hasher,
			Tokens:        issuer,
			RefreshTTL:    cfg.JWT.RefreshTTL,
		}},
		Refresh: authhttp.RefreshHandler{UC: app.RefreshSession{
			Users:         users,
			RefreshTokens: refreshes,
			Tokens:        issuer,
			RefreshTTL:    cfg.JWT.RefreshTTL,
		}},
	}

	server := httpserver.NewServer(cfg.Port)
	server.RegisterRoutes(func(mux *http.ServeMux) {
		mux.Handle("/auth/register", handlers.Register)
		mux.Handle("/auth/login", handlers.Login)
		mux.Handle("/auth/refresh", handlers.Refresh)
	})

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
