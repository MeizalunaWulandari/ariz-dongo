package app

import (
	"fmt"

	"github.com/MeizalunaWulandari/ariz-dongo/internal/client/fieldsa"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/config"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/handler"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/logger"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/middleware"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/repository"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/router"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func Run() {
	cfg := config.Load()

	if err := logger.Init(); err != nil {
		panic(err)
	}

	fieldsaClient := fieldsa.New(
		fieldsa.Config{
			BaseURL:         cfg.Fieldsa.BaseURL,
			DispatchBaseURL: cfg.Fieldsa.DispatchBaseURL,
			Username:        cfg.Fieldsa.Username,
			Password:        cfg.Fieldsa.Password,
			LoginBearer:     cfg.Fieldsa.LoginBearer,
		},
	)

	repo := repository.New()

	healthService := service.NewHealthService(
		cfg,
		repo,
	)

	dataService := service.NewDataService(
		fieldsaClient,
	)

	detailService := service.NewDetailService(
		fieldsaClient,
	)

	photoService := service.NewPhotoService(
		fieldsaClient,
	)

	searchService := service.NewSearchService(
		fieldsaClient,
		detailService,
		5,
	)

	healthHandler := handler.NewHealthHandler(
		healthService,
	)

	dataHandler := handler.NewDataHandler(
		dataService,
	)

	detailHandler := handler.NewDetailHandler(
		detailService,
	)

	photoHandler := handler.NewPhotoHandler(
		photoService,
	)

	searchHandler := handler.NewSearchHandler(
		searchService,
	)

	app := fiber.New()

	// CORS
	app.Use(
		middleware.CORS(cfg.AllowSite),
	)

	// API documentation
	app.Get("/docs", func(c fiber.Ctx) error {
		return c.SendFile("./docs/index.html")
	})

	app.Get(
		"/docs/*",
		static.New("./docs"),
	)

	router.Setup(
		app,
		healthHandler,
		dataHandler,
		detailHandler,
		photoHandler,
		searchHandler,
	)

	address := ":" + cfg.AppPort

	logger.Info(
		"%s running on %s",
		cfg.AppName,
		address,
	)

	if err := app.Listen(address); err != nil {
		logger.Error(
			"Server stopped: %v",
			err,
		)

		panic(
			fmt.Errorf(
				"server error: %w",
				err,
			),
		)
	}
}
