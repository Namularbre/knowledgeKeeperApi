package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/cohort/domain"
)

type FindByUserID struct{ Cohorts domain.Repository }
type FindByUserIDInput struct{ ID uint64 }

func (uc FindByUserID) Execute(ctx context.Context, in FindByUserIDInput) ([]domain.Cohort, error) {
	return uc.Cohorts.FindByUserID(ctx, in.ID)
}
