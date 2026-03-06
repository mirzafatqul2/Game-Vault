package review

import (
	"context"
)

type ReviewRepository interface {
	Create(ctx context.Context, review Review) error
}
