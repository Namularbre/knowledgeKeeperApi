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

	return uc.Roles.Create(ctx, label)
}

// normalizeLabel trims and lowercases the input, then validates it against the
// role labels supported by the domain.
func normalizeLabel(label string) (domain.RoleLabel, error) {
	normalized := domain.RoleLabel(strings.ToLower(strings.TrimSpace(label)))
	if !normalized.IsValid() {
		return "", domain.ErrInvalidRoleLabel
	}
	return normalized, nil
}
