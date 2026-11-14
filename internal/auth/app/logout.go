package app

import (
	"context"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

type LogoutSession struct {
	RefreshTokens domain.RefreshTokenRepository
	Tokens        domain.TokenIssuer
}

// Execute revokes the given refresh token. It is idempotent: an unknown or
// already revoked token is not an error, so a client logout always succeeds.
func (uc LogoutSession) Execute(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return domain.ErrInvalidRefresh
	}
	return uc.RefreshTokens.Revoke(ctx, uc.Tokens.HashRefreshToken(refreshToken))
}
