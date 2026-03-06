package dto

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ProfileResponse struct {
	FullName      string `json:"full_name"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	DepositAmount int    `json:"deposit_amount"`
}
