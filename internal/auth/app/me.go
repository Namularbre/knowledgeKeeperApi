package app

import (
	"context"
	"errors"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

type Me struct {
	Users domain.UserRepository
}

func (uc *Me) Execute(ctx context.Context, userId int64) (domain.User, error) {
	if userId < 0 {
		return domain.User{}, errors.New("userId is negative")
	}
	return uc.Users.FindByID(ctx, userId)
}
