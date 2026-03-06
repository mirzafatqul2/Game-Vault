package rental

import (
	"Mini-Project-Game-Vault-API/app/dto"
	"Mini-Project-Game-Vault-API/service/rental"
	"Mini-Project-Game-Vault-API/util/response"
	"Mini-Project-Game-Vault-API/util/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RentalController struct {
	RentalService rental.RentalService
}

func NewRentalController(RentalService rental.RentalService) *RentalController {
	return &RentalController{
		RentalService: RentalService,
	}
}

// @Summary Rental Checkout Game
// @Description User checkout rental game with deposit balance deduction
// @Tags Rentals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.RentalRequest true "Rental Checkout Request"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /rentals/checkout [post]
func (c *RentalController) RentalGames(ctx echo.Context) error {
	// Get user ID from JWT middleware context
	userID, ok := ctx.Get("id").(string)
	if !ok || userID == "" {
		return ctx.JSON(http.StatusUnauthorized,
			response.NewErrorResponse("unauthorized"))
	}

	var req dto.RentalRequest

	// Bind JSON request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse("invalid request body"))
	}

	// Validate request
	if err := validator.ValidateRentalRequest(req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse(err.Error()))
	}

	rental, err := c.RentalService.RentalGames(
		ctx.Request().Context(),
		userID,
		req.GameID,
		req.RentalDays,
	)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse(err.Error()))
	}

	resp := dto.RentalResponse{
		GameID:     rental.GameID,
		UserID:     rental.UserID,
		RentalDays: rental.RentalDays,
		TotalCost:  rental.TotalCost,
		Status:     rental.Status,
		RentedAt:   rental.RentedAt,
		DueDate:    rental.DueDate,
	}

	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("success", resp))
}

// GetHistoryRental godoc
// @Summary Get all user history rental
// @Description Retrieve list of history rental
// @Tags Rentals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /rentals [get]
func (c *RentalController) GetHistoryRental(ctx echo.Context) error {
	// Get user ID from JWT middleware context
	userID, ok := ctx.Get("id").(string)
	if !ok || userID == "" {
		return ctx.JSON(http.StatusUnauthorized,
			response.NewErrorResponse("unauthorized"))
	}

	rentals, err := c.RentalService.GetHistoryRental(
		ctx.Request().Context(),
		userID,
	)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse("failed to fetch data"))
	}

	var results []dto.RentalResponse
	for _, r := range rentals {
		results = append(results, dto.RentalResponse{
			GameID:     r.GameID,
			UserID:     r.UserID,
			RentalDays: r.RentalDays,
			TotalCost:  r.TotalCost,
			Status:     r.Status,
			RentedAt:   r.RentedAt,
			DueDate:    r.DueDate,
		})
	}

	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("success", results))
}
