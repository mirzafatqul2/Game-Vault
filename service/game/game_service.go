package game

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type gameService struct {
	gameRepo GameRepository
}

type GameService interface {
	GetGames(ctx context.Context, page int, pageSize int) ([]Game, error)
	ImportGames(ctx context.Context, externalID string, stock, rentalCost int) (Game, error)
	GetGameByID(ctx context.Context, id string) (Game, error)
	GetAllGamesReady(ctx context.Context) ([]Game, error)
}

func NewGameService(gameRepo GameRepository) GameService {
	return &gameService{
		gameRepo: gameRepo,
	}
}

func (s *gameService) GetGames(ctx context.Context, page int, limit int) ([]Game, error) {
	return s.gameRepo.GetGames(ctx, page, limit)
}

func (s *gameService) ImportGames(ctx context.Context, externalID string, stock, rentalCost int) (Game, error) {
	rawgGame, err := s.gameRepo.GetGameDetail(ctx, externalID)
	if err != nil {
		return Game{}, err
	}

	game := Game{
		ID:                uuid.New().String(),
		ExternalID:        externalID,
		Name:              rawgGame.Name,
		Category:          rawgGame.Category,
		ImageUrl:          rawgGame.ImageUrl,
		StockAvailability: stock,
		RentalCost:        rentalCost,
	}

	if existing, err := s.gameRepo.GetGameByExternalID(ctx, game.ExternalID); err == nil && existing.ExternalID != "" {
		return Game{}, errors.New("game already posted")
	}

	err = s.gameRepo.Create(ctx, game)
	if err != nil {
		return Game{}, err
	}

	return game, nil
}

func (s *gameService) GetGameByID(ctx context.Context, id string) (Game, error) {
	game, err := s.gameRepo.GetGameByID(ctx, id)
	if err != nil {
		return Game{}, err
	}

	return game, nil
}

func (s *gameService) GetAllGamesReady(ctx context.Context) ([]Game, error) {
	return s.gameRepo.GetGameReady(ctx)
}
