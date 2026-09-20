package model

import "time"

// ── Trạng thái đặt cọc ──────────────────────────────────────────────
// AWAITING_PAYMENT : vừa tạo, chưa thanh toán (giữ chỗ khung giờ)
// PENDING          : đã thanh toán, tiền nằm ở escrow, chờ môi giới xác nhận
// BROKER_REJECTED  : môi giới từ chối / quá 24h không phản hồi → hoàn 100%
// BROKER_CONFIRMED : môi giới đã xác nhận lịch xem nhà
// CHECKED_IN       : 2 bên đã check-in bằng OTP tại chỗ
// PENDING_PURCHASE_APPROVAL : 2 bên khai khách ĐÃ MUA, chờ admin duyệt tài liệu mua bán
//                     (tiền vẫn freeze, CHƯA trừ tồn kho dự án)
// VISITED_BOUGHT   : admin đã duyệt tài liệu → khách mua nhà → hoàn 100% + trừ tồn kho dự án
// VISITED_NOT_BUY  : khách đến nhưng không mua → hoàn (cọc - phí môi giới)
// NO_SHOW_CUSTOMER : khách không đến → mất cọc
// NO_SHOW_BROKER   : môi giới không đến → hoàn 100% + phạt môi giới
// DISPUTE          : tranh chấp, tiền bị freeze chờ admin xử lý
// REFUNDED         : đã hoàn tiền xong (kết thúc ở nhánh hoàn tiền)
// COMPLETED        : đã tất toán xong (kết thúc ở nhánh chuyển tiền môi giới)
const (
	DepositStatusAwaitingPayment = "AWAITING_PAYMENT"
	DepositStatusPending         = "PENDING"
	DepositStatusBrokerRejected  = "BROKER_REJECTED"
	DepositStatusBrokerConfirmed = "BROKER_CONFIRMED"
	DepositStatusCheckedIn       = "CHECKED_IN"
	DepositStatusPendingPurchase = "PENDING_PURCHASE_APPROVAL"
	DepositStatusVisitedBought   = "VISITED_BOUGHT"
	DepositStatusVisitedNotBuy   = "VISITED_NOT_BUY"
	DepositStatusNoShowCustomer  = "NO_SHOW_CUSTOMER"
	DepositStatusNoShowBroker    = "NO_SHOW_BROKER"
	DepositStatusDispute         = "DISPUTE"
	DepositStatusRefunded        = "REFUNDED"
	DepositStatusCompleted       = "COMPLETED"
	// CANCELLED : quá hạn thanh toán, hệ thống tự huỷ để trả lại khung giờ (chưa phát sinh tiền)
	DepositStatusCancelled = "CANCELLED"
)

// Deposit status đã kết thúc, không giữ khung giờ xem nhà nữa
var DepositFinishedStatuses = []string{
	DepositStatusBrokerRejected,
	DepositStatusRefunded,
	DepositStatusCompleted,
	DepositStatusCancelled,
}

// ── Báo cáo của từng bên (mỗi bên báo về CHÍNH MÌNH) ────────────────
// Giai đoạn CHECKED_IN    : BOUGHT / NOT_BUY  (khách có mua nhà hay không)
// Giai đoạn chưa check-in : ATTENDED / NO_SHOW (bên gửi báo cáo có mặt hay không)
const (
	ReportBought   = "BOUGHT"
	ReportNotBuy   = "NOT_BUY"
	ReportAttended = "ATTENDED"
	ReportNoShow   = "NO_SHOW"
)

// ── Loại tài liệu khách xuất trình khi khai ĐÃ MUA NHÀ ──────────────
// Khách là bên duy nhất hưởng lợi khi khai BOUGHT (hoàn 100% thay vì mất phí
// môi giới) nên phải chứng minh bằng tài liệu mua bán thật, không phải ảnh bất kỳ.
const (
	PurchaseProofContract    = "PURCHASE_CONTRACT" // hợp đồng mua bán
	PurchaseProofDepositSlip = "DEPOSIT_SLIP"      // phiếu đặt cọc mua nhà
	PurchaseProofPaymentSlip = "PAYMENT_SLIP"      // biên nhận chuyển tiền
)

// PurchaseProofTypes danh sách loại tài liệu hợp lệ
var PurchaseProofTypes = []string{
	PurchaseProofContract,
	PurchaseProofDepositSlip,
	PurchaseProofPaymentSlip,
}

type Deposit struct {
	ID uint64 `gorm:"primaryKey" json:"id"`

	// FK tới bảng users (khách đặt cọc / môi giới phụ trách)
	CustomerID   uint64 `gorm:"column:customer_id;index" json:"customer_id"`
	RealEstateID uint64 `gorm:"column:real_estate_id;index" json:"real_estate_id"`
	BrokerID     uint64 `gorm:"column:broker_id;index" json:"broker_id"`
	// Dự án của BĐS tại thời điểm đặt cọc — chốt sẵn để biết trừ tồn kho dự án nào
	// khi admin duyệt tài liệu mua (BĐS có thể được đổi dự án sau này).
	ProjectID *uint64 `gorm:"column:project_id;index" json:"project_id"`

	// ── Tiền ──
	Amount        float64  `gorm:"column:amount;type:decimal(15,2)" json:"amount"`
	BrokerFee     float64  `gorm:"column:broker_fee;type:decimal(15,2)" json:"broker_fee"`
	RefundAmount  *float64 `gorm:"column:refund_amount;type:decimal(15,2)" json:"refund_amount"`
	PenaltyAmount *float64 `gorm:"column:penalty_amount;type:decimal(15,2)" json:"penalty_amount"`

	// ── Lịch xem nhà (giờ lưu dạng "HH:MM" để so sánh trùng khung giờ) ──
	ViewingDate  time.Time `gorm:"column:viewing_date;type:date" json:"viewing_date"`
	ViewingStart string    `gorm:"column:viewing_start;size:5" json:"viewing_start"`
	ViewingEnd   string    `gorm:"column:viewing_end;size:5" json:"viewing_end"`

	Status string `gorm:"column:status;index;default:AWAITING_PAYMENT" json:"status"`

	// ── Thanh toán ──
	PaymentMethod string     `gorm:"column:payment_method;size:20" json:"payment_method"`
	PaymentRef    string     `gorm:"column:payment_ref;size:100;index" json:"payment_ref"`
	PaidAt        *time.Time `gorm:"column:paid_at" json:"paid_at"`

	// ── OTP check-in (môi giới generate, khách nhập) ──
	OTPHash      string     `gorm:"column:otp_hash" json:"-"`
	OTPExpiresAt *time.Time `gorm:"column:otp_expires_at" json:"otp_expires_at"`

	// ── Xác nhận 2 bên ──
	BrokerCheckin   *bool `gorm:"column:broker_checkin" json:"broker_checkin"`
	CustomerCheckin *bool `gorm:"column:customer_checkin" json:"customer_checkin"`

	// ── Báo cáo kết quả / điểm danh (xem hằng số Report*) ──
	BrokerReport      string     `gorm:"column:broker_report;size:20" json:"broker_report"`
	CustomerReport    string     `gorm:"column:customer_report;size:20" json:"customer_report"`
	BrokerReportedAt  *time.Time `gorm:"column:broker_reported_at" json:"broker_reported_at"`
	CustomerReportedAt *time.Time `gorm:"column:customer_reported_at" json:"customer_reported_at"`
	// Bằng chứng kèm theo báo cáo kết quả mua/không mua (JSON mảng URL ảnh).
	// Bắt buộc phải có để làm căn cứ cho admin khi 2 bên báo cáo lệch nhau.
	BrokerReportEvidence   string `gorm:"column:broker_report_evidence;type:json" json:"broker_report_evidence"`
	CustomerReportEvidence string `gorm:"column:customer_report_evidence;type:json" json:"customer_report_evidence"`
	// Loại tài liệu khách xuất trình khi khai đã mua nhà (xem hằng số PurchaseProof*)
	CustomerPurchaseProof string `gorm:"column:customer_purchase_proof;size:30" json:"customer_purchase_proof"`

	// ── Timestamps ──
	BrokerConfirmedAt *time.Time `gorm:"column:broker_confirmed_at" json:"broker_confirmed_at"`
	ReminderSentAt    *time.Time `gorm:"column:reminder_sent_at" json:"reminder_sent_at"`
	// Hạn cuối 2 bên phải báo cáo kết quả khi không check-in được (viewing_start + 2h + 24h)
	ReportDeadline *time.Time `gorm:"column:report_deadline" json:"report_deadline"`
	RejectReason   string     `gorm:"column:reject_reason;type:text" json:"reject_reason"`

	// ── Admin duyệt tài liệu mua nhà ──
	// Chỉ khi admin duyệt thì đơn mới sang VISITED_BOUGHT và tồn kho dự án mới bị trừ.
	PurchaseDecisionBy   *uint64    `gorm:"column:purchase_decision_by" json:"purchase_decision_by"`
	PurchaseDecisionAt   *time.Time `gorm:"column:purchase_decision_at" json:"purchase_decision_at"`
	PurchaseDecisionNote string     `gorm:"column:purchase_decision_note;type:text" json:"purchase_decision_note"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Quan hệ (không tạo FK constraint trong DB)
	Customer   *User       `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Broker     *User       `gorm:"foreignKey:BrokerID;references:ID" json:"broker,omitempty"`
	RealEstate *RealEstate `gorm:"foreignKey:RealEstateID;references:ID" json:"real_estate,omitempty"`
}

func (Deposit) TableName() string { return "deposits" }

// IsFinished cho biết deposit đã tất toán, không còn giữ khung giờ.
func (d *Deposit) IsFinished() bool {
	switch d.Status {
	case DepositStatusBrokerRejected, DepositStatusRefunded, DepositStatusCompleted, DepositStatusCancelled:
		return true
	}
	return false
}

// Transaction — lịch sử mọi dòng tiền của một deposit (escrow ledger).
const (
	TransactionDeposit          = "DEPOSIT"          // khách nạp tiền cọc vào escrow
	TransactionRefundFull       = "REFUND_FULL"      // hoàn 100% cho khách
	TransactionRefundPartial    = "REFUND_PARTIAL"   // hoàn một phần cho khách
	TransactionTransferToBroker = "TRANSFER_TO_BROKER" // chuyển phí/cọc cho môi giới
	TransactionPenaltyBroker    = "PENALTY_BROKER"   // phạt trừ vào bảo lãnh môi giới
)

type Transaction struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	DepositID uint64    `gorm:"column:deposit_id;index" json:"deposit_id"`
	Type      string    `gorm:"column:type;size:30" json:"type"`
	Amount    float64   `gorm:"column:amount;type:decimal(15,2)" json:"amount"`
	Note      string    `gorm:"column:note;type:text" json:"note"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Transaction) TableName() string { return "transactions" }

// Dispute — tranh chấp cần admin phân xử. Tiền bị freeze cho tới khi có quyết định.
const (
	DisputeStatusOpen      = "OPEN"
	DisputeStatusReviewing = "REVIEWING"
	DisputeStatusResolved  = "RESOLVED"

	DisputeRaisedByCustomer = "CUSTOMER"
	DisputeRaisedByBroker   = "BROKER"
	DisputeRaisedBySystem   = "SYSTEM"

	DisputeResolutionRefundCustomer  = "REFUND_CUSTOMER"
	DisputeResolutionTransferBroker  = "TRANSFER_BROKER"
	DisputeResolutionSplit           = "SPLIT"
)

type Dispute struct {
	ID           uint64  `gorm:"primaryKey" json:"id"`
	DepositID    uint64  `gorm:"column:deposit_id;index" json:"deposit_id"`
	RaisedBy     string  `gorm:"column:raised_by;size:20" json:"raised_by"`
	Reason       string  `gorm:"column:reason;type:text" json:"reason"`
	EvidenceURLs string  `gorm:"column:evidence_urls;type:json" json:"evidence_urls"`
	Status       string  `gorm:"column:status;size:20;default:OPEN" json:"status"`
	Resolution   string  `gorm:"column:resolution;size:30" json:"resolution"`
	ResolvedBy   *uint64 `gorm:"column:resolved_by" json:"resolved_by"`
	ResolvedAt   *time.Time `gorm:"column:resolved_at" json:"resolved_at"`
	// Hạn cuối 2 bên upload bằng chứng (mở dispute + 48h)
	EvidenceDeadline *time.Time `gorm:"column:evidence_deadline" json:"evidence_deadline"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`

	Deposit *Deposit `gorm:"foreignKey:DepositID;references:ID" json:"deposit,omitempty"`
}

func (Dispute) TableName() string { return "disputes" }

// DepositPolicy — chính sách mức cọc đề xuất theo khoảng giá BĐS.
// price_min/price_max dạng nửa khoảng [min, max): NULL = không giới hạn phía đó.
// Số tiền cọc KHÔNG lấy từ khách nhập mà luôn tra bảng này theo giá BĐS.
type DepositPolicy struct {
	ID            uint64   `gorm:"primaryKey" json:"id"`
	Label         string   `gorm:"column:label;size:100" json:"label"`
	PriceMin      *float64 `gorm:"column:price_min;type:decimal(15,2)" json:"price_min"`
	PriceMax      *float64 `gorm:"column:price_max;type:decimal(15,2)" json:"price_max"`
	DepositAmount float64  `gorm:"column:deposit_amount;type:decimal(15,2)" json:"deposit_amount"`
	BrokerFee     float64  `gorm:"column:broker_fee;type:decimal(15,2)" json:"broker_fee"`
	IsActive      int      `gorm:"column:is_active;default:1" json:"is_active"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (DepositPolicy) TableName() string { return "deposit_policies" }

// BrokerRating — đánh giá môi giới sau buổi xem nhà.
type BrokerRating struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	DepositID  uint64    `gorm:"column:deposit_id;index" json:"deposit_id"`
	BrokerID   uint64    `gorm:"column:broker_id;index" json:"broker_id"`
	CustomerID uint64    `gorm:"column:customer_id;index" json:"customer_id"`
	Rating     int       `gorm:"column:rating" json:"rating"`
	Comment    string    `gorm:"column:comment;type:text" json:"comment"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (BrokerRating) TableName() string { return "broker_ratings" }

// NotificationLog — log các thông báo (email) đã gửi cho 2 bên.
// Đặt tên bảng riêng để không đụng bảng notifications (thông báo tin đăng mới qua SSE).
type NotificationLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	DepositID *uint64   `gorm:"column:deposit_id;index" json:"deposit_id"`
	Channel   string    `gorm:"column:channel;size:20" json:"channel"`
	Trigger   string    `gorm:"column:trigger;size:50" json:"trigger"`
	Recipient string    `gorm:"column:recipient;size:255" json:"recipient"`
	Subject   string    `gorm:"column:subject" json:"subject"`
	Content   string    `gorm:"column:content;type:text" json:"content"`
	Status    string    `gorm:"column:status;size:20" json:"status"`
	Error     string    `gorm:"column:error;type:text" json:"error"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (NotificationLog) TableName() string { return "notification_logs" }
