package domain

import "testing"

func TestRoleLabelIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		label RoleLabel
		want  bool
	}{
		{label: RoleProf, want: true},
		{label: RoleAdmin, want: true},
		{label: "teacher", want: false},
		{label: "pas de role", want: false},
		{label: "", want: false},
	}

	for _, tt := range tests {
		t.Run(string(tt.label), func(t *testing.T) {
			t.Parallel()

			if got := tt.label.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %t, want %t", got, tt.want)
			}
		})
	}
}
