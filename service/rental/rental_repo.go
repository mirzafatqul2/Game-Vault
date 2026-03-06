package rental

import (
	"Mini-Project-Game-Vault-API/service/game"
	"Mini-Project-Game-Vault-API/service/user"
	"context"
)

type RentalRepository interface {
	CheckoutTransaction(ctx context.Context, user user.User, game game.Game, rental Rental) error
	GetAll(ctx context.Context, userID string) ([]Rental, error)
}
