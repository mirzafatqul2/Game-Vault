package review

import (
	"Mini-Project-Game-Vault-API/app/dto"
	"Mini-Project-Game-Vault-API/service/review"
	"Mini-Project-Game-Vault-API/util/response"
	"Mini-Project-Game-Vault-API/util/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReviewController struct {
	ReviewService review.ReviewService
}

func NewReviewController(ReviewService review.ReviewService) *ReviewController {
	return &ReviewController{
		ReviewService: ReviewService,
	}
}

// @Summary Create Game Review
// @Description User review rented game
// @Tags Review
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ReviewRequest true "Review Request"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /reviews [post]
func (c *ReviewController) CreateReview(ctx echo.Context) error {
	// Get user ID from JWT middleware context
	userID, ok := ctx.Get("id").(string)
	if !ok || userID == "" {
		return ctx.JSON(http.StatusUnauthorized,
			response.NewErrorResponse("unauthorized"))
	}

	var req dto.ReviewRequest

	// Bind JSON request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse("invalid request body"))
	}

	// Validate request
	if err := validator.ValidateReviewRequest(req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse(err.Error()))
	}

	request := review.Review{
		GameID:  req.GameID,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	review, err := c.ReviewService.CreateReview(
		ctx.Request().Context(),
		userID,
		request,
	)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse(err.Error()))
	}

	resp := dto.ReviewResponse{
		GameID:  review.GameID,
		Rating:  review.Rating,
		Comment: review.Comment,
	}

	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("success", resp))
}
