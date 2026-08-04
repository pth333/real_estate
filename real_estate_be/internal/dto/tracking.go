package dto

// TrackingSearchRequest payload cho tracking search
type TrackingSearchRequest struct {
	Query   string  `json:"query"`
	UserID  string  `json:"user_id"`
	Filters Filters `json:"filters"`
}

type Filters struct {
	Location   *Location   `json:"location,omitempty"`
	PriceRange *PriceRange `json:"price_range,omitempty"`
}

type Location struct {
	District *string `json:"district,omitempty"`
	City     *string `json:"city,omitempty"`
	Ward     *string `json:"ward,omitempty"`
	Street   *string `json:"street,omitempty"`
}

type PriceRange struct {
	MinPrice *float64 `json:"min_price,omitempty"`
	MaxPrice *float64 `json:"max_price,omitempty"`
}
