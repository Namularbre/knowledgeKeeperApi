package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type AddUserRole struct {
	Roles domain.Repository
}

type AddUserRoleInput struct {
	UserID uint64
	RoleID uint64
}

func (uc *AddUserRole) Execute(ctx context.Context, in AddUserRoleInput) ([]domain.Role, error) {
	return uc.Roles.AddUserRole(ctx, in.RoleID, in.UserID)
}
