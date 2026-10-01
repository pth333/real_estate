package main

import (
	"real_estate_be/internal/initialize"
	"real_estate_be/pkg/vntime"
)

func main() {
	// Ép cả tiến trình dùng giờ VN (UTC+7) trước khi đọc config / mở DB,
	// tránh VPS chạy UTC làm lệch ngày giờ nghiệp vụ 7 tiếng.
	vntime.SetDefault()

	initialize.Run()
}
