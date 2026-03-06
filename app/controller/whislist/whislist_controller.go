package whislist

import (
	"Mini-Project-Game-Vault-API/app/dto"
	"Mini-Project-Game-Vault-API/service/whistlist"
	"Mini-Project-Game-Vault-API/util/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type WhislistController struct {
	WhislistService whistlist.WhislistService
}

func NewWhislistController(WhislistService whistlist.WhislistService) *WhislistController {
	return &WhislistController{
		WhislistService: WhislistService,
	}
}

// @Summary Add Wishlist
// @Description Add game to wishlist
// @Tags Wishlist
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.WishlistRequest true "Wishlist Request"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /wishlists [post]
func (c *WhislistController) AddWhislist(ctx echo.Context) error {
	// Get user ID from JWT middleware context
	userID, ok := ctx.Get("id").(string)
	if !ok || userID == "" {
		return ctx.JSON(http.StatusUnauthorized,
			response.NewErrorResponse("unauthorized"))
	}

	var req dto.WishlistRequest

	// Bind JSON request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse("invalid request body"))
	}

	whislist, err := c.WhislistService.AddWhislist(
		ctx.Request().Context(),
		userID,
		req.GameID,
	)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse(err.Error()))
	}

	resp := dto.WishlistResponse{
		GameID: whislist.GameID,
	}

	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("success", resp))
}

// GetAllWhislist godoc
// @Summary Get all user whislist
// @Description Retrieve list of user whislist
// @Tags Whislists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /rentals [get]
func (c *WhislistController) GetAllWhislist(ctx echo.Context) error {
	// Get user ID from JWT middleware context
	userID, ok := ctx.Get("id").(string)
	if !ok || userID == "" {
		return ctx.JSON(http.StatusUnauthorized,
			response.NewErrorResponse("unauthorized"))
	}

	rentals, err := c.WhislistService.GetAllWhislist(
		ctx.Request().Context(),
		userID,
	)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse("failed to fetch data"))
	}

	var results []dto.WishlistResponse
	for _, r := range rentals {
		results = append(results, dto.WishlistResponse{
			GameID: r.GameID,
		})
	}

	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("success", results))
}
