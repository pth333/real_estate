package controller

import (
	"real_estate_be/internal/dto"
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

// ===== Summary =====
func (h *DashboardHandler) Summary(c *fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")

	result, err := h.service.Summary(from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//

	return c.JSON(result)
}

func (h *DashboardHandler) ListRealEstate(c *fiber.Ctx) error {
	var req dto.RealEstateSearchRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	data, total, err := h.service.ListRealEstate(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"total": total,
		"data":  data,
	})
}
