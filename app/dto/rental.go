package dto

import "time"

type RentalRequest struct {
	GameID     string `json:"game_id"`
	RentalDays int    `json:"rental_days" validate:"required,gt=0"`
}

type RentalResponse struct {
	GameID     string    `json:"game_id"`
	UserID     string    `json:"user_id"`
	RentalDays int       `json:"rental_days"`
	TotalCost  int       `json:"total_cost"`
	Status     string    `json:"status"`
	RentedAt   time.Time `json:"rented_at"`
	DueDate    time.Time `json:"due_date"`
}
