package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	rolesdomain "github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

func TestRequireAnyRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		userID     *int64
		roles      []rolesdomain.Role
		err        error
		wantStatus int
		wantCalled bool
	}{
		{name: "missing authenticated user", wantStatus: http.StatusUnauthorized},
		{name: "role lookup failure", userID: int64Ptr(1), err: errors.New("database unavailable"), wantStatus: http.StatusInternalServerError},
		{name: "role is not allowed", userID: int64Ptr(1), roles: []rolesdomain.Role{{Label: "member"}}, wantStatus: http.StatusForbidden},
		{name: "allowed role", userID: int64Ptr(1), roles: []rolesdomain.Role{{Label: "admin"}}, wantStatus: http.StatusNoContent, wantCalled: true},
		{name: "role comparison is normalized", userID: int64Ptr(1), roles: []rolesdomain.Role{{Label: " Admin "}}, wantStatus: http.StatusNoContent, wantCalled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusNoContent)
			})
			handler := RequireAnyRole(fakeRoleRepository{roles: tt.roles, err: tt.err}, "admin")(next)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.userID != nil {
				req = req.WithContext(context.WithValue(req.Context(), userIDKey, *tt.userID))
			}
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if res.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", res.Code, tt.wantStatus)
			}
			if called != tt.wantCalled {
				t.Fatalf("next called = %t, want %t", called, tt.wantCalled)
			}
		})
	}
}

func int64Ptr(value int64) *int64 { return &value }

type fakeRoleRepository struct {
	roles []rolesdomain.Role
	err   error
}

func (f fakeRoleRepository) Create(context.Context, string) (rolesdomain.Role, error) {
	return rolesdomain.Role{}, nil
}

func (f fakeRoleRepository) FetchByPage(context.Context, uint64, uint64) ([]rolesdomain.Role, error) {
	return nil, nil
}

func (f fakeRoleRepository) FindByID(context.Context, uint64) (rolesdomain.Role, error) {
	return rolesdomain.Role{}, nil
}

func (f fakeRoleRepository) SearchByLabel(context.Context, string) ([]rolesdomain.Role, error) {
	return nil, nil
}

func (f fakeRoleRepository) FindByUserID(context.Context, uint64) ([]rolesdomain.Role, error) {
	return f.roles, f.err
}

func (f fakeRoleRepository) AddUserRole(context.Context, uint64, uint64) ([]rolesdomain.Role, error) {
	return nil, nil
}

func (f fakeRoleRepository) RemoveUserRole(context.Context, uint64, uint64) ([]rolesdomain.Role, error) {
	return nil, nil
}
