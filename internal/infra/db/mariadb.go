package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MariaDB struct {
	db *sql.DB
}

func NewMariaDB(host, port, name, user, password string) (*MariaDB, error) {
	// Format DSN (go-sql-driver/mysql)
	// parseTim
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_general_ci",
		user, password, host, port, name,
	)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Pooling (à ajuster selon charge)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return &MariaDB{db: sqlDB}, nil
}

func (m *MariaDB) Ping(ctx context.Context) error {
	return m.db.PingContext(ctx)
}

func (m *MariaDB) Close() error {
	return m.db.Close()
}

// DB expose *sql.DB aux implémentations infra (repositories SQL).
// Évite de l’exposer au domain/application.
func (m *MariaDB) DB() *sql.DB {
	return m.db
}

// ApplySchema executes a SQL script composed of one or more statements
// separated by `;`. It is intended for idempotent bootstrap scripts
// (CREATE TABLE IF NOT EXISTS, ...). Statements are executed sequentially
// against the underlying connection.
//
// The go-sql-driver/mysql connector does not support multi-statement queries
// by default, hence the manual split.
func (m *MariaDB) ApplySchema(ctx context.Context, script string) error {
	for _, raw := range strings.Split(script, ";") {
		stmt := strings.TrimSpace(raw)
		if stmt == "" {
			continue
		}
		if _, err := m.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("apply schema: %w (statement: %q)", err, truncate(stmt, 80))
		}
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
