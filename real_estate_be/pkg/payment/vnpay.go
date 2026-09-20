package payment

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// VNPayGateway — adapter cổng VNPay (mặc định sandbox).
// Tham chiếu: https://sandbox.vnpayment.vn/apis/docs/thanh-toan-pay
type VNPayGateway struct {
	cfg Config
}

func NewVNPayGateway(cfg Config) Gateway {
	return &VNPayGateway{cfg: cfg}
}

func (g *VNPayGateway) Name() string { return MethodVNPay }

func (g *VNPayGateway) IsMock() bool { return false }

// CreatePaymentURL sinh URL thanh toán VNPay kèm chữ ký HMAC-SHA512.
func (g *VNPayGateway) CreatePaymentURL(req CreatePaymentRequest) (string, error) {
	now := time.Now()
	expire := now.Add(time.Duration(g.cfg.VNPay.ExpireMinutes) * time.Minute)
	if !req.ExpiresAt.IsZero() {
		expire = req.ExpiresAt
	}

	params := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    g.cfg.VNPay.TmnCode,
		"vnp_Locale":     g.cfg.VNPay.Locale,
		"vnp_CurrCode":   "VND",
		"vnp_TxnRef":     req.OrderRef,
		"vnp_OrderInfo":  orderInfoOrDefault(req.OrderInfo, req.OrderRef),
		"vnp_OrderType":  "other",
		"vnp_Amount":     formatAmount(req.Amount),
		"vnp_ReturnUrl":  g.cfg.ReturnURL,
		"vnp_IpAddr":     req.ClientIP,
		"vnp_CreateDate": now.Format("20060102150405"),
		"vnp_ExpireDate": expire.Format("20060102150405"),
	}
	if req.BankCode != "" {
		params["vnp_BankCode"] = req.BankCode
	}

	hashData := buildQuery(params, sortedKeys(params))
	secureHash := g.sign(hashData)

	return fmt.Sprintf("%s?%s&vnp_SecureHash=%s",
		g.cfg.VNPay.PaymentURL, hashData, secureHash), nil
}

// VerifyCallback xác thực chữ ký HMAC-SHA512 và đọc kết quả giao dịch.
func (g *VNPayGateway) VerifyCallback(params map[string]string) (*CallbackResult, error) {
	receivedHash := strings.ToLower(params["vnp_SecureHash"])
	if receivedHash == "" {
		return nil, fmt.Errorf("thiếu vnp_SecureHash")
	}

	// Loại bỏ 2 tham số chữ ký trước khi tính lại hash
	filtered := make(map[string]string, len(params))
	for k, v := range params {
		if k == "vnp_SecureHash" || k == "vnp_SecureHashType" {
			continue
		}
		filtered[k] = v
	}

	hashData := buildQuery(filtered, sortedKeys(filtered))
	expected := g.sign(hashData)

	if !hmac.Equal([]byte(expected), []byte(receivedHash)) {
		return nil, fmt.Errorf("chữ ký VNPay không hợp lệ")
	}

	responseCode := params["vnp_ResponseCode"]
	transactionStatus := params["vnp_TransactionStatus"]

	result := &CallbackResult{
		OrderRef:      params["vnp_TxnRef"],
		Amount:        parseAmount(params["vnp_Amount"]),
		ResponseCode:  responseCode,
		TransactionNo: params["vnp_TransactionNo"],
		Message:       params["vnp_OrderInfo"],
		Raw:           params,
	}
	// Giao dịch thành công khi cả response code và transaction status đều là "00"
	result.Success = responseCode == "00" && (transactionStatus == "" || transactionStatus == "00")

	return result, nil
}

func (g *VNPayGateway) sign(data string) string {
	mac := hmac.New(sha512.New, []byte(g.cfg.VNPay.HashSecret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func sortedKeys(params map[string]string) []string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ParseBoolInt đọc tham số kiểu số nguyên an toàn (dùng cho các cổng khác nếu cần).
func ParseBoolInt(raw string) bool {
	v, err := strconv.Atoi(raw)
	return err == nil && v == 1
}
