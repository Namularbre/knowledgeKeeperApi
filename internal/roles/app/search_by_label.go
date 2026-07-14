package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type SearchByLabel struct {
	Roles domain.Repository
}

type SearchByLabelInput struct {
	Label string
}

func (uc SearchByLabel) Execute(ctx context.Context, in SearchByLabelInput) ([]domain.Role, error) {
	return uc.Roles.SearchByLabel(ctx, in.Label)
}
