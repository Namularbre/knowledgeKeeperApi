package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Config struct {
	Port string
	DB   DBConfig
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
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.DB.Host == "" || cfg.DB.Port == "" || cfg.DB.Name == "" || cfg.DB.User == "" {
		return Config{}, fmt.Errorf("missing required DB env vars (DB_HOST/DB_PORT/DB_NAME/DB_USER)")
	}

	return cfg, nil
}
