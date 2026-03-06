package wallet

import "context"

type WalletRepository interface {
	Create(ctx context.Context, wallet WalletTransaction) error
}
