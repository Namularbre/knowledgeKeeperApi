package infra

//TODO: make the funct

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

// TODO: move this to a file that is not related to a specific repository, as all repo can have a use for it
const mysqlDuplicateEntryErrno = 1062

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Create(ctx context.Context, email, passwordHash string) (domain.User, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash) VALUES (?, ?)`,
		email, passwordHash,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrno {
			return domain.User{}, domain.ErrEmailAlreadyTaken
		}
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, fmt.Errorf("user last insert id: %w", err)
	}
	return r.FindByID(ctx, id)
}

func (r *MySQLUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.findOne(ctx,
		`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = ?`,
		strings.ToLower(strings.TrimSpace(email)),
	)
}

func (r *MySQLUserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	return r.findOne(ctx,
		`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = ?`,
		id,
	)
}

func (r *MySQLUserRepository) findOne(ctx context.Context, query string, arg any) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, arg).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("query user: %w", err)
	}
	return u, nil
}

type MySQLRefreshTokenRepository struct {
	db *sql.DB
}

func NewMySQLRefreshTokenRepository(db *sql.DB) *MySQLRefreshTokenRepository {
	return &MySQLRefreshTokenRepository{db: db}
}

func (r *MySQLRefreshTokenRepository) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) (domain.RefreshToken, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)`,
		userID, tokenHash, expiresAt.UTC(),
	)
	if err != nil {
		return domain.RefreshToken{}, fmt.Errorf("insert refresh token: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.RefreshToken{}, fmt.Errorf("refresh token last insert id: %w", err)
	}
	return domain.RefreshToken{
		ID:        id,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt.UTC(),
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (r *MySQLRefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	var t domain.RefreshToken
	var revokedAt sql.NullTime
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		 FROM refresh_tokens WHERE token_hash = ?`,
		tokenHash,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &revokedAt, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.RefreshToken{}, domain.ErrInvalidRefresh
	}
	if err != nil {
		return domain.RefreshToken{}, fmt.Errorf("query refresh token: %w", err)
	}
	if revokedAt.Valid {
		t.RevokedAt = &revokedAt.Time
	}
	return t, nil
}

func (r *MySQLRefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = CURRENT_TIMESTAMP
		 WHERE token_hash = ? AND revoked_at IS NULL`,
		tokenHash,
	)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}
