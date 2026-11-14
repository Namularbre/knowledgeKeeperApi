package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

// fakeRefreshRepo records revoked hashes; only Revoke is exercised here.
type fakeRefreshRepo struct {
	revoked []string
	err     error
}

func (f *fakeRefreshRepo) Create(context.Context, int64, string, time.Time) (domain.RefreshToken, error) {
	return domain.RefreshToken{}, nil
}

func (f *fakeRefreshRepo) FindByHash(context.Context, string) (domain.RefreshToken, error) {
	return domain.RefreshToken{}, domain.ErrInvalidRefresh
}

func (f *fakeRefreshRepo) Revoke(_ context.Context, hash string) error {
	if f.err != nil {
		return f.err
	}
	f.revoked = append(f.revoked, hash)
	return nil
}

// fakeIssuer hashes refresh tokens with a predictable prefix.
type fakeIssuer struct{}

func (fakeIssuer) IssueAccessToken(int64, time.Time) (string, time.Time, error) {
	return "", time.Time{}, nil
}
func (fakeIssuer) ParseAccessToken(string) (int64, error)        { return 0, nil }
func (fakeIssuer) GenerateRefreshToken() (string, string, error) { return "", "", nil }
func (fakeIssuer) HashRefreshToken(token string) string          { return "hash:" + token }

func TestLogoutSession_RevokesHashedToken(t *testing.T) {
	repo := &fakeRefreshRepo{}
	uc := LogoutSession{RefreshTokens: repo, Tokens: fakeIssuer{}}

	if err := uc.Execute(context.Background(), "abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.revoked) != 1 || repo.revoked[0] != "hash:abc" {
		t.Fatalf("expected hash:abc to be revoked, got %v", repo.revoked)
	}
}

func TestLogoutSession_EmptyToken(t *testing.T) {
	repo := &fakeRefreshRepo{}
	uc := LogoutSession{RefreshTokens: repo, Tokens: fakeIssuer{}}

	err := uc.Execute(context.Background(), "")
	if !errors.Is(err, domain.ErrInvalidRefresh) {
		t.Fatalf("expected ErrInvalidRefresh, got %v", err)
	}
	if len(repo.revoked) != 0 {
		t.Fatalf("nothing should be revoked, got %v", repo.revoked)
	}
}

func TestLogoutSession_RepoError(t *testing.T) {
	repo := &fakeRefreshRepo{err: errors.New("db down")}
	uc := LogoutSession{RefreshTokens: repo, Tokens: fakeIssuer{}}

	if err := uc.Execute(context.Background(), "abc"); err == nil {
		t.Fatal("expected repository error to be returned")
	}
}
