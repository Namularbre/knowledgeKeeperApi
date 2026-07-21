package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
)

type AddUserSubject struct{ Subjects domain.Repository }
type AddUserSubjectInput struct{ UserID, SubjectID uint64 }

func (uc AddUserSubject) Execute(ctx context.Context, in AddUserSubjectInput) ([]domain.Subject, error) {
	return uc.Subjects.AddUserSubject(ctx, in.SubjectID, in.UserID)
}
