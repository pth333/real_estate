package controller

import (
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"
	"real_estate_be/pkg/jwt"

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

	// Decode user_id từ access token nếu có (route này public, không bắt buộc đăng nhập)
	var userID uint64
	if authHeader := c.Get("Authorization"); authHeader != "" {
		tokenStr := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr = authHeader[7:]
		}
		if claims, err := jwt.ParseAccessToken(tokenStr); err == nil {
			userID = claims.UserID
		}
	}

	return response.OK(c, fiber.Map{
		"user_id":    userID,
		"categories": categoriesResponse,
	})
}
