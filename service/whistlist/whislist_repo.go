package whistlist

import (
	"context"
)

type WhislistRepository interface {
	Create(ctx context.Context, whislist Whislist) error
	GetAll(ctx context.Context, userID string) ([]Whislist, error)
}
