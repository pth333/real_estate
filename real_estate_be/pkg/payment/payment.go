package payment

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Phương thức thanh toán được hỗ trợ
const (
	MethodVNPay   = "VNPAY"
	MethodMomo    = "MOMO"
	MethodZaloPay = "ZALOPAY"
)

// CreatePaymentRequest — dữ liệu cần để sinh link thanh toán
type CreatePaymentRequest struct {
	// Mã đơn hàng duy nhất (chính là payment_ref lưu ở bảng deposits)
	OrderRef string
	// Số tiền (VNĐ, đơn vị đồng — không phải xu)
	Amount float64
	// Nội dung thanh toán hiển thị trên cổng
	OrderInfo string
	// IP của khách tạo giao dịch (VNPay bắt buộc)
	ClientIP string
	// Phương thức khách chọn: VNPAY / MOMO / ZALOPAY
	Method string
	// Ngân hàng chỉ định (không bắt buộc)
	BankCode string
	// Thời điểm hết hạn thanh toán
	ExpiresAt time.Time
}

// CallbackResult — kết quả trả về từ cổng thanh toán (return URL / IPN)
type CallbackResult struct {
	OrderRef    string
	Amount      float64
	Success     bool
	ResponseCode string
	TransactionNo string
	Message     string
	Raw         map[string]string
}

// Gateway — cổng thanh toán. Mọi cổng (VNPay/Momo/ZaloPay/Mock) đều implement interface này.
type Gateway interface {
	// Name tên cổng đang dùng
	Name() string
	// IsMock cho biết đây là cổng giả lập (chỉ dùng khi thiếu cấu hình merchant)
	IsMock() bool
	// CreatePaymentURL sinh URL để redirect khách sang trang thanh toán
	CreatePaymentURL(req CreatePaymentRequest) (string, error)
	// VerifyCallback xác thực chữ ký và đọc kết quả giao dịch từ query params
	VerifyCallback(params map[string]string) (*CallbackResult, error)
}

// Config — cấu hình cổng thanh toán (đọc từ file config / biến môi trường)
type Config struct {
	// BaseURL của frontend, dùng để ghép tham số trả về cho cổng mock
	ReturnURL string
	// URL backend nhận IPN (server-to-server)
	IPNURL string
	VNPay  VNPayConfig
}

type VNPayConfig struct {
	TmnCode    string
	HashSecret string
	PaymentURL string
	// Locale: vn / en
	Locale string
	// Thời gian sống của link thanh toán (phút)
	ExpireMinutes int
}

// NewGateway chọn cổng thanh toán theo cấu hình.
// Chưa cấu hình merchant (tmn_code/hash_secret) → dùng cổng Mock để chạy được toàn bộ luồng.
func NewGateway(cfg Config) Gateway {
	if cfg.VNPay.TmnCode == "" || cfg.VNPay.HashSecret == "" {
		return NewMockGateway(cfg)
	}
	if cfg.VNPay.PaymentURL == "" {
		cfg.VNPay.PaymentURL = "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"
	}
	if cfg.VNPay.Locale == "" {
		cfg.VNPay.Locale = "vn"
	}
	if cfg.VNPay.ExpireMinutes <= 0 {
		cfg.VNPay.ExpireMinutes = 15
	}
	return NewVNPayGateway(cfg)
}

// formatAmount VNPay yêu cầu số tiền là số nguyên, đơn vị đồng (nhân 100).
func formatAmount(amount float64) string {
	return strconv.FormatInt(int64(amount*100), 10)
}

// buildQuery ghép query string từ map đã sắp xếp khoá tăng dần.
func buildQuery(params map[string]string, keys []string) string {
	values := url.Values{}
	for _, k := range keys {
		if v, ok := params[k]; ok && v != "" {
			values.Set(k, v)
		}
	}
	// url.Values.Encode() đã sắp xếp theo key và escape đúng chuẩn VNPay
	return values.Encode()
}

func parseAmount(raw string) float64 {
	if raw == "" {
		return 0
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	// VNPay trả về số tiền đã nhân 100
	return v / 100
}

func orderInfoOrDefault(info, ref string) string {
	if info == "" {
		return fmt.Sprintf("Dat coc xem nha don hang %s", ref)
	}
	return info
}
