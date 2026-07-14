package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type RemoveUserRole struct {
	Roles domain.Repository
}

type RemoveUserRoleInput struct {
	RoleID uint64
	UserID uint64
}

func (uc *RemoveUserRole) Execute(ctx context.Context, in RemoveUserRoleInput) ([]domain.Role, error) {
	return uc.Roles.RemoveUserRole(ctx, in.RoleID, in.UserID)
}
