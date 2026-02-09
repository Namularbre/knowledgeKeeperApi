package db

import (
	"context"
	"database/sql"
	"fmt"
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
