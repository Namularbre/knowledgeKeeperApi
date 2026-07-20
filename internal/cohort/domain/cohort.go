// Package domain holds the cohort domain model and its ports.
package domain

import (
	"errors"
	"time"

	"github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
)

// Cohort represents a class that can be assigned to users.
type Cohort struct {
	ID        int64
	Name      string
	Subjects  []domain.Subject
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Domain errors related to cohorts.
var (
	ErrCohortNotFound      = errors.New("cohort not found")
	ErrCohortAlreadyExists = errors.New("cohort already exists")
	ErrInvalidCohortName   = errors.New("invalid cohort name")
	ErrUserCohortExists    = errors.New("user cohort already exists")
)
