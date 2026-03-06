package validator

import (
	"Mini-Project-Game-Vault-API/app/dto"

	"github.com/go-playground/validator/v10"
)

// ValidateRegisterRequest validates register request
func ValidateRegisterRequest(register dto.RegisterRequest) error {
	validate := validator.New()
	return validate.Struct(register)
}

// ValidateLoginRequest validates login request
func ValidateLoginRequest(login dto.LoginRequest) error {
	validate := validator.New()
	return validate.Struct(login)
}

func ValidateImportGameRequest(req dto.ImportGameRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

func ValidateDepositRequest(req dto.DepositRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

func ValidateRentalRequest(req dto.RentalRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}

func ValidateReviewRequest(req dto.ReviewRequest) error {
	validate := validator.New()
	return validate.Struct(req)
}
