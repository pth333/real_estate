package dto

// AIRequest là payload client gửi lên khi yêu cầu AI tạo nội dung.
// Chứa thông tin BĐS + văn phong để AI dựa vào đó viết tiêu đề & mô tả.
type AIRequest struct {
	Tone string `json:"tone"` // văn phong: "lich_su" hoặc "tre_trung"

	// ── Thông tin BĐS ──
	ListingType    string `json:"listing_type"`    // loại tin: "sell" hoặc "rent"
	RealEstateType string `json:"real_estate_type"` // loại BĐS: căn hộ, nhà riêng, đất...
	ProjectName    string `json:"project_name"`    // tên dự án (nếu có)
	Address        string `json:"address"`         // địa chỉ đầy đủ
	Area           float64 `json:"area"`           // diện tích (m²)
	Price          float64 `json:"price"`          // giá
	Unit           string  `json:"unit"`           // đơn vị giá: vnd/usd/eur
	Bedrooms       int     `json:"bedrooms"`       // số phòng ngủ
	Bathrooms      int     `json:"bathrooms"`      // số phòng tắm
	LegalDocs      string  `json:"legal_docs"`     // giấy tờ pháp lý
	Interior       string  `json:"interior"`       // nội thất
	HouseDirection string  `json:"house_direction"` // hướng nhà
	BalconyDirection string `json:"balcony_direction"` // hướng ban công
	ContactName    string  `json:"contact_name"`   // tên người liên hệ
	ContactPhone   string  `json:"contact_phone"`  // số điện thoại liên hệ
}

// AIContentResponse là nội dung AI trả về
type AIContentResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
