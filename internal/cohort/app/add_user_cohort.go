package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/cohort/domain"
)

type AddUserCohort struct {
	Cohort domain.Repository
}

type AddUserCohortInput struct {
	UserID   uint64
	CohortID uint64
}

func (uc *AddUserCohort) Execute(ctx context.Context, in AddUserCohortInput) ([]domain.Cohort, error) {
	return uc.Cohort.AddUserCohort(ctx, in.CohortID, in.UserID)
}
