package handler

import (
	"github.com/MeizalunaWulandari/ariz-dongo/internal/service"
	"github.com/gofiber/fiber/v3"
)

type HealthHandler struct {
	service *service.HealthService
}

func NewHealthHandler(service *service.HealthService) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

func (h *HealthHandler) Check(c fiber.Ctx) error {
	data := h.service.Check()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "API is running",
		"data":    data,
	})
}
