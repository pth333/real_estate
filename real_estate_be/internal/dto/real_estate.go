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
	City     string  `json:"city,omitempty"`
	Ward     string  `json:"ward,omitempty"`
	Slug     string  `json:"slug,omitempty"`
}

type RealEstateResponse struct {
	ID          uint64   `json:"id" gorm:"column:id"`
	Title       string   `json:"title" gorm:"column:title"`
	PriceVND    float64  `json:"price_vnd" gorm:"column:price_vnd"`
	Address     string   `json:"address" gorm:"column:address"`
	District    string   `json:"district" gorm:"column:district"`
	City        string   `json:"city" gorm:"column:city"`
	Acreage     float64  `json:"acreage" gorm:"column:acreage"`
	PricePerM2  float64  `json:"price_per_m2" gorm:"column:price_per_m2"`
	Images      []string `json:"images,omitempty" gorm:"-"`
	Bedrooms    *int     `json:"bedrooms,omitempty" gorm:"column:bedrooms"`
	Bathrooms   *int     `json:"bathrooms,omitempty" gorm:"column:bathrooms"`
	Description string   `json:"description,omitempty" gorm:"column:description"`
	AgentName   string   `json:"agent_name,omitempty" gorm:"column:agent_name"`
	AgentPhone  string   `json:"agent_phone,omitempty" gorm:"column:agent_phone"`
	Badge       string   `json:"badge,omitempty" gorm:"-"`
	CreatedAt   string   `json:"created_at" gorm:"column:created_at"`
	// Scan nội bộ — không export ra JSON, chứa URL ảnh pipe-separated từ GROUP_CONCAT
	ImageURLs string `json:"-" gorm:"column:image_urls"`
}

type ProvinceResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// TopCityResponse — 1 thành phố trong danh sách "Bất động sản theo khu vực"
type TopCityResponse struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Count int64  `json:"count"`
	Image string `json:"image"`
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
