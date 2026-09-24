package handler

import (
	"github.com/MeizalunaWulandari/ariz-dongo/internal/service"
	"github.com/gofiber/fiber/v3"
)

type PhotoHandler struct {
	service *service.PhotoService
}

func NewPhotoHandler(
	service *service.PhotoService,
) *PhotoHandler {
	return &PhotoHandler{
		service: service,
	}
}

func (h *PhotoHandler) GetPhotos(
	c fiber.Ctx,
) error {

	workOrderNumber :=
		c.Params("workOrderNumber")

	data, err := h.service.GetPhotos(
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
