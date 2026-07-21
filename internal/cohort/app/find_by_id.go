package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/cohort/domain"
)

type FindByID struct{ Cohorts domain.Repository }
type FindByIDInput struct{ ID uint64 }

func (uc FindByID) Execute(ctx context.Context, in FindByIDInput) (domain.Cohort, error) {
	return uc.Cohorts.FindByID(ctx, in.ID)
}
