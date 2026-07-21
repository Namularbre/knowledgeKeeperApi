package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
)

type SearchByName struct{ Subjects domain.Repository }
type SearchByNameInput struct{ Name string }

func (uc SearchByName) Execute(ctx context.Context, in SearchByNameInput) ([]domain.Subject, error) {
	return uc.Subjects.SearchByName(ctx, in.Name)
}
