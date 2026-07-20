package app

import (
	"context"
	"strings"

	"github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
)

type CreateSubject struct{ Subjects domain.Repository }
type CreateSubjectInput struct{ Name string }

func (uc CreateSubject) Execute(ctx context.Context, in CreateSubjectInput) (domain.Subject, error) {
	name := strings.ToLower(strings.TrimSpace(in.Name))
	if name == "" {
		return domain.Subject{}, domain.ErrInvalidSubjectName
	}
	return uc.Subjects.Create(ctx, name)
}
