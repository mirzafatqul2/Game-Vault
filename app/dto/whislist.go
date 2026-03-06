package dto

type WishlistRequest struct {
	GameID string `json:"game_id"`
}

type WishlistResponse struct {
	GameID string `json:"game_id"`
}
