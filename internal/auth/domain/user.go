// Package domain holds the authentication domain model and its ports.
//
// The domain layer has zero dependencies on infrastructure: no SQL, no HTTP,
// no JWT library. Application use cases orchestrate these primitives, and
// infrastructure adapters implement the interfaces defined here.
package domain

import (
	"errors"
	"time"
)

// User represents an authenticated principal of the API.
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Roles        []Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RefreshToken is the persisted view of an opaque refresh token. Only the
// SHA-256 hash of the token is stored, never the token itself.
type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// IsActive reports whether the refresh token can still be used.
func (t RefreshToken) IsActive(now time.Time) bool {
	return t.RevokedAt == nil && now.Before(t.ExpiresAt)
}

// Role represents the role of a user
type Role struct {
	ID        int64
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Domain errors. Handlers map these to HTTP status codes; use cases return
// them so callers can distinguish business-rule failures from infra errors.
var (
	ErrEmailAlreadyTaken  = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidRefresh     = errors.New("invalid or expired refresh token")
	ErrWeakPassword       = errors.New("password does not meet minimum requirements")
	ErrInvalidEmail       = errors.New("invalid email")
)
