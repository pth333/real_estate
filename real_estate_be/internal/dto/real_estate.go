package dto

type RealEstateSearchRequest struct {
	Slug   string `json:"slug,omitempty"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Filter Filter `json:"filter,omitempty"`
	Search string `json:"search,omitempty"`
}

type Filter struct {
	MinPrice   float64 `json:"min_price,omitempty"`
	MaxPrice   float64 `json:"max_price,omitempty"`
	MinAcreage float64 `json:"min_acreage,omitempty"`
	MaxAcreage float64 `json:"max_acreage,omitempty"`
	District   string  `json:"district,omitempty"`
	City       string  `json:"city,omitempty"`
	Ward       string  `json:"ward,omitempty"`
	Slug       string  `json:"slug,omitempty"`

	// ── Bộ lọc nâng cao (query string) ──
	// CategoryID: loại bất động sản (lấy id từ bảng categories)
	Bedrooms         *int   `json:"bedrooms,omitempty"`
	Bathrooms        *int   `json:"bathrooms,omitempty"`
	HouseDirection   string `json:"house_direction,omitempty"`
	BalconyDirection string `json:"balcony_direction,omitempty"`
	LegalDocs        string `json:"legal_docs,omitempty"`
	Interior         string `json:"interior,omitempty"`
}

type RealEstateResponse struct {
	ID             uint64          `json:"id" gorm:"column:id"`
	Title          string          `json:"title" gorm:"column:title"`
	Slug           string          `json:"slug,omitempty" gorm:"column:slug"`
	RealEstateType string          `json:"real_estate_type,omitempty" gorm:"column:real_estate_type"`
	PriceVND       float64         `json:"price_vnd" gorm:"column:price_vnd"`
	Address        string          `json:"address" gorm:"column:address"`
	District       string          `json:"district" gorm:"column:district"`
	City           string          `json:"city" gorm:"column:city"`
	Acreage        float64         `json:"acreage" gorm:"column:acreage"`
	PricePerM2     float64         `json:"price_per_m2" gorm:"column:price_per_m2"`
	Images         []ImageResponse `json:"images,omitempty" gorm:"-"`
	Bedrooms       *int            `json:"bedrooms,omitempty" gorm:"column:bedrooms"`
	Bathrooms      *int            `json:"bathrooms,omitempty" gorm:"column:bathrooms"`
	Description    string          `json:"description,omitempty" gorm:"column:description"`
	AgentName      string          `json:"agent_name,omitempty" gorm:"column:agent_name"`
	AgentEmail     string          `json:"agent_email,omitempty" gorm:"column:agent_email"`
	AgentPhone     string          `json:"agent_phone,omitempty" gorm:"column:agent_phone"`
	Badge          string          `json:"badge,omitempty" gorm:"-"`
	CreatedAt      string          `json:"created_at" gorm:"column:created_at"`

	// ── Thông tin bổ sung mục Đặc điểm BĐS ──
	HouseDirection   string   `json:"house_direction,omitempty" gorm:"column:house_direction"`
	BalconyDirection string   `json:"balcony_direction,omitempty" gorm:"column:balcony_direction"`
	Floors           *int     `json:"floors,omitempty" gorm:"column:floors"`
	LegalDocs        string   `json:"legal_docs,omitempty" gorm:"column:legal_docs"`
	Interior         string   `json:"interior,omitempty" gorm:"column:interior"`
	PriceElectricity *float64 `json:"price_electricity,omitempty" gorm:"column:price_electricity"`
	PriceWater       *float64 `json:"price_water,omitempty" gorm:"column:price_water"`
	PriceInternet    *float64 `json:"price_internet,omitempty" gorm:"column:price_internet"`
	Latitude         *float64 `json:"latitude,omitempty" gorm:"column:latitude"`
	Longitude        *float64 `json:"longitude,omitempty" gorm:"column:longitude"`
	// Tiện ích lưu JSON string trong DB (`["camera"]`), DTO trả mảng string
	Amenities []string `json:"amenities,omitempty" gorm:"-"`
	// AmenitiesRaw nhận giá trị raw từ cột amenities (JSON string) khi scan
	AmenitiesRaw string `json:"-" gorm:"column:amenities"`

	// Nhận dữ liệu từ SQL
	ImageURLsRaw string `json:"-" gorm:"column:image_urls"`

	// Trả FE
	ImageURLs []string `json:"image_urls,omitempty" gorm:"-"`
}

type ProvinceResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type ProjectResponse struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug,omitempty"`
	Status      string   `json:"status,omitempty"`
	FullAddress string   `json:"full_address,omitempty"`
	TotalAreaHA *float64 `json:"total_area_ha,omitempty"`
	TotalUnits  *uint32  `json:"total_units,omitempty"`
	PriceMin    *float64 `json:"price_min,omitempty"`
	PriceMax    *float64 `json:"price_max,omitempty"`
	ViewCount   uint32   `json:"view_count"`
	Thumbnail   string   `json:"thumbnail,omitempty"`
}

type ListProjectRequest struct {
	Province string `json:"province"`
	Ward     string `json:"ward"`
}

// TopCityResponse — 1 thành phố trong danh sách "Bất động sản theo khu vực"
// CategorySlug: category đầu tiên trong menu; CitySlug: tên thành phố đã slug hóa.
// FE ghép → /{CategorySlug}/{CitySlug} để server-driven.
type TopCityResponse struct {
	Name         string `json:"name"`
	CategorySlug string `json:"category_slug"`
	CitySlug     string `json:"city_slug"`
	Count        int64  `json:"count"`
	Image        string `json:"image"`
}

// CreateRealEstateRequest — request tạo tin đăng từ FE
type CreateRealEstateRequest struct {
	ID               *uint64  `json:"id,omitempty"`
	ListingType      string   `json:"listing_type"`
	Province         string   `json:"province"`
	Ward             string   `json:"ward"`
	DetailAddress    string   `json:"detail_address"`
	ProjectID        *uint64  `json:"project_id"`
	RealEstateType   string   `json:"real_estate_type"`
	Area             float64  `json:"area"`
	PricePerM2       float64  `json:"price_per_m2"`
	Unit             string   `json:"unit"`
	LegalDocs        string   `json:"legal_docs"`
	Interior         string   `json:"interior"`
	BathroomCount    *int     `json:"bathroom_count"`
	BedroomCount     *int     `json:"bedroom_count"`
	FloorCount       *int     `json:"floor_count"`
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
	Latitude         *float64 `json:"latitude,omitempty"`
	Longitude        *float64 `json:"longitude,omitempty"`
	// Id ảnh/video đã upload, để liên kết với tin đăng
	Images []ImageResponse `json:"images"`
}
