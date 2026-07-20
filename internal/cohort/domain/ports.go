package domain

import "context"

// Repository persists and retrieves cohorts.
type Repository interface {
	Create(ctx context.Context, name string) (Cohort, error)
	FetchByPage(ctx context.Context, page, perPage uint64) ([]Cohort, error)
	FindByID(ctx context.Context, id uint64) (Cohort, error)
	SearchByName(ctx context.Context, name string) ([]Cohort, error)
	FindByUserID(ctx context.Context, userID uint64) ([]Cohort, error)
	AddUserCohort(ctx context.Context, cohortID, userID uint64) ([]Cohort, error)
	RemoveUserCohort(ctx context.Context, cohortID, userID uint64) ([]Cohort, error)
}
