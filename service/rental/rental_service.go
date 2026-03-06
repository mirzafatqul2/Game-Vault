package rental

import (
	"Mini-Project-Game-Vault-API/service/game"
	"Mini-Project-Game-Vault-API/service/user"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type rentalService struct {
	rentalRepo RentalRepository
	gameRepo   game.GameRepository
	userRepo   user.UserRepository
}

type RentalService interface {
	RentalGames(ctx context.Context, userID string, gameID string, rentalDays int) (Rental, error)
	GetHistoryRental(ctx context.Context, userID string) ([]Rental, error)
}

func NewRentalService(rentalRepo RentalRepository, gameRepo game.GameRepository, userRepo user.UserRepository) RentalService {
	return &rentalService{rentalRepo: rentalRepo, userRepo: userRepo, gameRepo: gameRepo}
}

func (s *rentalService) RentalGames(ctx context.Context, userID string, gameID string, rentalDays int) (Rental, error) {
	getGame, err := s.gameRepo.GetGameByID(ctx, gameID)
	if err != nil {
		return Rental{}, err
	}

	getUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return Rental{}, err
	}

	if getGame.StockAvailability <= 0 {
		return Rental{}, errors.New("game out of stock")
	}

	totalCost := getGame.RentalCost * rentalDays

	if getUser.DepositAmount < totalCost {
		return Rental{}, errors.New("insufficient balance")
	}

	getUser.DepositAmount -= totalCost
	getGame.StockAvailability--

	rental := Rental{
		ID:         uuid.New().String(),
		UserID:     userID,
		GameID:     gameID,
		RentalDays: rentalDays,
		TotalCost:  totalCost,
		Status:     "ongoing",
		RentedAt:   time.Now(),
		DueDate:    time.Now().Add(time.Duration(rentalDays*24) * time.Hour),
	}

	s.rentalRepo.CheckoutTransaction(ctx, getUser, getGame, rental)

	return rental, nil
}

func (s *rentalService) GetHistoryRental(ctx context.Context, userID string) ([]Rental, error) {
	return s.rentalRepo.GetAll(ctx, userID)
}
