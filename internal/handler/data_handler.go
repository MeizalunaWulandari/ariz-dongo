package handler

import (
	"encoding/json"

	"github.com/MeizalunaWulandari/ariz-dongo/internal/service"
	"github.com/gofiber/fiber/v3"
)

type DataHandler struct {
	service *service.DataService
}

func NewDataHandler(
	service *service.DataService,
) *DataHandler {
	return &DataHandler{
		service: service,
	}
}

func (h *DataHandler) GetData(c fiber.Ctx) error {
	date := c.Query("date")
	pagination := c.Query("pagination")
	region := c.Query("region")

	data, err := h.service.GetData(
		c.Context(),
		date,
		pagination,
		region,
	)

	if err != nil {
		return c.Status(
			fiber.StatusBadGateway,
		).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	var jsonData any

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return c.Status(
			fiber.StatusBadGateway,
		).JSON(fiber.Map{
			"success": false,
			"message": "invalid response from Fieldsa",
		})
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(jsonData)
}
