package domain

import (
	"context"
	"time"
)

// UserRepository persists and retrieves users.
type UserRepository interface {
	// Create inserts a new user. Returns ErrEmailAlreadyTaken on email collision.
	Create(ctx context.Context, email, passwordHash string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id int64) (User, error)
}

// RefreshTokenRepository persists and retrieves refresh tokens by their hash.
type RefreshTokenRepository interface {
	// Create stores a new refresh token row.
	Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) (RefreshToken, error)
	// FindByHash returns the token row matching the given hash, or
	// ErrInvalidRefresh if missing.
	FindByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	// Revoke marks the token identified by hash as revoked.
	Revoke(ctx context.Context, tokenHash string) error
}

// PasswordHasher abstracts the password hashing algorithm.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

// TokenPair is the access + refresh token pair returned to clients on login
// and refresh.
type TokenPair struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string // opaque, returned ONCE to the client
	RefreshExpiresAt time.Time
}

// TokenIssuer mints and verifies access tokens, and generates opaque refresh
// tokens. The hashing of refresh tokens for storage is the issuer's
// responsibility so the algorithm stays in one place.
type TokenIssuer interface {
	IssueAccessToken(userID int64, now time.Time) (token string, expiresAt time.Time, err error)
	ParseAccessToken(token string) (userID int64, err error)
	GenerateRefreshToken() (token, hash string, err error)
	HashRefreshToken(token string) string
}

type RoleRepository interface {
	Create(ctx context.Context, label string) (Role, error)
	FetchByPage(ctx context.Context, page, perPage int) ([]Role, error)
	FindById(ctx context.Context, id int64) (Role, error)
	SearchByLabel(ctx context.Context, label string) ([]Role, error)
	FindByUserID(ctx context.Context, userID int64) (Role, error)
	AddUserRole(ctx context.Context, roleID int64, userID int64) error
	RemoveUserRole(ctx context.Context, roleID int64, userID int64) error
}
