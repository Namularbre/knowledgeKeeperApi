// Package domain holds the subject domain model and its ports.
package domain

import (
	"errors"
	"time"
)

// Subject represents a subject that can be assigned to users.
type Subject struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrSubjectNotFound      = errors.New("subject not found")
	ErrSubjectAlreadyExists = errors.New("subject already exists")
	ErrInvalidSubjectName   = errors.New("invalid subject name")
	ErrUserSubjectExists    = errors.New("user subject already exists")
)
