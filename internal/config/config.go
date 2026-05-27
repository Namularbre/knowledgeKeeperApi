package config

import (
	"fmt"
	"os"
	"time"
)

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type JWTConfig struct {
	Secret     string
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Config struct {
	Port string
	DB   DBConfig
	JWT  JWTConfig
}

func LoadFromEnv() (Config, error) {
	cfg := Config{
		Port: os.Getenv("PORT"),
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Issuer: os.Getenv("JWT_ISSUER"),
		},
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.DB.Host == "" || cfg.DB.Port == "" || cfg.DB.Name == "" || cfg.DB.User == "" {
		return Config{}, fmt.Errorf("missing required DB env vars (DB_HOST/DB_PORT/DB_NAME/DB_USER)")
	}

	if cfg.JWT.Secret == "" {
		return Config{}, fmt.Errorf("missing required env var JWT_SECRET")
	}
	if cfg.JWT.Issuer == "" {
		cfg.JWT.Issuer = "knowledgeKeeperApi"
	}

	accessTTL, err := parseDurationOrDefault(os.Getenv("JWT_TTL"), 15*time.Minute)
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_TTL: %w", err)
	}
	cfg.JWT.AccessTTL = accessTTL

	refreshTTL, err := parseDurationOrDefault(os.Getenv("JWT_REFRESH_TTL"), 7*24*time.Hour)
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_REFRESH_TTL: %w", err)
	}
	cfg.JWT.RefreshTTL = refreshTTL

	return cfg, nil
}

func parseDurationOrDefault(value string, fallback time.Duration) (time.Duration, error) {
	if value == "" {
		return fallback, nil
	}
	return time.ParseDuration(value)
}
