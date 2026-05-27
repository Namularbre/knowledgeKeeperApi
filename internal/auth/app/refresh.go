package app

import (
	"context"
	"errors"
	"time"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

type RefreshSession struct {
	Users         domain.UserRepository
	RefreshTokens domain.RefreshTokenRepository
	Tokens        domain.TokenIssuer
	RefreshTTL    time.Duration
}

func (uc RefreshSession) Execute(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	if refreshToken == "" {
		return domain.TokenPair{}, domain.ErrInvalidRefresh
	}

	hash := uc.Tokens.HashRefreshToken(refreshToken)

	existing, err := uc.RefreshTokens.FindByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidRefresh) {
			return domain.TokenPair{}, domain.ErrInvalidRefresh
		}
		return domain.TokenPair{}, err
	}

	now := Now().UTC()
	if !existing.IsActive(now) {
		return domain.TokenPair{}, domain.ErrInvalidRefresh
	}

	// Rotate: revoke the presented token before minting a new pair, so a
	// replay of the same token after rotation is rejected.
	if err := uc.RefreshTokens.Revoke(ctx, hash); err != nil {
		return domain.TokenPair{}, err
	}

	return IssueTokenPair(ctx, uc.Tokens, uc.RefreshTokens, existing.UserID, uc.RefreshTTL)
}
