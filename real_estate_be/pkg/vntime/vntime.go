// Package vntime — mọi mốc thời gian nghiệp vụ dùng GIỜ VIỆT NAM (UTC+7),
// không phụ thuộc múi giờ của máy/VPS (container mặc định là UTC).
//

package vntime

import "time"

// Location — UTC+7 cố định (Việt Nam không có DST).
// Dùng FixedZone thay vì LoadLocation("Asia/Ho_Chi_Minh") để KHÔNG phụ thuộc tzdata:
// image alpine/distroless thường thiếu /usr/share/zoneinfo nên LoadLocation sẽ lỗi.
var Location = time.FixedZone("UTC+7", 7*60*60)

// Now — thời điểm hiện tại theo giờ VN. Dùng thay time.Now() ở mọi chỗ nghiệp vụ.
func Now() time.Time {
	return time.Now().In(Location)
}

// SetDefault — ép TOÀN BỘ tiến trình dùng UTC+7.
//
// Sau lệnh này, mọi chỗ gọi time.Now()/time.Local trong code (kể cả thư viện) đều theo giờ VN:
//   - code tự viết dùng time.Now() / time.ParseInLocation(..., time.Local)
//   - driver MySQL với DSN `loc=Local` khi đọc/ghi cột DATETIME
//   - GORM khi tự set created_at/updated_at
//
// Gọi ngay đầu main(), TRƯỚC khi mở kết nối DB và đọc config.
func SetDefault() {
	time.Local = Location
}
