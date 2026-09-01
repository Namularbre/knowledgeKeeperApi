package app

import (
	"errors"
	"testing"

	"github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

func TestNormalizeLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    domain.RoleLabel
		wantErr bool
	}{
		{name: "prof", input: "prof", want: domain.RoleProf},
		{name: "admin is normalized", input: " ADMIN ", want: domain.RoleAdmin},
		{name: "teacher is rejected", input: "teacher", wantErr: true},
		{name: "no role is rejected", input: "pas de role", wantErr: true},
		{name: "empty is rejected", input: " ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeLabel(tt.input)
			if tt.wantErr {
				if !errors.Is(err, domain.ErrInvalidRoleLabel) {
					t.Fatalf("normalizeLabel() error = %v, want %v", err, domain.ErrInvalidRoleLabel)
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeLabel() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("normalizeLabel() = %q, want %q", got, tt.want)
			}
		})
	}
}
