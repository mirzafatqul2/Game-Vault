package wallet

import (
	"Mini-Project-Game-Vault-API/service/wallet"
	"context"

	"gorm.io/gorm"
)

type WalletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		DB: db,
	}
}

func (r *WalletRepository) Create(ctx context.Context, wallet wallet.WalletTransaction) error {
	return r.DB.WithContext(ctx).Create(&wallet).Error
}
