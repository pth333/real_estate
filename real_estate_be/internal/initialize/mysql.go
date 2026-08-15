package initialize

import (
	"fmt"
	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() {
	m := global.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.DBName,
	)
	// DisableForeignKeyConstraintWhenMigrating: giữ relationship field để dùng Preload,
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Image{})
	db.AutoMigrate(&model.RealEstate{})
	db.AutoMigrate(&model.Province{})
	db.AutoMigrate(&model.SearchHistory{})
	db.AutoMigrate(&model.ViewHistory{})
	db.AutoMigrate(&model.RealEstateProject{})
	db.AutoMigrate(&model.FilterRange{})

	if err != nil {
		panic(err)
	}

	// Seed filter_ranges nếu bảng rỗng — menu giá/diện tích cho slug SEO.
	// FE build URL theo đúng slug này (server-driven), slug phải khớp 1-1.
	seedFilterRanges(db)

	global.DB = db
}

// seedFilterRanges chèn menu khoảng giá (price) + diện tích (area) khi bảng
// filter_ranges chưa có dữ liệu. Đơn vị: price = VNĐ, area = m².
// NULL ở một đầu = không giới hạn phía đó (VD "Trên 10 tỷ" không có max).
func seedFilterRanges(db *gorm.DB) {
	var count int64
	db.Model(&model.FilterRange{}).Count(&count)
	if count > 0 {
		return
	}

	// Helpers dùng pointer để biểu diễn min/max "không giới hạn"
	toPtr := func(v float64) *float64 { return &v }

	ranges := []model.FilterRange{
		// Giá (type=price) — VNĐ
		{Type: "price", Label: "Dưới 1 tỷ", Slug: "gia-duoi-1-ty", MaxVal: toPtr(1_000_000_000)},
		{Type: "price", Label: "Từ 1 đến 3 tỷ", Slug: "gia-1-den-3-ty", MinVal: toPtr(1_000_000_000), MaxVal: toPtr(3_000_000_000)},
		{Type: "price", Label: "Từ 3 đến 5 tỷ", Slug: "gia-3-den-5-ty", MinVal: toPtr(3_000_000_000), MaxVal: toPtr(5_000_000_000)},
		{Type: "price", Label: "Từ 5 đến 10 tỷ", Slug: "gia-5-den-10-ty", MinVal: toPtr(5_000_000_000), MaxVal: toPtr(10_000_000_000)},
		{Type: "price", Label: "Trên 10 tỷ", Slug: "gia-tren-10-ty", MinVal: toPtr(10_000_000_000)},

		// Diện tích (type=area) — m²
		{Type: "area", Label: "Dưới 30m²", Slug: "dien-tich-duoi-30", MaxVal: toPtr(30)},
		{Type: "area", Label: "Từ 30 đến 50m²", Slug: "dien-tich-30-50", MinVal: toPtr(30), MaxVal: toPtr(50)},
		{Type: "area", Label: "Từ 50 đến 100m²", Slug: "dien-tich-50-100", MinVal: toPtr(50), MaxVal: toPtr(100)},
		{Type: "area", Label: "Từ 100 đến 200m²", Slug: "dien-tich-100-200", MinVal: toPtr(100), MaxVal: toPtr(200)},
		{Type: "area", Label: "Trên 200m²", Slug: "dien-tich-tren-200", MinVal: toPtr(200)},
	}

	if err := db.Create(&ranges).Error; err != nil {
		panic(err)
	}
}
