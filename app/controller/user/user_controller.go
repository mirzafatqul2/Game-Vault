package user

import (
	"Mini-Project-Game-Vault-API/app/dto"
	"Mini-Project-Game-Vault-API/service/user"
	"Mini-Project-Game-Vault-API/util/response"
	"Mini-Project-Game-Vault-API/util/validator"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	logger *slog.Logger

	userService user.UserService
}

func NewUserController(logger *slog.Logger, userService user.UserService) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with email, password, and fullname
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "User registration request"
// @Success      201 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /users/register [post]
// Register handles POST /users/register
func (c *UserController) Register(ctx echo.Context) error {
	var req dto.RegisterRequest

	// Bind JSON request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse("invalid request body"))
	}

	// Validate request
	if err := validator.ValidateRegisterRequest(req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse(err.Error()))
	}

	// Mapping dto to user entity
	userEntity := user.User{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Call service layer
	id, err := c.userService.Register(ctx.Request().Context(), userEntity)

	if err != nil {
		if strings.Contains(err.Error(), "registered") {
			return ctx.JSON(http.StatusBadRequest,
				response.NewErrorResponse(err.Error()))
		}
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse(err.Error()))

	}

	return ctx.JSON(http.StatusCreated,
		response.NewSuccessResponse(
			"user registered successfully",
			id,
		),
	)
}

func (c UserController) VerifyEmail(ctx echo.Context) error {
	encCode := ctx.Param("code")

	err := c.userService.VerifyEmail(ctx.Request().Context(), encCode)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			return ctx.JSON(http.StatusUnauthorized,
				response.NewErrorResponse(err.Error()))
		}
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse(err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.NewSuccessResponse("email success to verify", nil))
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User Login Request"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router /users/login [post]
func (c *UserController) Login(ctx echo.Context) error {
	var req dto.LoginRequest

	// Bind JSON request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse("invalid request body"))
	}

	// Validate request
	if err := validator.ValidateLoginRequest(req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse(err.Error()))
	}

	// Call service layer for authentication
	token, err := c.userService.Login(
		ctx.Request().Context(),
		req.Email,
		req.Password,
	)

	// Authentication Failed
	if err != nil {
		if strings.Contains(err.Error(), "email address") {
			return ctx.JSON(http.StatusUnauthorized,
				response.NewErrorResponse(err.Error()))
		}
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse(err.Error()))
	}

	// Return success response with JWT token
	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("login success", response.TokenResponse{
			Token: token,
		}),
	)
}

// GetUserProfile godoc
// @Summary		Get authenticated user profile
// @Description	Get user profile information based on authenticated JWT/session user id
// @Tags		Users
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Success		200	{object}	response.SuccessResponse
// @Failure		401	{object}	response.ErrorResponse
// @Failure		500	{object}	response.ErrorResponse
// @Router		/users/me [get]
// GetUserProfile returns authenticated user profile information.
func (c *UserController) GetUserProfile(ctx echo.Context) error {
	// Get user ID from JWT middleware context
	userID, ok := ctx.Get("id").(string)
	if !ok || userID == "" {
		return ctx.JSON(http.StatusUnauthorized,
			response.NewErrorResponse("unauthorized"))
	}

	// Call service layer to fetch user profile data
	data, err := c.userService.GetUserProfile(
		ctx.Request().Context(),
		userID,
	)

	// Mapping entity to DTO
	newData := dto.ProfileResponse{
		FullName:      data.FullName,
		Username:      data.Username,
		Email:         data.Email,
		DepositAmount: data.DepositAmount,
	}

	// Handle internal server error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse("failed to fetch user profile"))
	}

	// Return success response with user profile data
	return ctx.JSON(http.StatusOK,
		response.NewSuccessResponse("success", newData))
}

// @Summary User Deposit Balance
// @Description Add balance to user deposit_amount
// @Tags Wallet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.DepositRequest true "Deposit Request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/deposit [post]
func (c *UserController) DepositAmount(ctx echo.Context) error {
	// Get user ID from JWT middleware context
	userID, ok := ctx.Get("id").(string)
	if !ok || userID == "" {
		return ctx.JSON(http.StatusUnauthorized,
			response.NewErrorResponse("unauthorized"))
	}

	var req dto.DepositRequest

	// Bind JSON request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse("invalid request body"))
	}

	// Validate request
	if err := validator.ValidateDepositRequest(req); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.NewErrorResponse(err.Error()))
	}

	deposit, err := c.userService.DepositAmount(
		ctx.Request().Context(),
		userID,
		req.Description,
		req.Amount,
	)

	if err != nil {
		if strings.Contains(err.Error(), "registered") {
			return ctx.JSON(http.StatusBadRequest,
				response.NewErrorResponse(err.Error()))
		}
		return ctx.JSON(http.StatusInternalServerError,
			response.NewErrorResponse(err.Error()))

	}

	resp := dto.WalletResponse{
		UserID:      deposit.ID,
		Amount:      req.Amount,
		Description: req.Description,
	}

	// response JSON
	return ctx.JSON(http.StatusCreated,
		response.NewSuccessResponse("success", resp),
	)
}
