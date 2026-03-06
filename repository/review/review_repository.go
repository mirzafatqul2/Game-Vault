package review

import (
	"Mini-Project-Game-Vault-API/service/review"
	"context"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	DB *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{
		DB: db,
	}
}

func (r *ReviewRepository) Create(ctx context.Context, review review.Review) error {
	return r.DB.WithContext(ctx).Create(&review).Error
}
