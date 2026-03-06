package dto

type ReviewRequest struct {
	GameID  string `json:"game_id"`
	Rating  int    `json:"rating" validate:"required"`
	Comment string `json:"comment" validate:"required"`
}

type ReviewResponse struct {
	GameID  string `json:"game_id"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}
