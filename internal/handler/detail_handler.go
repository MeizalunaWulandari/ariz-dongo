package handler

import (
	"github.com/MeizalunaWulandari/ariz-dongo/internal/service"
	"github.com/gofiber/fiber/v3"
)

type DetailHandler struct {
	service *service.DetailService
}

func NewDetailHandler(
	service *service.DetailService,
) *DetailHandler {
	return &DetailHandler{
		service: service,
	}
}

func (h *DetailHandler) GetDetail(
	c fiber.Ctx,
) error {

	workOrderNumber := c.Params(
		"workOrderNumber",
	)

	data, err := h.service.GetDetail(
		c.Context(),
		workOrderNumber,
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
