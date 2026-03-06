package dto

type GameResponse struct {
	ExternalID        string `json:"external_id"`
	Name              string `json:"name"`
	Category          string `json:"category"`
	StockAvailability int    `json:"stock_availability"`
	RentalCost        int    `json:"rental_cost"`
	ImageURL          string `json:"image_url"`
}

type ExploreGamesResponse struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Data  []GameResponse `json:"data"`
}

type ImportGameRequest struct {
	ExternalID        string `json:"external_id" validate:"required"`
	StockAvailability int    `json:"stock_availability" validate:"required,gt=0"`
	RentalCost        int    `json:"rental_cost" validate:"required,gt=0"`
}
