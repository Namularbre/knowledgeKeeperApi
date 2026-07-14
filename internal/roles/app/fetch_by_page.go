package app

import (
	"context"
	"fmt"

	"github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type FetchByPage struct {
	Roles domain.Repository
}

type FetchByPageInput struct {
	Page  uint64
	Limit uint64
}

func (uc FetchByPage) Execute(ctx context.Context, in FetchByPageInput) ([]domain.Role, error) {
	if in.Limit == 0 {
		return []domain.Role{}, fmt.Errorf("limit is zero")
	}

	return uc.Roles.FetchByPage(ctx, in.Page, in.Limit)
}
