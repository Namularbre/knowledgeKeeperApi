package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	fileEnv, err := godotenv.Read(".env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("read .env: %w", err)
	}
	getEnv := func(key string) string {
		if value, exists := fileEnv[key]; exists {
			return value
		}
		return os.Getenv(key)
	}

	cfg := Config{
		Port: getEnv("PORT"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			Name:     getEnv("DB_NAME"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET"),
			Issuer: getEnv("JWT_ISSUER"),
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

	accessTTL, err := parseDurationOrDefault(getEnv("JWT_TTL"), 15*time.Minute)
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_TTL: %w", err)
	}
	cfg.JWT.AccessTTL = accessTTL

	refreshTTL, err := parseDurationOrDefault(getEnv("JWT_REFRESH_TTL"), 7*24*time.Hour)
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
