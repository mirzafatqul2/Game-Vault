package review

import "time"

type Review struct {
	ID        string
	UserID    string
	GameID    string
	Rating    int
	Comment   string
	CreatedAt time.Time
}
