package rental

import "time"

type Rental struct {
	ID         string
	UserID     string
	GameID     string
	RentalDays int
	TotalCost  int
	Status     string
	RentedAt   time.Time
	DueDate    time.Time
	ReturnedAt *time.Time
	PenaltyFee int
}
