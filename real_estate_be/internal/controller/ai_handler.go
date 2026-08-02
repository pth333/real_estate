package controller

import (
	"real_estate_be/internal/dto"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type AIHandler struct {
	service usecase.IAIService
}

func NewAIHandler(service usecase.IAIService) *AIHandler {
	return &AIHandler{service: service}
}

// GenerateContent xử lý yêu cầu AI tạo tiêu đề & mô tả theo văn phong
func (h *AIHandler) GenerateContent(c *fiber.Ctx) error {
	var req dto.AIRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu không hợp lệ", err.Error())
	}
	if req.Tone == "" {
		return response.BadRequest(c, "Thiếu văn phong (tone)", nil)
	}

	content, err := h.service.GenerateContent(req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, content)
}
