package whistlist

import "time"

type Whislist struct {
	ID        string
	UserID    string
	GameID    string
	CreatedAt time.Time
}
