package domain

import "context"

// Repository persists and retrieves roles.
type Repository interface {
	Create(ctx context.Context, label string) (Role, error)
	FetchByPage(ctx context.Context, page, perPage uint64) ([]Role, error)
	FindByID(ctx context.Context, id int64) (Role, error)
	SearchByLabel(ctx context.Context, label string) ([]Role, error)
	FindByUserID(ctx context.Context, userID int64) (Role, error)
	AddUserRole(ctx context.Context, roleID int64, userID int64) error
	RemoveUserRole(ctx context.Context, roleID int64, userID int64) error
}
