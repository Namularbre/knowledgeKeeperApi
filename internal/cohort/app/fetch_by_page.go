package app

import (
	"context"
	"fmt"

	"github.com/Namularbre/knowledgeKeeperApi/internal/cohort/domain"
)

type FetchByPage struct{ Cohorts domain.Repository }
type FetchByPageInput struct{ Page, Limit uint64 }

func (uc FetchByPage) Execute(ctx context.Context, in FetchByPageInput) ([]domain.Cohort, error) {
	if in.Limit == 0 {
		return []domain.Cohort{}, fmt.Errorf("limit is zero")
	}
	return uc.Cohorts.FetchByPage(ctx, in.Page, in.Limit)
}
