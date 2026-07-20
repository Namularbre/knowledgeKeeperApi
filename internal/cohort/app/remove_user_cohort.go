package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/cohort/domain"
)

type RemoveUserCohort struct{ Cohorts domain.Repository }
type RemoveUserCohortInput struct{ UserID, CohortID uint64 }

func (uc RemoveUserCohort) Execute(ctx context.Context, in RemoveUserCohortInput) ([]domain.Cohort, error) {
	return uc.Cohorts.RemoveUserCohort(ctx, in.CohortID, in.UserID)
}
