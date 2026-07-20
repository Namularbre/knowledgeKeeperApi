package app

import (
	"context"
	"fmt"

	"github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
)

type FetchByPage struct{ Subjects domain.Repository }
type FetchByPageInput struct{ Page, Limit uint64 }

func (uc FetchByPage) Execute(ctx context.Context, in FetchByPageInput) ([]domain.Subject, error) {
	if in.Limit == 0 {
		return []domain.Subject{}, fmt.Errorf("limit is zero")
	}
	return uc.Subjects.FetchByPage(ctx, in.Page, in.Limit)
}
