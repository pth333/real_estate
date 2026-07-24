package dto

type RealEstateSearchRequest struct {
	Slug   string `json:"slug,omitempty"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Filter Filter `json:"filter,omitempty"`
	Search string `json:"search,omitempty"`
}

type Filter struct {
	MinPrice float64 `json:"min_price,omitempty"`
	MaxPrice float64 `json:"max_price,omitempty"`
	District string  `json:"district,omitempty"`
	Slug     string  `json:"slug,omitempty"`
}

type RealEstateResponse struct {
	ID               uint64   `json:"id"`
	Title            string   `json:"title"`
	PriceVND         float64  `json:"price_vnd"`
	Address          string   `json:"address"`
	District         string   `json:"district"`
	City             string   `json:"city"`
	Acreage          float64  `json:"acreage"`
	PricePerM2       float64  `json:"price_per_m2"`
	TypeOfRealEstate string   `json:"type_of_real_estate"`
	Images           []string `json:"images,omitempty"`
	Bedrooms         *int     `json:"bedrooms,omitempty"`
	Bathrooms        *int     `json:"bathrooms,omitempty"`
	Description      string   `json:"description,omitempty"`
	AgentName        string   `json:"agent_name,omitempty"`
	AgentPhone       string   `json:"agent_phone,omitempty"`
	Badge            string   `json:"badge,omitempty"`
	CreatedAt        string   `json:"created_at"`
}

type ProvinceResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
