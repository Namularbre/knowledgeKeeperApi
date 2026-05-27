package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

type LoginUser struct {
	Users         domain.UserRepository
	RefreshTokens domain.RefreshTokenRepository
	Hasher        domain.PasswordHasher
	Tokens        domain.TokenIssuer
	RefreshTTL    time.Duration
}

type LoginInput struct {
	Email    string
	Password string
}

func (uc LoginUser) Execute(ctx context.Context, in LoginInput) (domain.TokenPair, error) {
	email, err := normalizeEmail(in.Email)
	if err != nil {
		return domain.TokenPair{}, domain.ErrInvalidCredentials
	}

	user, err := uc.Users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.TokenPair{}, domain.ErrInvalidCredentials
		}
		return domain.TokenPair{}, err
	}

	if err := uc.Hasher.Compare(user.PasswordHash, in.Password); err != nil {
		return domain.TokenPair{}, domain.ErrInvalidCredentials
	}

	return IssueTokenPair(ctx, uc.Tokens, uc.RefreshTokens, user.ID, uc.RefreshTTL)
}

// IssueTokenPair mints an access token, generates a fresh opaque refresh
// token, persists its hash, and returns both to the caller. Shared by login
// and refresh use cases.
func IssueTokenPair(
	ctx context.Context,
	tokens domain.TokenIssuer,
	refreshes domain.RefreshTokenRepository,
	userID int64,
	refreshTTL time.Duration,
) (domain.TokenPair, error) {
	now := Now().UTC()

	access, accessExp, err := tokens.IssueAccessToken(userID, now)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("issue access token: %w", err)
	}

	refresh, refreshHash, err := tokens.GenerateRefreshToken()
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	refreshExp := now.Add(refreshTTL)
	if _, err := refreshes.Create(ctx, userID, refreshHash, refreshExp); err != nil {
		return domain.TokenPair{}, fmt.Errorf("persist refresh token: %w", err)
	}

	return domain.TokenPair{
		AccessToken:      access,
		AccessExpiresAt:  accessExp,
		RefreshToken:     refresh,
		RefreshExpiresAt: refreshExp,
	}, nil
}
