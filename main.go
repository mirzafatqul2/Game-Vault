package main

import (
	gameCtrl "Mini-Project-Game-Vault-API/app/controller/game"
	rentalCtrl "Mini-Project-Game-Vault-API/app/controller/rental"
	reviewCtrl "Mini-Project-Game-Vault-API/app/controller/review"
	userCtrl "Mini-Project-Game-Vault-API/app/controller/user"
	whislistCtrl "Mini-Project-Game-Vault-API/app/controller/whislist"
	"Mini-Project-Game-Vault-API/app/router"
	"Mini-Project-Game-Vault-API/database"
	_ "Mini-Project-Game-Vault-API/docs"
	gameRepo "Mini-Project-Game-Vault-API/repository/game"
	mailjet "Mini-Project-Game-Vault-API/repository/mailjet"
	rentalRepo "Mini-Project-Game-Vault-API/repository/rental"
	reviewRepo "Mini-Project-Game-Vault-API/repository/review"
	userRepo "Mini-Project-Game-Vault-API/repository/user"
	walletRepo "Mini-Project-Game-Vault-API/repository/wallet"
	whislistRepo "Mini-Project-Game-Vault-API/repository/whislist"
	gameSvc "Mini-Project-Game-Vault-API/service/game"
	rentalSvc "Mini-Project-Game-Vault-API/service/rental"
	reviewSvc "Mini-Project-Game-Vault-API/service/review"
	userSvc "Mini-Project-Game-Vault-API/service/user"
	whislistSvc "Mini-Project-Game-Vault-API/service/whistlist"
	"Mini-Project-Game-Vault-API/util/response"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	cfg "github.com/pobyzaarif/go-config"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var loggerOption = slog.HandlerOptions{AddSource: true}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))

type Config struct {
	AppPort                 string `env:"APP_PORT"`
	AppHost                 string `env:"APP_HOST"`
	AppDeploymentUrl        string `env:"APP_DEPLOYMENT_URL"`
	AppEmailVerificationKey string `env:"APP_EMAIL_VERIFICATION_KEY"`
	AppJWTSecret            string `env:"APP_JWT_SECRET"`

	MailjetBaseUrl           string `env:"MAILJET_BASE_URL"`
	MailjetBasicAuthUsername string `env:"MAILJET_BASIC_AUTH_USERNAME"`
	MailjetBasicAuthPassword string `env:"MAILJET_BASIC_AUTH_PASSWORD"`
	MailjetSenderEmail       string `env:"MAILJET_SENDER_EMAIL"`
	MailjetSenderName        string `env:"MAILJET_SENDER_NAME"`

	RawgBaseUrl string `env:"RAWG_BASE_URL"`
	RawgAppKey  string `env:"RAWG_APP_KEY"`
}

// @title Library Management API
// @version 1.0
// @description Backend API for Library System
// @host 127.0.0.1:8000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load config
	config := Config{}
	err := cfg.LoadConfig(&config)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.Info("Config loaded")

	// Database connection
	db := database.InitDatabase()
	logger.Info("Database connected")

	walletRepository := walletRepo.NewWalletRepository(db)
	userRepository := userRepo.NewUserRepository(db)
	mailjetRepository := mailjet.NewMailjetRepository(
		logger,
		mailjet.MailjetConfig{
			MailjetBaseURL:           config.MailjetBaseUrl,
			MailjetBasicAuthUsername: config.MailjetBasicAuthUsername,
			MailjetBasicAuthPassword: config.MailjetBasicAuthPassword,
			MailjetSenderEmail:       config.MailjetSenderEmail,
			MailjetSenderName:        config.MailjetSenderName,
		},
	)
	userService := userSvc.NewUserService(
		logger,
		userRepository,
		config.AppDeploymentUrl,
		config.AppJWTSecret,
		config.AppEmailVerificationKey,
		mailjetRepository,
		walletRepository,
	)

	userController := userCtrl.NewUserController(
		logger,
		userService,
	)

	rawgRepo := gameRepo.NewGameRepository(
		logger,
		gameRepo.RAWGConfig{
			BaseURL: config.RawgBaseUrl,
			APIKey:  config.RawgAppKey,
		},
		db,
	)
	gameService := gameSvc.NewGameService(rawgRepo)
	gameController := gameCtrl.NewGameController(gameService)

	rentalRepository := rentalRepo.NewRentalRepository(db)
	rentalService := rentalSvc.NewRentalService(
		rentalRepository,
		rawgRepo,
		userRepository,
	)
	rentalController := rentalCtrl.NewRentalController(rentalService)

	whislistRepository := whislistRepo.NewWhislistRepository(db)
	whisllistService := whislistSvc.NewWhislistService(whislistRepository)
	whislistController := whislistCtrl.NewWhislistController(whisllistService)

	reviewRepository := reviewRepo.NewReviewRepository(db)
	reviewService := reviewSvc.NewReviewService(reviewRepository)
	reviewController := reviewCtrl.NewReviewController(reviewService)

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Skipper: middleware.DefaultSkipper,
			Format: `{"time":"${time_rfc3339_nano}","level":"INFO","id":"${id}","remote_ip":"${remote_ip}",` +
				`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
				`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
				`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
			CustomTimeFormat: "2006-01-02 15:04:05.00000",
		},
	))
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.NewSuccessResponse("Game Vault API running", nil))
	})

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())

	router.RegisterPath(
		e,
		config.AppJWTSecret,
		userController,
		gameController,
		rentalController,
		whislistController,
		reviewController,
	)

	// Start server
	address := config.AppHost + ":" + config.AppPort
	go func() {
		if err := e.Start(address); err != http.ErrServerClosed {
			log.Fatal("Failed on http server " + config.AppPort)
		}
	}()

	logger.Info("Api service running in " + address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutting down echo server", "err", err)
	} else {
		logger.Info("Successfully shutting down echo server")
	}
}
