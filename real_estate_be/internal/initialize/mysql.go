package initialize

import (
	"fmt"
	"log"
	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	"time"

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

	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			break
		}
		fmt.Printf("⏳ [MySQL] Waiting for DB... attempt %d/5\n", i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		panic(fmt.Sprintf("❌ [MySQL] Cannot connect after 5 attempts: %v", err))
	}

	// db.AutoMigrate(&model.User{})
	// db.AutoMigrate(&model.Category{})
	// db.AutoMigrate(&model.Image{})
	// db.AutoMigrate(&model.RealEstate{})
	// db.AutoMigrate(&model.Province{})
	// db.AutoMigrate(&model.SearchHistory{})
	// db.AutoMigrate(&model.ViewHistory{})
	// db.AutoMigrate(&model.RealEstateProject{})
	// db.AutoMigrate(&model.FilterRange{})
	// db.AutoMigrate(&model.Favorite{})
	// db.AutoMigrate(&model.ImageProject{})

	// if db.Migrator().HasColumn(&model.RealEstateProject{}, "province_id") {
	// 	_ = db.Migrator().DropColumn(&model.RealEstateProject{}, "province_id")
	// }
	// if db.Migrator().HasColumn(&model.RealEstateProject{}, "ward_id") {
	// 	_ = db.Migrator().DropColumn(&model.RealEstateProject{}, "ward_id")
	// }

	// seedFilterRanges(db)

	global.DB = db
}

// MigrateDb tự động migrate các bảng.
func MigrateDb(db *gorm.DB) {
	if err := db.AutoMigrate(
		&model.User{},
		&model.RealEstate{},
		&model.Notification{},
		// Bảng dự án — cần có để tạo cột sold_units (tồn kho đã bán)
		&model.RealEstateProject{},
		// Luồng đặt cọc escrow
		&model.Deposit{},
		&model.Transaction{},
		&model.Dispute{},
		&model.BrokerRating{},
		&model.NotificationLog{},
		&model.DepositPolicy{},
	); err != nil {
		log.Fatalf("❌ DB migration failed: %v", err)
	}
	seedAdminUser(db)
	seedDepositPolicies(db)
	log.Println("✅ DB migration completed")
}

// seedDepositPolicies chèn chính sách mức cọc theo khoảng giá khi bảng còn trống.
// Khoản cọc này là cọc GIỮ LỊCH XEM NHÀ (chống bùng + trả công môi giới),
// không phải cọc mua bán → mức tiền nhỏ và cố định theo phân khúc, không scale theo giá nhà.
// Khoảng giá dạng nửa khoảng [price_min, price_max): NULL = không giới hạn phía đó.
func seedDepositPolicies(db *gorm.DB) {
	var count int64
	db.Model(&model.DepositPolicy{}).Count(&count)
	if count > 0 {
		return
	}

	billion := func(value float64) *float64 { return &value }

	policies := []model.DepositPolicy{
		{Label: "Dưới 1 tỷ", PriceMax: billion(1_000_000_000), DepositAmount: 2_000_000, BrokerFee: 200_000},
		{Label: "Từ 1 đến 3 tỷ", PriceMin: billion(1_000_000_000), PriceMax: billion(3_000_000_000), DepositAmount: 3_000_000, BrokerFee: 300_000},
		{Label: "Từ 3 đến 5 tỷ", PriceMin: billion(3_000_000_000), PriceMax: billion(5_000_000_000), DepositAmount: 5_000_000, BrokerFee: 500_000},
		{Label: "Từ 5 đến 10 tỷ", PriceMin: billion(5_000_000_000), PriceMax: billion(10_000_000_000), DepositAmount: 8_000_000, BrokerFee: 800_000},
		{Label: "Trên 10 tỷ", PriceMin: billion(10_000_000_000), DepositAmount: 15_000_000, BrokerFee: 1_500_000},
	}

	if err := db.Create(&policies).Error; err != nil {
		panic(err)
	}
	log.Printf("✅ [Deposit] đã seed %d chính sách mức cọc theo khoảng giá", len(policies))
}

// seedAdminUser nâng quyền ADMIN cho tài khoản cấu hình ở admin.email.
// Dùng để có sẵn 1 admin xử lý tranh chấp mà không phải sửa DB thủ công.
func seedAdminUser(db *gorm.DB) {
	email := global.Config.Admin.Email
	if email == "" {
		log.Println("ℹ️ [Admin] chưa cấu hình admin.email — bỏ qua bước nâng quyền ADMIN")
		return
	}

	result := db.Model(&model.User{}).
		Where("email = ?", email).
		Update("role", model.RoleAdmin)
	if result.Error != nil {
		log.Printf("⚠️ [Admin] nâng quyền ADMIN cho %s thất bại: %v", email, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		log.Printf("⚠️ [Admin] không tìm thấy tài khoản %s để nâng quyền ADMIN", email)
		return
	}
	log.Printf("✅ [Admin] đã nâng quyền ADMIN cho %s", email)
}

// seedFilterRanges chèn menu khoảng giá (price) + diện tích (area) khi bảng
// filter_ranges chưa có dữ liệu. Đơn vị: price = VNĐ, area = m².
// NULL ở một đầu = không giới hạn phía đó (VD "Trên 10 tỷ" không có max).
// func seedFilterRanges(db *gorm.DB) {
// 	var count int64
// 	db.Model(&model.FilterRange{}).Count(&count)
// 	if count > 0 {
// 		return
// 	}

// 	// Helpers dùng pointer để biểu diễn min/max "không giới hạn"
// 	toPtr := func(v float64) *float64 { return &v }

// 	ranges := []model.FilterRange{
// 		// Giá (type=price) — VNĐ
// 		{Type: "price", Label: "Dưới 1 tỷ", Slug: "gia-duoi-1-ty", MaxVal: toPtr(1_000_000_000)},
// 		{Type: "price", Label: "Từ 1 đến 3 tỷ", Slug: "gia-1-den-3-ty", MinVal: toPtr(1_000_000_000), MaxVal: toPtr(3_000_000_000)},
// 		{Type: "price", Label: "Từ 3 đến 5 tỷ", Slug: "gia-3-den-5-ty", MinVal: toPtr(3_000_000_000), MaxVal: toPtr(5_000_000_000)},
// 		{Type: "price", Label: "Từ 5 đến 10 tỷ", Slug: "gia-5-den-10-ty", MinVal: toPtr(5_000_000_000), MaxVal: toPtr(10_000_000_000)},
// 		{Type: "price", Label: "Trên 10 tỷ", Slug: "gia-tren-10-ty", MinVal: toPtr(10_000_000_000)},

// 		// Diện tích (type=area) — m²
// 		{Type: "area", Label: "Dưới 30m²", Slug: "dien-tich-duoi-30", MaxVal: toPtr(30)},
// 		{Type: "area", Label: "Từ 30 đến 50m²", Slug: "dien-tich-30-50", MinVal: toPtr(30), MaxVal: toPtr(50)},
// 		{Type: "area", Label: "Từ 50 đến 100m²", Slug: "dien-tich-50-100", MinVal: toPtr(50), MaxVal: toPtr(100)},
// 		{Type: "area", Label: "Từ 100 đến 200m²", Slug: "dien-tich-100-200", MinVal: toPtr(100), MaxVal: toPtr(200)},
// 		{Type: "area", Label: "Trên 200m²", Slug: "dien-tich-tren-200", MinVal: toPtr(200)},
// 	}

// 	if err := db.Create(&ranges).Error; err != nil {
// 		panic(err)
// 	}
// }
