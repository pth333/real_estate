package usecase

import (
	"context"
	"fmt"
	"strings"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/repo"
)

type AIService struct {
	repo repo.IAIRepository
}

type IAIService interface {
	GenerateContent(req dto.AIRequest) (dto.AIContentResponse, error)
}

func NewAIService(aiRepo repo.IAIRepository) IAIService {
	return &AIService{repo: aiRepo}
}

// GenerateContent dựng prompt từ thông tin BĐS theo văn phong rồi gọi AI qua repo
func (s *AIService) GenerateContent(req dto.AIRequest) (dto.AIContentResponse, error) {
	prompt, err := buildPrompt(req)
	if err != nil {
		return dto.AIContentResponse{}, err
	}

	return s.repo.GenerateContent(context.Background(), prompt)
}

// buildPrompt dựng prompt chi tiết từ các trường thông tin BĐS, yêu cầu AI trả JSON { title, description }
func buildPrompt(req dto.AIRequest) (string, error) {
	// Văn phong
	toneDesc := map[string]string{
		"lich_su":   "LỊCH SỰ, trang trọng, chuyên nghiệp",
		"tre_trung": "TRẺ TRUNG, gần gũi, năng động",
	}
	tone, ok := toneDesc[req.Tone]
	if !ok {
		return "", fmt.Errorf("văn phong không hợp lệ: %s", req.Tone)
	}

	// Hành động theo loại tin: bán hay cho thuê
	listingVerb := "Bán"
	if req.ListingType == "rent" {
		listingVerb = "Cho thuê"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Bạn là chuyên gia viết tin đăng bất động sản tiếng Việt.\n")
	fmt.Fprintf(&sb, "Viết tin đăng %s bất động sản với văn phong %s.\n", strings.ToLower(listingVerb), tone)
	fmt.Fprintf(&sb, "Sử dụng thông tin sau đây của BĐS:\n\n")

	// ── Các thông tin BĐS ──
	addField := func(label string, value any) {
		if value != nil && fmt.Sprintf("%v", value) != "" {
			fmt.Fprintf(&sb, "- %s: %v\n", label, value)
		}
	}

	addField("Loại tin", listingVerb)
	addField("Loại bất động sản", req.RealEstateType)
	addField("Dự án", req.ProjectName)
	addField("Địa chỉ", req.Address)
	addField("Diện tích", fmt.Sprintf("%.0f m²", req.Area))
	addField("Giá", formatPrice(req.Price, req.Unit))
	addField("Số phòng ngủ", formatQuantity(req.Bedrooms, "PN"))
	addField("Số phòng tắm", formatQuantity(req.Bathrooms, "WC"))
	addField("Giấy tờ pháp lý", req.LegalDocs)
	addField("Nội thất", req.Interior)
	addField("Hướng nhà", req.HouseDirection)
	addField("Hướng ban công", req.BalconyDirection)
	addField("Tên liên hệ", req.ContactName)
	addField("Số điện thoại", req.ContactPhone)

	// ── Yêu cầu output ──
	sb.WriteString("\n\nViết theo cấu trúc sau:\n")
	sb.WriteString("1. TIÊU ĐỀ ngắn gọn, súc tích (dưới 99 ký tự), nêu rõ: hành động, loại BĐS, vị trí, giá.\n")
	sb.WriteString("2. MÔ TẢ: đoạn mở đầu tóm tắt BĐS, sau đó liệt kê các đặc điểm, tiện ích xung quanh, pháp lý, và kêu gọi liên hệ.\n")
	sb.WriteString("\nTrả về DUY NHẤT JSON đúng định dạng, không kèm giải thích:\n")
	sb.WriteString(`{"title": "<tiêu đề>", "description": "<mô tả>"}`)

	return sb.String(), nil
}

// formatPrice định dạng giá theo đơn vị (vnd/usd/eur)
func formatPrice(price float64, unit string) string {
	if price <= 0 {
		return ""
	}
	unitLabel := map[string]string{
		"vnd": "VND",
		"usd": "USD",
		"eur": "EUR",
	}
	label := unitLabel[unit]
	if label == "" {
		label = "VND"
	}
	return fmt.Sprintf("%.0f %s", price, label)
}

// formatQuantity hiển thị số lượng + đơn vị, bỏ qua nếu = 0
func formatQuantity(count int, suffix string) string {
	if count <= 0 {
		return ""
	}
	return fmt.Sprintf("%d%s", count, suffix)
}
