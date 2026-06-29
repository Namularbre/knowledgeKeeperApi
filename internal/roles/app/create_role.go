package app

import (
	"context"
	"strings"

	"github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type CreateRole struct {
	Roles domain.Repository
}

type CreateRoleInput struct {
	Label string
}

func (uc CreateRole) Execute(ctx context.Context, in CreateRoleInput) (domain.Role, error) {
	label, err := normalizeLabel(in.Label)
	if err != nil {
		return domain.Role{}, err
	}
	if len(in.Label) == 0 {
		return domain.Role{}, domain.ErrInvalidRoleLabel
	}

	return uc.Roles.Create(ctx, label)
}

// normalizeLabel aims to make roles labels look clean with all letters in lower and by trimming the input string
func normalizeLabel(label string) (string, error) {
	trimmed := strings.ToLower(strings.TrimSpace(label))
	if trimmed == "" {
		return "", domain.ErrInvalidRoleLabel
	}
	return trimmed, nil
}
