package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/cohort/domain"
)

type SearchByName struct{ Cohorts domain.Repository }
type SearchByNameInput struct{ Name string }

func (uc SearchByName) Execute(ctx context.Context, in SearchByNameInput) ([]domain.Cohort, error) {
	return uc.Cohorts.SearchByName(ctx, in.Name)
}
