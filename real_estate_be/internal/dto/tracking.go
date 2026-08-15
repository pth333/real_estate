package dto

// TrackingSearchRequest payload cho tracking search
type TrackingSearchRequest struct {
	Query     string  `json:"query"`
	UserID    string  `json:"user_id,omitempty"`    // Chấp nhận rỗng nếu là Guest
	SessionID string  `json:"session_id,omitempty"` // UUID của Guest
	Filters   Filters `json:"filters"`
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

// TrackingViewRequest payload cho tracking thời gian xem tin chi tiết
type TrackingViewRequest struct {
	RealEstateID    uint64 `json:"real_estate_id"`
	DurationSeconds int    `json:"duration_seconds"`
	SessionID       string `json:"session_id"`
	UserID          uint64 `json:"user_id,omitempty"`
}

// MergeSessionRequest payload khi sáp nhập session từ Guest sang User đã đăng nhập
type MergeSessionRequest struct {
	SessionID string `json:"session_id"`
}
