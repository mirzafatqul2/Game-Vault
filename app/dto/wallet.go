package dto

type DepositRequest struct {
	Amount      int    `json:"amount" validate:"required,gt=0"`
	Description string `json:"description" validate:"required"`
}

type WalletResponse struct {
	UserID      string `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}
