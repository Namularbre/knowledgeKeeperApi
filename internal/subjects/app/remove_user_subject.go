package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
)

type RemoveUserSubject struct{ Subjects domain.Repository }
type RemoveUserSubjectInput struct{ UserID, SubjectID uint64 }

func (uc RemoveUserSubject) Execute(ctx context.Context, in RemoveUserSubjectInput) ([]domain.Subject, error) {
	return uc.Subjects.RemoveUserSubject(ctx, in.SubjectID, in.UserID)
}
