package model

import (
	"time"

	"gorm.io/gorm"
)

// RealEstateProject đại diện cho bảng real_estate_project trong database.
// Các trường địa lý (province_id, district_id, ward_id) chỉ lưu id,
// không tạo relationship/FK vì dữ liệu có thể đến từ crawler.
type RealEstateProject struct {
	ID uint64 `gorm:"primaryKey"`

	// Tên và định danh dự án
	Name            string `gorm:"column:name;index"`
	Slug            string `gorm:"column:slug;uniqueIndex"`
	AlternativeName string `gorm:"column:alternative_name"`
	CategoryID      *int64 `gorm:"column:category_id;index" json:"category_id"` // Liên kết với danh mục dự án trên Menu

	// Trạng thái dự án
	Status string `gorm:"column:status;index"`

	// Địa chỉ dự án — lưu MÃ tỉnh/phường dạng string (giữ số 0 đầu, VD "000331")
	// để map đúng với bảng provinces/wards khi hiển thị name vị trí.
	FullAddress  string  `gorm:"column:full_address"`
	ProvinceCode string  `gorm:"column:province_code;index"`
	WardCode     string  `gorm:"column:ward_code;index"`
	Latitude     *float64 `gorm:"column:latitude"`
	Longitude    *float64 `gorm:"column:longitude"`

	// Quy mô dự án
	TotalAreaHA         *float64 `gorm:"column:total_area_ha"`
	ConstructionDensity *float64 `gorm:"column:construction_density"`
	TotalUnits          *uint32  `gorm:"column:total_units"`
	TotalFloors         *uint32  `gorm:"column:total_floors"`
	TotalBlocks         *uint32  `gorm:"column:total_blocks"`
	ExpectedPopulation  *uint32  `gorm:"column:expected_population"`

	// Giá dự án
	PriceMin      *float64 `gorm:"column:price_min"`
	PriceMax      *float64 `gorm:"column:price_max"`
	PricePerM2Min *float64 `gorm:"column:price_per_m2_min"`
	PricePerM2Max *float64 `gorm:"column:price_per_m2_max"`
	InvestorID    *uint64  `gorm:"column:investor_id"`
	LegalStatus   string   `gorm:"column:legal_status"`

	// Mốc thời gian
	ConstructionStartDate *time.Time `gorm:"column:construction_start_date"`
	HandoverDate          *time.Time `gorm:"column:handover_date"`

	// Thống kê & SEO
	ListingCount    int64  `gorm:"column:listing_count;default:0"`
	ViewCount       int64  `gorm:"column:view_count;default:0"`
	MetaTitle       string `gorm:"column:meta_title"`
	MetaDescription string `gorm:"column:meta_description;type:text"`

	// Soft delete: GORM tự fill deleted_at
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
