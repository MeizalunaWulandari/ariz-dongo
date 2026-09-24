package router

import (
	"github.com/MeizalunaWulandari/ariz-dongo/internal/handler"
	"github.com/gofiber/fiber/v3"
)

func Setup(
	app *fiber.App,
	healthHandler *handler.HealthHandler,
	dataHandler *handler.DataHandler,
	detailHandler *handler.DetailHandler,
	photoHandler *handler.PhotoHandler,
	searchHandler *handler.SearchHandler,

) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get(
		"/health",
		healthHandler.Check,
	)

	data := v1.Group("/data")

	data.Get(
		"/",
		dataHandler.GetData,
	)

	v1.Get(
		"/search/:circuitId",
		searchHandler.Search,
	)

	// HARUS sebelum /:workOrderNumber
	data.Get(
		"/photos/:workOrderNumber",
		photoHandler.GetPhotos,
	)

	data.Get(
		"/:workOrderNumber",
		detailHandler.GetDetail,
	)
}
