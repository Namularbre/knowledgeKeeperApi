package app

import (
	"context"
	"strings"

	"github.com/Namularbre/knowledgeKeeperApi/internal/cohort/domain"
)

type CreateCohort struct{ Cohorts domain.Repository }
type CreateCohortInput struct{ Name string }

func (uc CreateCohort) Execute(ctx context.Context, in CreateCohortInput) (domain.Cohort, error) {
	name := strings.ToLower(strings.TrimSpace(in.Name))
	if name == "" {
		return domain.Cohort{}, domain.ErrInvalidCohortName
	}
	return uc.Cohorts.Create(ctx, name)
}
