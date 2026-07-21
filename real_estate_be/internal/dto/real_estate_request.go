package dto

type RealEstateSearchRequest struct {
	Slug     string  `json:"slug,omitempty"`
	Page     int     `json:"page"`
	CursorID int64   `json:"cursor_id,omitempty"` // id bản ghi cuối trang trước dùng để keyset
	Size     int     `json:"size"`
	Filter   Filter  `json:"filter,omitempty"`
}

type Filter struct {
	MinPrice float64 `json:"min_price,omitempty"`
	MaxPrice float64 `json:"max_price,omitempty"`
	District string  `json:"district,omitempty"`
}
