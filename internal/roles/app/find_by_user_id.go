package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type FindByUserID struct {
	Roles domain.Repository
}

type FindByUserIDInput struct {
	ID uint64
}

func (uc *FindByUserID) Execute(ctx context.Context, in FindByUserIDInput) ([]domain.Role, error) {
	return uc.Roles.FindByUserID(ctx, in.ID)
}
