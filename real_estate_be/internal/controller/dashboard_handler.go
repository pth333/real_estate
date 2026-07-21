package controller

import (
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service usecase.IDashboardService
}

func NewDashboardHandler(service usecase.IDashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

func (h *DashboardHandler) Summary(c *fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")

	result, err := h.service.Summary(from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
