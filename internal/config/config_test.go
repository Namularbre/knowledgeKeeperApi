package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoadFromEnv(t *testing.T) {
	for _, tt := range []struct {
		name    string
		file    *string
		wantErr string
	}{
		{name: "OS environment without file"},
		{name: "file values take priority with OS fallback", file: envText("PORT=9090\nDB_HOST=file-host\nDB_PASSWORD=\"\"\nJWT_TTL=30m\n")},
		{name: "invalid file", file: envText("DB_HOST=\"unterminated\n"), wantErr: "read .env:"},
		{name: "empty required file value", file: envText("DB_HOST=\n"), wantErr: "missing required DB env vars"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			for key, value := range map[string]string{
				"PORT": "8081", "DB_HOST": "os-host", "DB_PORT": "3306",
				"DB_NAME": "knowledge", "DB_USER": "app", "DB_PASSWORD": "os-password",
				"JWT_SECRET": "test-secret", "JWT_ISSUER": "test-issuer",
				"JWT_TTL": "15m", "JWT_REFRESH_TTL": "168h",
			} {
				t.Setenv(key, value)
			}
			if tt.file != nil {
				if err := os.WriteFile(".env", []byte(*tt.file), 0600); err != nil {
					t.Fatal(err)
				}
			}
			cfg, err := LoadFromEnv()
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := Config{
				Port: "8081",
				DB:   DBConfig{Host: "os-host", Port: "3306", Name: "knowledge", User: "app", Password: "os-password"},
				JWT:  JWTConfig{Secret: "test-secret", Issuer: "test-issuer", AccessTTL: 15 * time.Minute, RefreshTTL: 168 * time.Hour},
			}
			if tt.file != nil {
				want.Port = "9090"
				want.DB.Host = "file-host"
				want.DB.Password = ""
				want.JWT.AccessTTL = 30 * time.Minute
			}
			if cfg != want {
				t.Error("configuration did not match expected file and OS values")
			}
			if os.Getenv("DB_HOST") != "os-host" {
				t.Error("loading configuration modified the OS environment")
			}
		})
	}
}

func envText(value string) *string { return &value }
