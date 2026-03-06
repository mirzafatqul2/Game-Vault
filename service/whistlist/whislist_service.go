package whistlist

import (
	"context"

	"github.com/google/uuid"
)

type whislistService struct {
	whislitslRepo WhislistRepository
}

type WhislistService interface {
	AddWhislist(ctx context.Context, userID, gameID string) (Whislist, error)
	GetAllWhislist(ctx context.Context, userID string) ([]Whislist, error)
}

func NewWhislistService(whislistRepo WhislistRepository) WhislistService {
	return &whislistService{whislitslRepo: whislistRepo}
}

func (s *whislistService) AddWhislist(ctx context.Context, userID, gameID string) (Whislist, error) {
	whislist := Whislist{
		ID:     uuid.New().String(),
		UserID: userID,
		GameID: gameID,
	}

	err := s.whislitslRepo.Create(ctx, whislist)
	if err != nil {
		return Whislist{}, err
	}

	return whislist, nil
}

func (s *whislistService) GetAllWhislist(ctx context.Context, userID string) ([]Whislist, error) {
	return s.whislitslRepo.GetAll(ctx, userID)
}
