package router

import (
	"Mini-Project-Game-Vault-API/app/controller/game"
	"Mini-Project-Game-Vault-API/app/controller/rental"
	"Mini-Project-Game-Vault-API/app/controller/review"
	"Mini-Project-Game-Vault-API/app/controller/user"
	"Mini-Project-Game-Vault-API/app/controller/whislist"
	"Mini-Project-Game-Vault-API/app/middleware"
	"Mini-Project-Game-Vault-API/util/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterPath(
	e *echo.Echo,
	jwtSecret string,
	ctrlUser *user.UserController,
	CtrlGame *game.GameController,
	CtrlRental *rental.RentalController,
	CtrlWhislist *whislist.WhislistController,
	CtrlReview *review.ReviewController,
) {
	// Setup routes
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.NewSuccessResponse(
			"pong", nil,
		))
	})

	// Init JWT
	jwtMiddleware := middleware.JWTMiddleware(jwtSecret)

	// Init ACL
	userNAdminAccess := middleware.ACLMiddleware(map[string]bool{
		"admin": true,
		"user":  true,
	})
	adminAccess := middleware.ACLMiddleware(map[string]bool{
		"admin": true,
	})
	// superadminAccess := middleware.ACLMiddleware(map[string]bool{
	// 	"superadmin": true,
	// })

	// Public endpoint
	publicEndpoint := e.Group("/users")
	publicEndpoint.POST("/register", ctrlUser.Register)
	publicEndpoint.POST("/login", ctrlUser.Login)
	publicEndpoint.GET("/email-verification/:code", ctrlUser.VerifyEmail)

	// User endpoint
	userEndpoint := e.Group("/users", jwtMiddleware)
	userEndpoint.GET("/profile", ctrlUser.GetUserProfile, userNAdminAccess)
	userEndpoint.POST("/deposit", ctrlUser.DepositAmount, userNAdminAccess)

	// Game endpoint
	gameEndpoint := e.Group("/games", jwtMiddleware)
	gameEndpoint.GET("/explore", CtrlGame.GetGames)
	gameEndpoint.GET("/ready", CtrlGame.GetAllGamesReady)
	gameEndpoint.GET("/:id", CtrlGame.GetGameByID)
	gameEndpoint.POST("/import", CtrlGame.ImportGames, adminAccess)

	// Rental endpoint
	rentalEndpoint := e.Group("/rentals", jwtMiddleware)
	rentalEndpoint.POST("/checkout", CtrlRental.RentalGames, userNAdminAccess)
	rentalEndpoint.GET("", CtrlRental.GetHistoryRental, userNAdminAccess)

	// Whislist endpoint
	whislistEndpoint := e.Group("/whislists", jwtMiddleware)
	whislistEndpoint.POST("", CtrlWhislist.AddWhislist, userNAdminAccess)
	whislistEndpoint.GET("", CtrlWhislist.GetAllWhislist, userNAdminAccess)

	// Review endpoint
	reviewEndpoint := e.Group("/reviews", jwtMiddleware)
	reviewEndpoint.POST("", CtrlReview.CreateReview, userNAdminAccess)
}
