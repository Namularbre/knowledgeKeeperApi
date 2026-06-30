package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type FindByID struct {
	Roles domain.Repository
}

type FindByIDInput struct {
	ID uint64
}

func (uc *FindByID) Execute(ctx context.Context, in FindByIDInput) (domain.Role, error) {
	return uc.Roles.FindByID(ctx, in.ID)
}
