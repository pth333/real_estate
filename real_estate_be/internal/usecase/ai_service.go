package usecase

import (
	"context"
	"fmt"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/repo"
)

type AIService struct {
	repo repo.IAIRepository
}

type IAIService interface {
	GenerateContent(tone string) (dto.AIContentResponse, error)
}

func NewAIService(aiRepo repo.IAIRepository) IAIService {
	return &AIService{repo: aiRepo}
}

// GenerateContent dựng prompt theo văn phong rồi gọi AI qua repo
func (s *AIService) GenerateContent(tone string) (dto.AIContentResponse, error) {
	prompt, err := buildPrompt(tone)
	if err != nil {
		return dto.AIContentResponse{}, err
	}

	return s.repo.GenerateContent(context.Background(), prompt)
}

// buildPrompt dựng prompt yêu cầu AI trả JSON { title, description }
func buildPrompt(tone string) (string, error) {
	switch tone {
	case "lich_su":
		return `Viết tin đăng bán bất động sản với văn phong LỊCH SỰ, trang trọng.
Trả về JSON đúng định dạng: {"title": "<tiêu đề>", "description": "<mô tả>"}`, nil
	case "tre_trung":
		return `Viết tin đăng cho thuê bất động sản với văn phong TRẺ TRUNG, gần gũi.
Trả về JSON đúng định dạng: {"title": "<tiêu đề>", "description": "<mô tả>"}`, nil
	default:
		return "", fmt.Errorf("văn phong không hợp lệ: %s", tone)
	}
}
