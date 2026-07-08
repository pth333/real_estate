package dto

type RealEstateSearchRequest struct {
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Filter Filter `json:"filter,omitempty"`
}

type Filter struct {
	MinPrice float64 `json:"min_price,omitempty"`
	MaxPrice float64 `json:"max_price,omitempty"`
	District string  `json:"district,omitempty"`
}
