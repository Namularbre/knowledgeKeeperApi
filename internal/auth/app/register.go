// Package app contains the authentication use cases. Use cases orchestrate
// the domain ports; they hold no infrastructure code themselves.
package app

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

const minPasswordLength = 8

type RegisterUser struct {
	Users  domain.UserRepository
	Hasher domain.PasswordHasher
}

type RegisterInput struct {
	Email    string
	Password string
}

func (uc RegisterUser) Execute(ctx context.Context, in RegisterInput) (domain.User, error) {
	email, err := normalizeEmail(in.Email)
	if err != nil {
		return domain.User{}, err
	}
	if len(in.Password) < minPasswordLength {
		return domain.User{}, domain.ErrWeakPassword
	}

	hash, err := uc.Hasher.Hash(in.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash password: %w", err)
	}

	return uc.Users.Create(ctx, email, hash)
}

func normalizeEmail(raw string) (string, error) {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if trimmed == "" {
		return "", domain.ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(trimmed); err != nil {
		return "", domain.ErrInvalidEmail
	}
	return trimmed, nil
}

// Now is the clock used by use cases. Overridable in tests.
var Now = time.Now
