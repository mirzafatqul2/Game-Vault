package rental

import (
	"Mini-Project-Game-Vault-API/service/game"
	"Mini-Project-Game-Vault-API/service/rental"
	"Mini-Project-Game-Vault-API/service/user"
	"context"

	"gorm.io/gorm"
)

type RentalRepository struct {
	DB *gorm.DB
}

func NewRentalRepository(db *gorm.DB) *RentalRepository {
	return &RentalRepository{DB: db}
}

func (r *RentalRepository) CheckoutTransaction(ctx context.Context, user user.User, game game.Game, rental rental.Rental) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		if err := tx.Save(&game).Error; err != nil {
			return err
		}

		if err := tx.Create(&rental).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *RentalRepository) GetAll(ctx context.Context, userID string) ([]rental.Rental, error) {
	var rental []rental.Rental

	err := r.DB.WithContext(ctx).Where("user_id =?", userID).Find(&rental).Error
	if err != nil {
		return nil, err
	}

	return rental, nil
}
