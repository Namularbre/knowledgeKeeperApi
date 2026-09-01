// Package domain holds the role domain model and its ports.
package domain

import (
	"errors"
	"time"
)

// Role represents a role that can be assigned to users.
type RoleLabel string

const (
	RoleProf  RoleLabel = "prof"
	RoleAdmin RoleLabel = "admin"
)

// IsValid reports whether the label is one of the roles persisted by the
// application. A user with no role has no users_roles association; it is not
// represented by a RoleLabel value.
func (l RoleLabel) IsValid() bool {
	switch l {
	case RoleProf, RoleAdmin:
		return true
	default:
		return false
	}
}

type Role struct {
	ID        int64
	Label     RoleLabel
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Domain errors related to roles.
var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")
	ErrInvalidRoleLabel  = errors.New("invalid role label")
)
