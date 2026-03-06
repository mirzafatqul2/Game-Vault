package wallet

import "time"

type WalletTransaction struct {
	ID          string
	UserID      string
	Amount      int
	Type        string
	Description string
	CreatedAt   time.Time
}
