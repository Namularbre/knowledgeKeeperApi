// Package domain holds the role domain model and its ports.
package domain

import (
	"errors"
	"time"
)

// Role represents a role that can be assigned to users.
type Role struct {
	ID        int64
	Label     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Domain errors related to roles.
var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")
	ErrInvalidRoleLabel  = errors.New("invalid role label")
)
