package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
)

type FindByID struct{ Subjects domain.Repository }
type FindByIDInput struct{ ID uint64 }

func (uc FindByID) Execute(ctx context.Context, in FindByIDInput) (domain.Subject, error) {
	return uc.Subjects.FindByID(ctx, in.ID)
}
