package game

import "time"

type Game struct {
	ID                string
	ExternalID        string
	Name              string
	Category          string
	StockAvailability int
	RentalCost        int
	ImageUrl          string
	CreatedAt         time.Time
}
