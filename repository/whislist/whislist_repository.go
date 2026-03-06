package whislist

import (
	"Mini-Project-Game-Vault-API/service/whistlist"
	"context"

	"gorm.io/gorm"
)

type WhislistRepository struct {
	DB *gorm.DB
}

func NewWhislistRepository(db *gorm.DB) *WhislistRepository {
	return &WhislistRepository{
		DB: db,
	}
}

func (r *WhislistRepository) Create(ctx context.Context, whislist whistlist.Whislist) error {
	return r.DB.WithContext(ctx).Create(&whislist).Error
}

func (r *WhislistRepository) GetAll(ctx context.Context, userID string) ([]whistlist.Whislist, error) {
	var whis []whistlist.Whislist

	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&whis).Error
	if err != nil {
		return whis, err
	}

	return whis, nil
}
