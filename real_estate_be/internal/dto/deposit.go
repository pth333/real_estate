package dto

// ── Request ──────────────────────────────────────────────

// CreateDepositRequest — khách điền form đặt lịch xem nhà.
// KHÔNG nhận số tiền từ client: hệ thống tự tra phí môi giới theo giá BĐS
// để tránh khách trả số tiền vô nghĩa làm mất tác dụng chống bùng.
type CreateDepositRequest struct {
	RealEstateID uint64 `json:"real_estate_id"`
	ViewingDate  string `json:"viewing_date"`  // "2006-01-02"
	ViewingStart string `json:"viewing_start"` // "HH:MM"
	ViewingEnd   string `json:"viewing_end"`   // "HH:MM"
	// Liên hệ của buổi xem — khách xác nhận lại trước khi trả phí (có thể khác hồ sơ)
	ContactName   string `json:"contact_name"`
	ContactPhone  string `json:"contact_phone"`
	PaymentMethod string `json:"payment_method"` // VNPAY / MOMO / ZALOPAY
	BankCode      string `json:"bank_code"`
}

// BookingOptionsResponse — phí môi giới hệ thống áp dụng cho 1 BĐS.
// FE chỉ hiển thị, khách không sửa được.
type BookingOptionsResponse struct {
	RealEstateID    uint64  `json:"real_estate_id"`
	RealEstateTitle string  `json:"real_estate_title"`
	PriceVND        float64 `json:"price_vnd"`
	BrokerFee       float64 `json:"broker_fee"`
	// Nhãn phân khúc giá đang áp dụng, VD "Từ 3 đến 5 tỷ"
	PolicyLabel string `json:"policy_label"`
	// true khi không khớp chính sách nào và phải dùng mức mặc định trong config
	IsFallback bool `json:"is_fallback"`
}

// RejectDepositRequest — môi giới từ chối lịch kèm lý do
type RejectDepositRequest struct {
	Reason string `json:"reason"`
}

// CheckinLocation — vị trí hiện tại của người bấm check-in (GPS từ trình duyệt).
// Là tín hiệu bổ trợ cho OTP: 2 bên gần nhau ⇒ buổi xem đã diễn ra.
type CheckinLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	// Sai số GPS (mét) — 0 nghĩa là client không báo
	Accuracy float64 `json:"accuracy"`
}

// CheckinRequest — khách nhập OTP do môi giới hiển thị tại chỗ (kèm vị trí nếu có)
type CheckinRequest struct {
	OTP string `json:"otp"`
	CheckinLocation
}

// CheckinLocationRequest — 1 bên xác nhận đã tới nơi kèm vị trí
type CheckinLocationRequest struct {
	CheckinLocation
}

// ReportResultRequest — 1 bên báo cáo kết quả buổi xem.
// Giai đoạn đã check-in (báo mua/không mua) BẮT BUỘC kèm bằng chứng ảnh.
// Riêng khách khai BOUGHT còn phải chọn loại tài liệu mua bán (PurchaseProof).
type ReportResultRequest struct {
	Report       string   `json:"report"`
	EvidenceURLs []string `json:"evidence_urls"`
	// Chỉ dùng khi khách khai BOUGHT: PURCHASE_CONTRACT / DEPOSIT_SLIP / PAYMENT_SLIP
	PurchaseProof string `json:"purchase_proof"`
}

// CreateDisputeRequest — một bên mở tranh chấp
type CreateDisputeRequest struct {
	Reason       string   `json:"reason"`
	EvidenceURLs []string `json:"evidence_urls"`
}

// AddEvidenceRequest — upload thêm bằng chứng vào dispute đang mở
type AddEvidenceRequest struct {
	EvidenceURLs []string `json:"evidence_urls"`
}

// ResolveDisputeRequest — admin ra quyết định xử lý tiền
type ResolveDisputeRequest struct {
	Resolution string `json:"resolution"` // REFUND_CUSTOMER / TRANSFER_BROKER / SPLIT
	Note       string `json:"note"`
	// Chỉ dùng khi resolution = SPLIT: % phí môi giới hoàn cho khách (0-100)
	SplitCustomerPercent int `json:"split_customer_percent"`
}

// PurchaseDecisionRequest — admin duyệt hoặc từ chối tài liệu mua nhà.
// Chỉ khi duyệt thì đơn mới sang VISITED_BOUGHT và tồn kho dự án mới bị trừ 1 căn.
type PurchaseDecisionRequest struct {
	Approved bool   `json:"approved"`
	Note     string `json:"note"`
}

// RateBrokerRequest — khách đánh giá môi giới sau buổi xem
type RateBrokerRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// ── Response ─────────────────────────────────────────────

// DepositResponse — dữ liệu 1 deposit trả về cho FE (đã làm phẳng, không lộ otp_hash)
type DepositResponse struct {
	ID           uint64 `json:"id"`
	Status       string `json:"status"`
	CustomerID   uint64 `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	CustomerEmail string `json:"customer_email"`
	BrokerID     uint64 `json:"broker_id"`
	BrokerName   string `json:"broker_name"`
	BrokerPhone  string `json:"broker_phone"`
	BrokerEmail  string `json:"broker_email"`

	RealEstateID      uint64  `json:"real_estate_id"`
	RealEstateTitle   string  `json:"real_estate_title"`
	RealEstateAddress string  `json:"real_estate_address"`
	RealEstateSlug    string  `json:"real_estate_slug"`
	RealEstateThumb   string  `json:"real_estate_thumbnail"`
	ProjectID         *uint64 `json:"project_id"`

	// Amount = số tiền khách đã trả cho buổi xem, chính là PHÍ MÔI GIỚI (không còn tiền cọc)
	Amount        float64  `json:"amount"`
	BrokerFee     float64  `json:"broker_fee"`
	RefundAmount  *float64 `json:"refund_amount"`
	PenaltyAmount *float64 `json:"penalty_amount"`

	ViewingDate  string `json:"viewing_date"`
	ViewingStart string `json:"viewing_start"`
	ViewingEnd   string `json:"viewing_end"`

	PaymentMethod string `json:"payment_method"`
	PaymentRef    string `json:"payment_ref"`
	PaidAt        string `json:"paid_at"`

	BrokerCheckin      bool   `json:"broker_checkin"`
	CustomerCheckin    bool   `json:"customer_checkin"`
	// Bằng chứng vị trí lúc check-in: thời điểm 2 bên xác nhận + khoảng cách giữa 2 bên.
	// CheckinMatched = true nghĩa là 2 bên đã ở gần nhau ⇒ buổi xem chắc chắn diễn ra.
	BrokerCheckinAt       string   `json:"broker_checkin_at"`
	CustomerCheckinAt     string   `json:"customer_checkin_at"`
	BrokerCheckinAcc      *float64 `json:"broker_checkin_accuracy"`
	CustomerCheckinAcc    *float64 `json:"customer_checkin_accuracy"`
	CheckinDistanceMeters *float64 `json:"checkin_distance_meters"`
	CheckinMatched        bool     `json:"checkin_matched"`
	// Chất lượng bằng chứng vị trí của từng bên: "" (chưa thao tác) / AT_ESTATE / FAR / NO_LOCATION.
	// FE dựa vào đây để hiển thị cảnh báo THƯỜNG TRỰC (không phải toast, không mất khi mở lại đơn).
	BrokerCheckinEvidence   string `json:"broker_checkin_evidence"`
	CustomerCheckinEvidence string `json:"customer_checkin_evidence"`
	BrokerReport       string `json:"broker_report"`
	CustomerReport     string `json:"customer_report"`
	// Bằng chứng kèm báo cáo mua/không mua của từng bên
	BrokerReportEvidence   []string `json:"broker_report_evidence"`
	CustomerReportEvidence []string `json:"customer_report_evidence"`
	// Loại tài liệu khách xuất trình khi khai đã mua nhà (giữ cho dữ liệu đơn cũ)
	CustomerPurchaseProof string `json:"customer_purchase_proof"`
	// Thời điểm khách đặt cọc MUA bất động sản (rỗng = chưa đặt cọc).
	// Trong thời gian giữ phí mà có mốc này ⇒ phí môi giới được hoàn 100% cho khách.
	PurchaseDepositAt string `json:"purchase_deposit_at"`
	BrokerConfirmedAt      string   `json:"broker_confirmed_at"`
	RejectReason           string   `json:"reject_reason"`
	ReportDeadline         string   `json:"report_deadline"`

	// Cờ tiện dụng cho FE
	CanConfirm bool `json:"can_confirm"` // môi giới còn được xác nhận/từ chối
	CanCheckin bool `json:"can_checkin"` // được phép check-in OTP
	CanReport  bool `json:"can_report"`  // được phép báo cáo kết quả
	// Admin còn phải duyệt tài liệu mua nhà của đơn này
	CanApprovePurchase bool `json:"can_approve_purchase"`
	HasDispute         bool `json:"has_dispute"`
	HasRating          bool `json:"has_rating"`

	// Kết quả admin duyệt tài liệu mua nhà
	PurchaseDecisionBy   *uint64 `json:"purchase_decision_by"`
	PurchaseDecisionAt   string  `json:"purchase_decision_at"`
	PurchaseDecisionNote string  `json:"purchase_decision_note"`

	Dispute *DisputeResponse `json:"dispute,omitempty"`
	Rating  *BrokerRatingResponse `json:"rating,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// WithdrawalBreakdown — chi tiết dòng tiền của 1 deposit (hiển thị cho admin)
type TransactionResponse struct {
	ID        uint64  `json:"id"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	Note      string  `json:"note"`
	CreatedAt string  `json:"created_at"`
}

type DisputeResponse struct {
	ID               uint64   `json:"id"`
	DepositID        uint64   `json:"deposit_id"`
	RaisedBy         string   `json:"raised_by"`
	Reason           string   `json:"reason"`
	EvidenceURLs     []string `json:"evidence_urls"`
	Status           string   `json:"status"`
	Resolution       string   `json:"resolution"`
	ResolvedBy       *uint64  `json:"resolved_by"`
	ResolvedAt       string   `json:"resolved_at"`
	EvidenceDeadline string   `json:"evidence_deadline"`
	CreatedAt        string   `json:"created_at"`
	// Thông tin deposit kèm theo cho màn admin
	Deposit *DepositResponse `json:"deposit,omitempty"`
}

type BrokerRatingResponse struct {
	ID         uint64 `json:"id"`
	DepositID  uint64 `json:"deposit_id"`
	BrokerID   uint64 `json:"broker_id"`
	CustomerID uint64 `json:"customer_id"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment"`
	CreatedAt  string `json:"created_at"`
}

// EscrowSummaryResponse — số liệu tổng quan quỹ escrow cho admin
type EscrowSummaryResponse struct {
	// Tổng tiền đang giữ (phí môi giới đã thanh toán, chưa release)
	HoldingAmount float64 `json:"holding_amount"`
	// Tổng đã hoàn cho khách
	RefundedAmount float64 `json:"refunded_amount"`
	// Tổng đã chuyển cho môi giới
	TransferredAmount float64 `json:"transferred_amount"`
	// Tổng tiền phạt thu từ bảo lãnh môi giới
	PenaltyAmount float64 `json:"penalty_amount"`
	// Số deposit theo từng trạng thái
	StatusCounts map[string]int64 `json:"status_counts"`
	// Số dispute đang mở
	OpenDisputes int64 `json:"open_disputes"`
}

// PaymentCallbackResponse — kết quả xác nhận thanh toán trả cho FE
type PaymentCallbackResponse struct {
	DepositID uint64 `json:"deposit_id"`
	Status    string `json:"status"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

// CreateDepositResponse — kết quả tạo đặt lịch: kèm URL để redirect khách đi thanh toán
type CreateDepositResponse struct {
	Deposit    DepositResponse `json:"deposit"`
	PaymentURL string          `json:"payment_url"`
	Gateway    string          `json:"gateway"`
	IsMock     bool            `json:"is_mock"`
}

// CheckinOTPResponse — OTP check-in hiển thị cho môi giới tại chỗ
type CheckinOTPResponse struct {
	OTP       string `json:"otp"`
	ExpiresAt string `json:"expires_at"`
	// Cảnh báo khi vị trí gửi lên chưa đủ tin cậy (ngoài bán kính BĐS / sai số GPS lớn).
	// Rỗng nghĩa là vị trí đã được ghi nhận làm bằng chứng.
	LocationWarning string `json:"location_warning"`
}
