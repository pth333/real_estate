package controller

import (
	"real_estate_be/internal/dto"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type RealEstateHandler struct {
	service usecase.IRealEstateService
}

func NewRealEstateHandler(service usecase.IRealEstateService) *RealEstateHandler {
	return &RealEstateHandler{service: service}
}

func (h *RealEstateHandler) List(c *fiber.Ctx) error {
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

func (h *RealEstateHandler) ListRealEsateByCategory(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid category slug",
		})
	}

	var req dto.RealEstateSearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	req.Slug = slug
	data, total, err := h.service.ListRealEstateByCategory(req)
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
