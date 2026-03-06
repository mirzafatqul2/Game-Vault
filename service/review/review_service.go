package review

import (
	"context"

	"github.com/google/uuid"
)

type reviewService struct {
	reviewRepo ReviewRepository
}

type ReviewService interface {
	CreateReview(ctx context.Context, userID string, review Review) (Review, error)
}

func NewReviewService(reviewRepo ReviewRepository) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, userID string, review Review) (Review, error) {
	review = Review{
		ID:      uuid.New().String(),
		UserID:  userID,
		GameID:  review.GameID,
		Rating:  review.Rating,
		Comment: review.Comment,
	}

	err := s.reviewRepo.Create(ctx, review)
	if err != nil {
		return Review{}, err
	}

	return review, nil
}
