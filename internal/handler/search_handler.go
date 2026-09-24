package handler

import (
	"github.com/MeizalunaWulandari/ariz-dongo/internal/service"
	"github.com/gofiber/fiber/v3"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(
	service *service.SearchService,
) *SearchHandler {
	return &SearchHandler{
		service: service,
	}
}

func (h *SearchHandler) Search(
	c fiber.Ctx,
) error {

	circuitID := c.Params("circuitId")

	data, err :=
		h.service.Search(
			c.Context(),
			circuitID,
		)

	if err != nil {

		return c.Status(
			fiber.StatusBadGateway,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}
