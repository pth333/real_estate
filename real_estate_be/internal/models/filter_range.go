package model

// FilterRange — khoảng giá/diện tích dùng cho slug SEO. Thay cho việc parse
// chuỗi giá thủ công: URL "gia-1-den-3-ty" → WHERE slug = 'gia-1-den-3-ty'
// → lấy min_val/max_val để query BETWEEN (giống batdongsan.com.vn).
type FilterRange struct {
	ID     uint    `gorm:"primaryKey"`
	// Type: 'price' hoặc 'area' (giá / diện tích)
	Type   string  `gorm:"column:type;size:20;index"`
	Label  string  `gorm:"column:label;size:100"`
	Slug   string  `gorm:"column:slug;size:100;uniqueIndex"`
	// Nil = không giới hạn ở phía đó
	MinVal *float64 `gorm:"column:min_val"`
	MaxVal *float64 `gorm:"column:max_val"`
}