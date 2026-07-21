package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
)

type FindByUserID struct{ Subjects domain.Repository }
type FindByUserIDInput struct{ ID uint64 }

func (uc FindByUserID) Execute(ctx context.Context, in FindByUserIDInput) ([]domain.Subject, error) {
	return uc.Subjects.FindByUserID(ctx, in.ID)
}
