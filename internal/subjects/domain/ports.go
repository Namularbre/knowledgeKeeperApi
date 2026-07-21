package domain

import "context"

// Repository persists and retrieves subjects.
type Repository interface {
	Create(ctx context.Context, name string) (Subject, error)
	FetchByPage(ctx context.Context, page, perPage uint64) ([]Subject, error)
	FindByID(ctx context.Context, id uint64) (Subject, error)
	SearchByName(ctx context.Context, name string) ([]Subject, error)
	FindByUserID(ctx context.Context, userID uint64) ([]Subject, error)
	AddUserSubject(ctx context.Context, subjectID, userID uint64) ([]Subject, error)
	RemoveUserSubject(ctx context.Context, subjectID, userID uint64) ([]Subject, error)
}
