package user

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	UpdatEmailVerification(ctx context.Context, user User) error
	UpdateDeposit(ctx context.Context, userID string, newBalance int) error
}
