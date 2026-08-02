package dto

// AIRequest là payload client gửi lên khi yêu cầu AI tạo nội dung
type AIRequest struct {
	Tone string `json:"tone"` // văn phong: "lich_su" hoặc "tre_trung"
}

// AIContentResponse là nội dung AI trả về
type AIContentResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
