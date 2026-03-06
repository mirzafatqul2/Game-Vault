package user

import "time"

type User struct {
	ID            string
	FullName      string
	Username      string
	Email         string
	Password      string
	DepositAmount int
	Role          string
	IsVerified    bool
	CreatedAt     time.Time
}
