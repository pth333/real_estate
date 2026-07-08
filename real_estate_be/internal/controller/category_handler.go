package controller

import (
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	service usecase.ICategoryService
}

func NewCategoryHandler(service usecase.ICategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetAll()
	if err != nil {
		return response.InternalServerError(c, "Failed to get categories", err.Error())
	}

	categoriesResponse := h.service.BuildCategoriesResponse(categories)
	return response.OK(c, categoriesResponse)
}
