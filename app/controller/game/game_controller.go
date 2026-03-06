package game

import (
	"Mini-Project-Game-Vault-API/app/dto"
	"Mini-Project-Game-Vault-API/service/game"
	"Mini-Project-Game-Vault-API/util/response"
	"Mini-Project-Game-Vault-API/util/validator"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type GameController struct {
	GameService game.GameService
}

func NewGameController(GameService game.GameService) *GameController {
	return &GameController{
		GameService: GameService,
	}
}

// GetGames godoc
// @Summary Explore games
// @Description Get list games from RAWG API
// @Tags Games
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /games/explore [get]
func (c *GameController) GetGames(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	games, err := c.GameService.GetGames(ctx.Request().Context(), page, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse("failed to fetch games"))
	}

	var results []dto.GameResponse
	for _, g := range games {
		results = append(results, dto.GameResponse{
			ExternalID: g.ExternalID,
			Name:       g.Name,
			Category:   g.Category,
			ImageURL:   g.ImageUrl,
		})
	}

	listGames := dto.ExploreGamesResponse{
		Page:  page,
		Limit: limit,
		Data:  results,
	}

	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("success", listGames))
}

// ImportGame godoc
// @Summary Import game from RAWG
// @Tags Games
// @Accept json
// @Produce json
// @Param request body dto.ImportGameRequest true "Import Game"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /games/import [post]
func (c *GameController) ImportGames(ctx echo.Context) error {
	var req dto.ImportGameRequest

	// Bind JSON request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse("invalid request body"))
	}

	// Validate request
	if err := validator.ValidateImportGameRequest(req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse(err.Error()))
	}

	// Call service layer
	game, err := c.GameService.ImportGames(
		ctx.Request().Context(),
		req.ExternalID,
		req.StockAvailability,
		req.RentalCost,
	)

	if err != nil {
		if strings.Contains(err.Error(), "registered") {
			return ctx.JSON(http.StatusBadRequest,
				response.NewErrorResponse(err.Error()))
		}
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse(err.Error()))

	}

	resp := dto.GameResponse{
		ExternalID:        game.ExternalID,
		Name:              game.Name,
		Category:          game.Category,
		StockAvailability: game.StockAvailability,
		RentalCost:        game.RentalCost,
		ImageURL:          game.ImageUrl,
	}

	// response JSON
	return ctx.JSON(http.StatusCreated,
		response.NewSuccessResponse("success", resp),
	)
}

// GetGameByID godoc
// @Summary Get detail game from database
// @Tags Games
// @Accept json
// @Produce json
// @Param id path int true "Game ID"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /games/{id} [get]
func (c *GameController) GetGameByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse("invalid game id"))
	}

	game, err := c.GameService.GetGameByID(ctx.Request().Context(), id)

	if err != nil {
		return ctx.JSON(http.StatusNotFound,
			response.NewErrorResponse("game not found"))
	}

	resp := dto.GameResponse{
		ExternalID:        game.ExternalID,
		Name:              game.Name,
		Category:          game.Category,
		StockAvailability: game.StockAvailability,
		RentalCost:        game.RentalCost,
		ImageURL:          game.ImageUrl,
	}

	// response JSON
	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("success", resp),
	)
}

// GetAllGamesReady godoc
// @Summary Get all games ready for rental
// @Description Retrieve list of games that have stock availability and rental cost information
// @Tags Games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]dto.GameResponse}
// @Failure 500 {object} response.ErrorResponse
// @Router /games/ready [get]
func (c *GameController) GetAllGamesReady(ctx echo.Context) error {
	games, err := c.GameService.GetAllGamesReady(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse("failed to fetch games"))
	}

	var result []dto.GameResponse
	for _, g := range games {
		result = append(result, dto.GameResponse{
			ExternalID:        g.ExternalID,
			Name:              g.Name,
			Category:          g.Category,
			StockAvailability: g.StockAvailability,
			RentalCost:        g.RentalCost,
			ImageURL:          g.ImageUrl,
		})
	}

	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("success", result))
}
