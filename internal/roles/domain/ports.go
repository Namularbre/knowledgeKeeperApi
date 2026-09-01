package domain

import "context"

// Repository persists and retrieves roles.
type Repository interface {
	Create(ctx context.Context, label RoleLabel) (Role, error)
	FetchByPage(ctx context.Context, page, perPage uint64) ([]Role, error)
	FindByID(ctx context.Context, id uint64) (Role, error)
	SearchByLabel(ctx context.Context, label string) ([]Role, error)
	FindByUserID(ctx context.Context, userID uint64) ([]Role, error)
	AddUserRole(ctx context.Context, roleID, userID uint64) ([]Role, error)
	RemoveUserRole(ctx context.Context, roleID, userID uint64) ([]Role, error)
}
