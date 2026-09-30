package model

import "time"

// ── Trạng thái đặt lịch ──────────────────────────────────────────────
// AWAITING_PAYMENT : vừa tạo, chưa thanh toán (giữ chỗ khung giờ)
// PENDING          : đã thanh toán PHÍ MÔI GIỚI, tiền nằm ở escrow, chờ môi giới xác nhận
// BROKER_REJECTED  : môi giới từ chối / quá 24h không phản hồi → hoàn 100% phí
// BROKER_CONFIRMED : môi giới đã xác nhận lịch xem nhà
// CHECKED_IN       : 2 bên đã check-in bằng OTP tại chỗ
// PENDING_PURCHASE_APPROVAL : 2 bên khai khách ĐÃ MUA, chờ admin duyệt tài liệu mua bán
//                     (phí vẫn freeze, CHƯA trừ tồn kho dự án)
// VISITED_BOUGHT   : admin đã duyệt tài liệu → khách mua nhà → hoàn 100% phí + trừ tồn kho dự án
// VISITED_NOT_BUY  : khách đến nhưng không mua → môi giới nhận phí, khách không được hoàn
// NO_SHOW_CUSTOMER : khách không đến → môi giới nhận toàn bộ phí
// NO_SHOW_BROKER   : môi giới không đến → hoàn 100% phí + phạt môi giới
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

// Deposit — một lượt ĐẶT LỊCH XEM NHÀ (bảng `deposits`). Kèm phần giữ phí môi giới (escrow).
//
// Vòng đời tiền: khách trả PHÍ MÔI GIỚI (= Amount = BrokerFee) → platform giữ →
// trong RefundWindowDays (mặc định 4 ngày) sau buổi xem:
//   - khách ĐẶT CỌC MUA BĐS  ⇒ hoàn 100% phí cho khách (REFUNDED)
//   - không đặt cọc          ⇒ phí thuộc về môi giới (COMPLETED)
type Deposit struct {
	// Khoá chính
	ID uint64 `gorm:"primaryKey" json:"id"`

	// ── Ai liên quan ──
	// Khách đặt lịch (người trả phí) — FK tới bảng users
	CustomerID uint64 `gorm:"column:customer_id;index" json:"customer_id"`
	// Bất động sản được xem — FK tới bảng real_estates (cũng dùng để đối chiếu toạ độ BĐS khi check-in)
	RealEstateID uint64 `gorm:"column:real_estate_id;index" json:"real_estate_id"`
	// Môi giới phụ trách — FK tới bảng users, lấy từ chủ tin lúc tạo đơn
	BrokerID uint64 `gorm:"column:broker_id;index" json:"broker_id"`
	// Dự án của BĐS tại thời điểm đặt lịch — chốt sẵn để biết trừ tồn kho dự án nào.
	// Hiện chỉ dùng cho luồng "khách khai đã mua nhà" (đang tạm bỏ, chờ luồng đặt cọc mua làm sau).
	ProjectID *uint64 `gorm:"column:project_id;index" json:"project_id"`

	// ── Tiền ──
	// SỐ TIỀN THỰC THU của đơn = phí môi giới (đơn cũ có thể là tiền cọc cũ).
	// Mọi khoản tất toán (hoàn khách / chuyển môi giới / phạt) đều tính theo Amount ⇒ BẮT BUỘC gán khi tạo đơn.
	Amount float64 `gorm:"column:amount;type:decimal(15,2)" json:"amount"`
	// Mức phí môi giới theo chính sách giá BĐS (bảng deposit_policies) — dùng lúc tạo đơn
	BrokerFee float64 `gorm:"column:broker_fee;type:decimal(15,2)" json:"broker_fee"`
	// Số tiền đã hoàn cho khách — NULL = chưa tất toán
	RefundAmount *float64 `gorm:"column:refund_amount;type:decimal(15,2)" json:"refund_amount"`
	// Tiền phạt trừ vào bảo lãnh môi giới (ca môi giới không đến)
	PenaltyAmount *float64 `gorm:"column:penalty_amount;type:decimal(15,2)" json:"penalty_amount"`

	// ── Liên hệ của buổi xem (khách nhập lúc đặt lịch) ──
	// Lưu snapshot để môi giới gọi đúng số khách để lại, không phụ thuộc hồ sơ user sau này.
	// Họ tên người liên hệ khách điền trong form đặt lịch
	ContactName string `gorm:"column:contact_name;size:100" json:"contact_name"`
	// SĐT người liên hệ khách điền (đã chuẩn hoá về dạng 0xxxxxxxxx)
	ContactPhone string `gorm:"column:contact_phone;size:20" json:"contact_phone"`

	// ── Lịch xem nhà (giờ lưu dạng "HH:MM" để so sánh trùng khung giờ) ──
	// Ngày xem (date)
	ViewingDate time.Time `gorm:"column:viewing_date;type:date" json:"viewing_date"`
	// Giờ bắt đầu khung xem — khung cố định 1 tiếng, mốc giờ chẵn trong 08:00–18:00
	ViewingStart string `gorm:"column:viewing_start;size:5" json:"viewing_start"`
	// Giờ kết thúc = ViewingStart + 1 tiếng (khách không sửa được)
	ViewingEnd string `gorm:"column:viewing_end;size:5" json:"viewing_end"`

	// Trạng thái đơn — xem hằng số DepositStatus* ở đầu file
	Status string `gorm:"column:status;index;default:AWAITING_PAYMENT" json:"status"`

	// ── Thanh toán ──
	// Cổng khách chọn: VNPAY / MOMO / ZALOPAY
	PaymentMethod string `gorm:"column:payment_method;size:20" json:"payment_method"`
	// Mã giao dịch gửi sang cổng thanh toán (cổng trả lại trong callback để tìm đơn)
	PaymentRef string `gorm:"column:payment_ref;size:100;index" json:"payment_ref"`
	// Thời điểm cổng xác nhận đã thu phí — NULL = chưa thanh toán
	PaidAt *time.Time `gorm:"column:paid_at" json:"paid_at"`

	// ── OTP check-in (môi giới sinh, khách nhập) ──
	// Hash bcrypt của mã OTP 6 số (không bao giờ trả ra JSON)
	OTPHash string `gorm:"column:otp_hash" json:"-"`
	// Hạn hiệu lực OTP (OTPValidMinutes, mặc định 10 phút)
	OTPExpiresAt *time.Time `gorm:"column:otp_expires_at" json:"otp_expires_at"`

	// ── Cờ 2 bên ĐÃ THAO TÁC check-in (chỉ để hiển thị; tiền quyết định theo toạ độ bên dưới) ──
	// true khi môi giới bấm "Tôi đã tới — sinh mã OTP" trong cửa sổ check-in
	BrokerCheckin *bool `gorm:"column:broker_checkin" json:"broker_checkin"`
	// true khi khách nhập đúng OTP hoặc bấm "Tôi đã tới"
	CustomerCheckin *bool `gorm:"column:customer_checkin" json:"customer_checkin"`

	// ── Bằng chứng VỊ TRÍ khi check-in (đây là căn cứ để tất toán tiền) ──
	// Toạ độ đã làm tròn ~100m trước khi lưu (hạn chế lưu vị trí chính xác của người dùng).
	// Chống gian lận: chỉ khi toạ độ nằm trong bán kính BĐS (checkinMatchRadiusMeters = 300m)
	// thì bên đó mới được coi là ĐÃ TỚI (xem checkinEvidenceAtEstate ở usecase).
	// Vĩ độ môi giới lúc sinh OTP
	BrokerCheckinLat *float64 `gorm:"column:broker_checkin_lat" json:"broker_checkin_lat"`
	// Kinh độ môi giới lúc sinh OTP
	BrokerCheckinLng *float64 `gorm:"column:broker_checkin_lng" json:"broker_checkin_lng"`
	// Sai số GPS (mét) của môi giới — chỉ dùng để cảnh báo chất lượng vị trí
	BrokerCheckinAcc *float64 `gorm:"column:broker_checkin_accuracy" json:"broker_checkin_accuracy"`
	// Thời điểm môi giới thao tác check-in (có/không có toạ độ đều ghi)
	BrokerCheckinAt *time.Time `gorm:"column:broker_checkin_at" json:"broker_checkin_at"`
	// Vĩ độ khách lúc check-in
	CustomerCheckinLat *float64 `gorm:"column:customer_checkin_lat" json:"customer_checkin_lat"`
	// Kinh độ khách lúc check-in
	CustomerCheckinLng *float64 `gorm:"column:customer_checkin_lng" json:"customer_checkin_lng"`
	// Sai số GPS (mét) của khách
	CustomerCheckinAcc *float64 `gorm:"column:customer_checkin_accuracy" json:"customer_checkin_accuracy"`
	// Thời điểm khách check-in (nhập OTP hoặc bấm "Tôi đã tới")
	CustomerCheckinAt *time.Time `gorm:"column:customer_checkin_at" json:"customer_checkin_at"`
	// Khoảng cách giữa 2 bên lúc check-in (mét) — NULL khi thiếu vị trí 1 bên
	CheckinDistanceMeters *float64 `gorm:"column:checkin_distance_meters" json:"checkin_distance_meters"`
	// true khi 2 bên ở gần nhau (≤ 300m) ⇒ coi như đã gặp mặt, hệ thống tự chuyển CHECKED_IN (không cần OTP)
	CheckinMatched bool `gorm:"column:checkin_matched;default:0" json:"checkin_matched"`

	// ── Dấu hiệu khách đã ĐẶT CỌC MUA bất động sản ──
	// Luồng đặt cọc mua (nút "Đặt cọc" ở trang BĐS) sẽ set field này khi khách đặt cọc thành công.
	// Khác NULL + còn trong thời gian giữ phí ⇒ hoàn 100% phí môi giới cho khách.
	PurchaseDepositAt *time.Time `gorm:"column:purchase_deposit_at" json:"purchase_deposit_at"`

	// ── Báo cáo ĐIỂM DANH (chỉ dùng khi 2 bên không check-in được bằng OTP/vị trí) ──
	// Giá trị: ATTENDED (có mặt) / NO_SHOW (không đến) — xem hằng số Report*.
	// Hằng số BOUGHT/NOT_BUY chỉ còn cho dữ liệu đơn cũ (luồng mua/không mua đã bỏ).
	// Môi giới báo: ATTENDED/NO_SHOW
	BrokerReport string `gorm:"column:broker_report;size:20" json:"broker_report"`
	// Khách báo: ATTENDED/NO_SHOW
	CustomerReport string `gorm:"column:customer_report;size:20" json:"customer_report"`
	// Thời điểm môi giới báo (mỗi bên chỉ báo 1 lần, không sửa)
	BrokerReportedAt *time.Time `gorm:"column:broker_reported_at" json:"broker_reported_at"`
	// Thời điểm khách báo
	CustomerReportedAt *time.Time `gorm:"column:customer_reported_at" json:"customer_reported_at"`
	// Ảnh bằng chứng môi giới gửi kèm báo cáo (JSON mảng URL ảnh)
	BrokerReportEvidence string `gorm:"column:broker_report_evidence;type:json" json:"broker_report_evidence"`
	// Ảnh bằng chứng khách gửi kèm báo cáo (JSON mảng URL ảnh)
	CustomerReportEvidence string `gorm:"column:customer_report_evidence;type:json" json:"customer_report_evidence"`
	// Loại tài liệu khách xuất trình khi KHAI đã mua nhà (PURCHASE_CONTRACT / DEPOSIT_SLIP / PAYMENT_SLIP)
	// — chỉ còn cho đơn cũ, luồng mới không yêu cầu nữa
	CustomerPurchaseProof string `gorm:"column:customer_purchase_proof;size:30" json:"customer_purchase_proof"`

	// ── Mốc thời gian ──
	// Thời điểm môi giới bấm "Xác nhận lịch" (đơn sang BROKER_CONFIRMED)
	BrokerConfirmedAt *time.Time `gorm:"column:broker_confirmed_at" json:"broker_confirmed_at"`
	// Đã gửi mail nhắc lịch trước 24h chưa (tránh gửi trùng)
	ReminderSentAt *time.Time `gorm:"column:reminder_sent_at" json:"reminder_sent_at"`
	// Hạn xử lý của đơn, mang 2 nghĩa theo trạng thái:
	//   - BROKER_CONFIRMED (chưa check-in): hạn 2 bên báo cáo điểm danh = giờ hẹn + 2h ân hạn + 24h
	//   - CHECKED_IN: hạn GIỮ PHÍ = hết buổi xem + RefundWindowDays (4 ngày), hết hạn thì tất toán
	ReportDeadline *time.Time `gorm:"column:report_deadline" json:"report_deadline"`
	// Lý do môi giới từ chối lịch, hoặc lý do hệ thống tự huỷ đơn quá hạn thanh toán
	RejectReason string `gorm:"column:reject_reason;type:text" json:"reject_reason"`

	// ── Admin duyệt tài liệu mua nhà (chỉ còn cho đơn cũ) ──
	// Chỉ khi admin duyệt thì đơn mới sang VISITED_BOUGHT và tồn kho dự án mới bị trừ.
	// Admin nào duyệt
	PurchaseDecisionBy *uint64 `gorm:"column:purchase_decision_by" json:"purchase_decision_by"`
	// Thời điểm duyệt
	PurchaseDecisionAt *time.Time `gorm:"column:purchase_decision_at" json:"purchase_decision_at"`
	// Ghi chú khi duyệt/từ chối tài liệu
	PurchaseDecisionNote string `gorm:"column:purchase_decision_note;type:text" json:"purchase_decision_note"`

	// Thời điểm tạo đơn (đơn AWAITING_PAYMENT quá hạn tính từ đây để tự huỷ)
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	// Lần cập nhật gần nhất
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Quan hệ (không tạo FK constraint trong DB) — nạp bằng Preload khi dựng response
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
	TransactionDeposit          = "DEPOSIT"          // khách nạp phí môi giới vào escrow
	TransactionRefundFull       = "REFUND_FULL"      // hoàn 100% cho khách
	TransactionRefundPartial    = "REFUND_PARTIAL"   // hoàn một phần cho khách
	TransactionTransferToBroker = "TRANSFER_TO_BROKER" // chuyển phí môi giới cho môi giới
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

// DepositPolicy — chính sách phí môi giới theo khoảng giá BĐS.
// price_min/price_max dạng nửa khoảng [min, max): NULL = không giới hạn phía đó.
// Số tiền KHÔNG lấy từ khách nhập mà luôn tra bảng này theo giá BĐS.
// Cột deposit_amount cũ đã bỏ: khách chỉ trả phí môi giới, không còn tiền cọc.
type DepositPolicy struct {
	ID        uint64   `gorm:"primaryKey" json:"id"`
	Label     string   `gorm:"column:label;size:100" json:"label"`
	PriceMin  *float64 `gorm:"column:price_min;type:decimal(15,2)" json:"price_min"`
	PriceMax  *float64 `gorm:"column:price_max;type:decimal(15,2)" json:"price_max"`
	BrokerFee float64  `gorm:"column:broker_fee;type:decimal(15,2)" json:"broker_fee"`
	IsActive  int      `gorm:"column:is_active;default:1" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
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
