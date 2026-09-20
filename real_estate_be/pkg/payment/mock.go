package payment

import (
	"fmt"
	"net/url"
)

// MockGateway — cổng giả lập, dùng khi chưa cấu hình merchant thật.
// Luôn trả về URL trỏ về trang kết quả thanh toán của FE kèm cờ mock=1,
// tại đó FE gọi API xác nhận thanh toán để mô phỏng IPN của cổng thật.
type MockGateway struct {
	cfg Config
}

func NewMockGateway(cfg Config) Gateway {
	return &MockGateway{cfg: cfg}
}

func (g *MockGateway) Name() string { return "MOCK" }

func (g *MockGateway) IsMock() bool { return true }

func (g *MockGateway) CreatePaymentURL(req CreatePaymentRequest) (string, error) {
	base := g.cfg.ReturnURL
	if base == "" {
		return "", fmt.Errorf("thiếu payment.return_url trong cấu hình")
	}

	values := url.Values{}
	values.Set("mock", "1")
	values.Set("vnp_TxnRef", req.OrderRef)
	values.Set("vnp_Amount", formatAmount(req.Amount))
	values.Set("vnp_OrderInfo", orderInfoOrDefault(req.OrderInfo, req.OrderRef))

	return fmt.Sprintf("%s?%s", base, values.Encode()), nil
}

// VerifyCallback — với cổng mock, chữ ký không tồn tại: chỉ đọc tham số.
func (g *MockGateway) VerifyCallback(params map[string]string) (*CallbackResult, error) {
	orderRef := params["vnp_TxnRef"]
	if orderRef == "" {
		return nil, fmt.Errorf("thiếu vnp_TxnRef")
	}

	success := params["vnp_ResponseCode"] == "" || params["vnp_ResponseCode"] == "00"

	return &CallbackResult{
		OrderRef:     orderRef,
		Amount:       parseAmount(params["vnp_Amount"]),
		Success:      success,
		ResponseCode: params["vnp_ResponseCode"],
		Message:      params["vnp_OrderInfo"],
		Raw:          params,
	}, nil
}
