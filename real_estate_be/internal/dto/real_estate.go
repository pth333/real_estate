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

// CreateRealEstateRequest — payload tạo tin đăng từ FE
type CreateRealEstateRequest struct {
	ListingType      string   `json:"listing_type"`
	Province         string   `json:"province"`
	Ward             string   `json:"ward"`
	DetailAddress    string   `json:"detail_address"`
	RealEstateType   string   `json:"real_estate_type"`
	Area             float64  `json:"area"`
	Price            float64  `json:"price"`
	Unit             string   `json:"unit"`
	LegalDocs        string   `json:"legal_docs"`
	Interior         string   `json:"interior"`
	BathroomCount    *int     `json:"bathroom_count"`
	BedroomCount     *int     `json:"bedroom_count"`
	HouseDirection   string   `json:"house_direction"`
	BalconyDirection string   `json:"balcony_direction"`
	MoveInTime       string   `json:"move_in_time"`
	PriceElectricity *float64 `json:"price_electricity"`
	PriceWater       *float64 `json:"price_water"`
	PriceInternet    *float64 `json:"price_internet"`
	Amenities        []string `json:"amenities"`
	ContactName      string   `json:"contact_name"`
	ContactEmail     string   `json:"contact_email"`
	ContactPhone     string   `json:"contact_phone"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	// Id ảnh/video đã upload (từ /upload/confirm), để liên kết với tin đăng
	ImageIDs []uint64 `json:"image_ids"`
}
