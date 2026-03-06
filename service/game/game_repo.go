package game

import (
	"context"
)

type GameRepository interface {
	GetGames(ctx context.Context, page int, pageSize int) ([]Game, error)
	GetGameDetail(ctx context.Context, externalID string) (Game, error)
	Create(ctx context.Context, game Game) error
	GetGameByID(ctx context.Context, id string) (Game, error)
	GetGameByExternalID(ctx context.Context, externalID string) (Game, error)
	GetGameReady(ctx context.Context) ([]Game, error)
}
