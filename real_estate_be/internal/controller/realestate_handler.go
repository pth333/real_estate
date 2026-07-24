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

func (h *RealEstateHandler) ListCity(c *fiber.Ctx) error {
	provinces, err := h.service.GetListCity()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	options := make([]dto.ProvinceResponse, len(provinces))
	for i, province := range provinces {
		options[i] = dto.ProvinceResponse{
			Name: province.Name,
			Code: province.Code,
		}
	}

	return c.JSON(fiber.Map{
		"data": options,
	})
}

func (h *RealEstateHandler) ListWard(c *fiber.Ctx) error {
	provinceCode := c.Query("code")
	if provinceCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing district code",
		})
	}

	wards, err := h.service.GetListWard(provinceCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	options := make([]dto.ProvinceResponse, len(wards))
	for i, ward := range wards {
		options[i] = dto.ProvinceResponse{
			Name: ward.Name,
			Code: ward.Code,
		}
	}

	return c.JSON(fiber.Map{
		"data": options,
	})
}
