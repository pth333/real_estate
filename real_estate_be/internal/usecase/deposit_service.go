package usecase

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"
	"time"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	"real_estate_be/pkg/mailer"
	"real_estate_be/pkg/payment"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	dateLayout = "2006-01-02"
	timeLayout = "15:04"

	// Cho phép môi giới sinh OTP sớm trước giờ hẹn
	otpEarlyMinutes = 30
)

// IDepositService — toàn bộ nghiệp vụ đặt cọc escrow (plan mục 2 → 8)
type IDepositService interface {
	// ── Khách hàng ──
	GetBookingOptions(realEstateID uint64) (*dto.BookingOptionsResponse, error)
	CreateDeposit(customerID uint64, req dto.CreateDepositRequest, clientIP string) (*dto.CreateDepositResponse, error)
	HandlePaymentCallback(params map[string]string) (*dto.PaymentCallbackResponse, error)
	ListCustomerDeposits(customerID uint64, status string, page, size int) ([]dto.DepositResponse, int64, error)
	CustomerCheckin(depositID, customerID uint64, otp string) (*dto.DepositResponse, error)
	RateBroker(depositID, customerID uint64, req dto.RateBrokerRequest) error

	// ── Môi giới ──
	ListBrokerDeposits(brokerID uint64, status string, page, size int) ([]dto.DepositResponse, int64, error)
	ConfirmDeposit(depositID, brokerID uint64) (*dto.DepositResponse, error)
	RejectDeposit(depositID, brokerID uint64, reason string) (*dto.DepositResponse, error)
	GenerateCheckinOTP(depositID, brokerID uint64) (*dto.CheckinOTPResponse, error)

	// ── Dùng chung 2 bên ──
	GetDepositDetail(depositID, requesterID uint64) (*dto.DepositResponse, error)
	GetDepositForAdmin(depositID uint64) (*dto.DepositResponse, error)
	SubmitReport(depositID, userID uint64, req dto.ReportResultRequest) (*dto.DepositResponse, error)

	// ── Tranh chấp ──
	OpenDispute(depositID, userID uint64, req dto.CreateDisputeRequest) (*dto.DisputeResponse, error)
	AddDisputeEvidence(disputeID, userID uint64, req dto.AddEvidenceRequest) (*dto.DisputeResponse, error)
	ListDisputes(status string, page, size int) ([]dto.DisputeResponse, int64, error)
	GetDisputeDetail(disputeID uint64) (*dto.DisputeResponse, error)
	ResolveDispute(disputeID, adminID uint64, req dto.ResolveDisputeRequest) (*dto.DisputeResponse, error)

	// ── Admin ──
	ListAllDeposits(status string, page, size int) ([]dto.DepositResponse, int64, error)
	DecidePurchase(depositID, adminID uint64, req dto.PurchaseDecisionRequest) (*dto.DepositResponse, error)
	GetEscrowSummary() (*dto.EscrowSummaryResponse, error)

	// ── Cron ──
	RunScheduledTasks()
	PaymentGatewayName() string
}

type depositService struct {
	depositRepo      repo.IDepositRepository
	disputeRepo      repo.IDisputeRepository
	transactionRepo  repo.ITransactionRepository
	ratingRepo       repo.IBrokerRatingRepository
	notifyLogRepo    repo.INotificationLogRepository
	policyRepo       repo.IDepositPolicyRepository
	realEstateRepo   repo.RealEstateRepository
	userRepo         repo.IUserRepository
	gateway          payment.Gateway
	mailer           mailer.Mailer
	cfg              global.DepositConfig
}

func NewDepositService(
	depositRepo repo.IDepositRepository,
	disputeRepo repo.IDisputeRepository,
	transactionRepo repo.ITransactionRepository,
	ratingRepo repo.IBrokerRatingRepository,
	notifyLogRepo repo.INotificationLogRepository,
	policyRepo repo.IDepositPolicyRepository,
	realEstateRepo repo.RealEstateRepository,
	userRepo repo.IUserRepository,
	gateway payment.Gateway,
	mailService mailer.Mailer,
) IDepositService {
	return &depositService{
		depositRepo:     depositRepo,
		disputeRepo:     disputeRepo,
		transactionRepo: transactionRepo,
		ratingRepo:      ratingRepo,
		notifyLogRepo:   notifyLogRepo,
		policyRepo:      policyRepo,
		realEstateRepo:  realEstateRepo,
		userRepo:        userRepo,
		gateway:         gateway,
		mailer:          mailService,
		cfg:             normalizeDepositConfig(global.Config.Deposit),
	}
}

func (s *depositService) PaymentGatewayName() string { return s.gateway.Name() }

// normalizeDepositConfig điền giá trị mặc định theo plan khi file cấu hình thiếu.
func normalizeDepositConfig(cfg global.DepositConfig) global.DepositConfig {
	if cfg.DefaultAmount <= 0 {
		cfg.DefaultAmount = 5_000_000
	}
	if cfg.DefaultBrokerFee < 0 {
		cfg.DefaultBrokerFee = 0
	}
	if cfg.BrokerConfirmHours <= 0 {
		cfg.BrokerConfirmHours = 24
	}
	if cfg.OTPValidMinutes <= 0 {
		cfg.OTPValidMinutes = 10
	}
	if cfg.CheckinGraceHours <= 0 {
		cfg.CheckinGraceHours = 2
	}
	if cfg.ReportWindowHours <= 0 {
		cfg.ReportWindowHours = 24
	}
	if cfg.DisputeEvidenceHours <= 0 {
		cfg.DisputeEvidenceHours = 48
	}
	if cfg.PaymentTimeoutMinutes <= 0 {
		cfg.PaymentTimeoutMinutes = 30
	}
	return cfg
}

// ══════════════════════════════════════════════════════════
// 2.1 Luồng đặt cọc
// ══════════════════════════════════════════════════════════

// GetBookingOptions trả mức cọc + phí môi giới hệ thống đề xuất cho 1 BĐS.
// FE gọi API này để hiển thị, khách không được sửa số tiền.
func (s *depositService) GetBookingOptions(realEstateID uint64) (*dto.BookingOptionsResponse, error) {
	estate, err := s.realEstateRepo.GetModelByID(realEstateID)
	if err != nil {
		return nil, errors.New("không tìm thấy bất động sản")
	}
	return s.resolveBookingOptions(estate)
}

// ensureProjectHasStock chặn đặt cọc khi dự án đã bán hết căn.
// BĐS không thuộc dự án nào (project_id NULL) thì không giới hạn.
func (s *depositService) ensureProjectHasStock(projectID *uint64) error {
	if projectID == nil {
		return nil
	}

	project, err := s.realEstateRepo.GetProjectByID(*projectID)
	if err != nil {
		return nil // không tra được dự án thì bỏ qua, không chặn oan khách
	}
	if project.TotalUnits != nil && project.SoldUnits >= *project.TotalUnits {
		return errors.New("dự án này đã hết căn, vui lòng chọn bất động sản khác")
	}
	return nil
}

// resolveBookingOptions tra bảng deposit_policies theo giá BĐS.// BĐS chưa có giá (price <= 0) → dùng mức mặc định trong config, không tra bảng
// (nếu tra sẽ khớp nhầm khoảng "Dưới 1 tỷ" và báo sai phân khúc cho khách).
func (s *depositService) resolveBookingOptions(estate *model.RealEstate) (*dto.BookingOptionsResponse, error) {
	response := &dto.BookingOptionsResponse{
		RealEstateID:    estate.ID,
		RealEstateTitle: estate.Title,
		PriceVND:        estate.PriceVND,
	}

	if estate.PriceVND > 0 {
		if policy, err := s.policyRepo.GetByPrice(estate.PriceVND); err == nil {
			response.Amount = policy.DepositAmount
			response.BrokerFee = policy.BrokerFee
			response.PolicyLabel = policy.Label
			return response, nil
		}
	}

	// Fallback: BĐS không có giá hoặc chưa cấu hình chính sách cho phân khúc đó
	response.Amount = s.cfg.DefaultAmount
	response.BrokerFee = s.cfg.DefaultBrokerFee
	response.PolicyLabel = "Mức cọc mặc định"
	response.IsFallback = true

	if response.Amount <= 0 || response.BrokerFee >= response.Amount {
		return nil, errors.New("chưa cấu hình mức cọc hợp lệ cho bất động sản này")
	}
	return response, nil
}

// CreateDeposit — khách điền form đặt cọc, hệ thống giữ chỗ khung giờ và trả URL thanh toán.
func (s *depositService) CreateDeposit(customerID uint64, req dto.CreateDepositRequest, clientIP string) (*dto.CreateDepositResponse, error) {
	estate, err := s.realEstateRepo.GetModelByID(req.RealEstateID)
	if err != nil {
		return nil, errors.New("không tìm thấy bất động sản")
	}
	if estate.UserID == nil {
		return nil, errors.New("bất động sản chưa có môi giới phụ trách")
	}

	viewingDate, start, end, err := parseViewingSlot(req.ViewingDate, req.ViewingStart, req.ViewingEnd)
	if err != nil {
		return nil, err
	}

	// Số tiền cọc + phí môi giới do hệ thống tra theo giá BĐS, KHÔNG nhận từ client
	options, err := s.resolveBookingOptions(estate)
	if err != nil {
		return nil, err
	}
	amount := options.Amount
	brokerFee := options.BrokerFee

	method := strings.ToUpper(strings.TrimSpace(req.PaymentMethod))
	switch method {
	case payment.MethodVNPay, payment.MethodMomo, payment.MethodZaloPay:
	default:
		return nil, errors.New("phương thức thanh toán không hợp lệ")
	}

	// Validate trùng khung giờ của chính BĐS đó (plan mục 5)
	overlap, err := s.depositRepo.CountOverlapSlot(estate.ID, viewingDate, start, end)
	if err != nil {
		return nil, err
	}
	if overlap > 0 {
		return nil, errors.New("khung giờ này đã có người đặt, vui lòng chọn khung giờ khác")
	}

	// BĐS thuộc dự án đã hết căn thì không nhận đặt cọc nữa
	if err := s.ensureProjectHasStock(estate.ProjectID); err != nil {
		return nil, err
	}

	deposit := &model.Deposit{
		CustomerID:    customerID,
		RealEstateID:  estate.ID,
		BrokerID:      *estate.UserID,
		ProjectID:     estate.ProjectID,
		Amount:        amount,
		BrokerFee:     brokerFee,
		ViewingDate:   viewingDate,
		ViewingStart:  start,
		ViewingEnd:    end,
		Status:        model.DepositStatusAwaitingPayment,
		PaymentMethod: method,
		PaymentRef:    buildPaymentRef(),
	}
	if err := s.depositRepo.Create(deposit); err != nil {
		return nil, err
	}

	// Sinh URL thanh toán; tiền sẽ vào tài khoản platform (escrow), KHÔNG vào môi giới
	paymentURL, err := s.gateway.CreatePaymentURL(payment.CreatePaymentRequest{
		OrderRef:  deposit.PaymentRef,
		Amount:    deposit.Amount,
		OrderInfo: fmt.Sprintf("Dat coc xem nha don %d", deposit.ID),
		ClientIP:  clientIP,
		Method:    method,
		BankCode:  req.BankCode,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.PaymentTimeoutMinutes) * time.Minute),
	})
	if err != nil {
		return nil, err
	}

	detail, err := s.GetDepositDetail(deposit.ID, customerID)
	if err != nil {
		return nil, err
	}

	return &dto.CreateDepositResponse{
		Deposit:    *detail,
		PaymentURL: paymentURL,
		Gateway:    s.gateway.Name(),
		IsMock:     s.gateway.IsMock(),
	}, nil
}

// HandlePaymentCallback — webhook/IPN xác nhận thanh toán → deposit chuyển sang PENDING.
func (s *depositService) HandlePaymentCallback(params map[string]string) (*dto.PaymentCallbackResponse, error) {
	result, err := s.gateway.VerifyCallback(params)
	if err != nil {
		return nil, err
	}

	deposit, err := s.depositRepo.GetByPaymentRef(result.OrderRef)
	if err != nil {
		return nil, errors.New("không tìm thấy đơn đặt cọc theo mã giao dịch")
	}

	// Đã xử lý trước đó (cổng gọi cả return URL và IPN) → trả lại kết quả, không ghi trùng
	if deposit.PaidAt != nil {
		return &dto.PaymentCallbackResponse{
			DepositID: deposit.ID,
			Status:    deposit.Status,
			Success:   true,
			Message:   "Giao dịch đã được xác nhận trước đó",
		}, nil
	}

	// Đơn đã bị huỷ do quá hạn thanh toán → không nhận tiền nữa
	if deposit.Status == model.DepositStatusCancelled {
		return &dto.PaymentCallbackResponse{
			DepositID: deposit.ID,
			Status:    deposit.Status,
			Success:   false,
			Message:   "Đơn đặt cọc đã quá hạn thanh toán và bị huỷ",
		}, nil
	}

	if !result.Success {
		return &dto.PaymentCallbackResponse{
			DepositID: deposit.ID,
			Status:    deposit.Status,
			Success:   false,
			Message:   "Thanh toán không thành công hoặc đã bị huỷ",
		}, nil
	}

	now := time.Now()
	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"status":  model.DepositStatusPending,
		"paid_at": now,
	}); err != nil {
		return nil, err
	}

	// Ghi nhận dòng tiền vào escrow
	if err := s.transactionRepo.Create(&model.Transaction{
		DepositID: deposit.ID,
		Type:      model.TransactionDeposit,
		Amount:    deposit.Amount,
		Note:      fmt.Sprintf("Khách thanh toán qua %s, mã giao dịch %s", deposit.PaymentMethod, result.TransactionNo),
		CreatedAt: now,
	}); err != nil {
		return nil, err
	}

	deposit.Status = model.DepositStatusPending
	deposit.PaidAt = &now

	// Thông báo môi giới: có 24h để xác nhận lịch
	s.notifyBoth(deposit, "deposit_paid",
		fmt.Sprintf("Khách đã đặt cọc xem nhà #%d", deposit.ID),
		fmt.Sprintf("Bạn đã đặt cọc thành công %s cho lịch xem nhà ngày %s (%s - %s). Tiền đang được platform giữ, sẽ hoàn nếu môi giới từ chối.",
			formatMoney(deposit.Amount), formatDate(deposit.ViewingDate), deposit.ViewingStart, deposit.ViewingEnd),
		fmt.Sprintf("Khách đã đặt cọc %s cho BĐS #%d. Vui lòng xác nhận hoặc từ chối trong %d giờ.",
			formatMoney(deposit.Amount), deposit.RealEstateID, s.cfg.BrokerConfirmHours),
	)

	return &dto.PaymentCallbackResponse{
		DepositID: deposit.ID,
		Status:    deposit.Status,
		Success:   true,
		Message:   "Thanh toán thành công, tiền đã được giữ tại escrow của platform",
	}, nil
}

// ══════════════════════════════════════════════════════════
// 2.2 Môi giới xác nhận / từ chối
// ══════════════════════════════════════════════════════════

func (s *depositService) ConfirmDeposit(depositID, brokerID uint64) (*dto.DepositResponse, error) {
	deposit, err := s.requireSide(depositID, brokerID, model.RoleBroker)
	if err != nil {
		return nil, err
	}
	if deposit.Status != model.DepositStatusPending {
		return nil, errors.New("chỉ xác nhận được đơn đang chờ xử lý")
	}

	now := time.Now()
	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"status":              model.DepositStatusBrokerConfirmed,
		"broker_confirmed_at": now,
	}); err != nil {
		return nil, err
	}

	s.notifyBoth(deposit, "broker_confirmed",
		fmt.Sprintf("Môi giới đã xác nhận lịch xem nhà #%d", deposit.ID),
		fmt.Sprintf("Lịch xem nhà đã được xác nhận.\nĐịa chỉ: %s\nThời gian: %s (%s - %s)\nMôi giới: %s - %s",
			estateAddress(deposit), formatDate(deposit.ViewingDate), deposit.ViewingStart, deposit.ViewingEnd,
			deposit.Broker.Name, deposit.Broker.Phone),
		fmt.Sprintf("Bạn đã xác nhận lịch xem nhà #%d. Hệ thống sẽ nhắc trước 24h.", deposit.ID),
	)

	return s.GetDepositDetail(deposit.ID, brokerID)
}

func (s *depositService) RejectDeposit(depositID, brokerID uint64, reason string) (*dto.DepositResponse, error) {
	deposit, err := s.requireSide(depositID, brokerID, model.RoleBroker)
	if err != nil {
		return nil, err
	}
	if deposit.Status != model.DepositStatusPending {
		return nil, errors.New("chỉ từ chối được đơn đang chờ xử lý")
	}
	if strings.TrimSpace(reason) == "" {
		return nil, errors.New("vui lòng nhập lý do từ chối")
	}

	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"reject_reason": reason,
	}); err != nil {
		return nil, err
	}
	deposit.RejectReason = reason

	// Từ chối → hoàn 100% tiền cọc cho khách
	if err := s.settle(deposit, settlementPlan{
		Status: model.DepositStatusBrokerRejected,
		Refund: deposit.Amount,
		Note:   "Môi giới từ chối lịch: " + reason,
	}); err != nil {
		return nil, err
	}

	s.notifyBoth(deposit, "broker_rejected",
		fmt.Sprintf("Môi giới từ chối lịch xem nhà #%d", deposit.ID),
		fmt.Sprintf("Môi giới đã từ chối lịch xem nhà. Lý do: %s.\nSố tiền %s sẽ được hoàn về tài khoản thanh toán của bạn.",
			reason, formatMoney(deposit.Amount)),
		fmt.Sprintf("Bạn đã từ chối đơn #%d. Tiền cọc đã được hoàn cho khách.", deposit.ID),
	)

	return s.GetDepositDetail(deposit.ID, brokerID)
}

// ══════════════════════════════════════════════════════════
// 2.3 OTP check-in chống gian lận
// ══════════════════════════════════════════════════════════

// GenerateCheckinOTP — môi giới mở app sinh OTP 6 số (hiệu lực 10 phút, dùng 1 lần).
func (s *depositService) GenerateCheckinOTP(depositID, brokerID uint64) (*dto.CheckinOTPResponse, error) {
	deposit, err := s.requireSide(depositID, brokerID, model.RoleBroker)
	if err != nil {
		return nil, err
	}
	if deposit.Status != model.DepositStatusBrokerConfirmed {
		return nil, errors.New("chỉ sinh được OTP khi lịch đã được xác nhận")
	}
	if !s.withinCheckinWindow(deposit, time.Now()) {
		return nil, errors.New("chưa tới thời gian check-in của buổi xem nhà")
	}

	otp := buildNumericOTP()
	hash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(time.Duration(s.cfg.OTPValidMinutes) * time.Minute)

	brokerChecked := true
	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"otp_hash":       string(hash),
		"otp_expires_at": expiresAt,
		"broker_checkin": brokerChecked,
	}); err != nil {
		return nil, err
	}

	return &dto.CheckinOTPResponse{
		OTP:       otp,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// CustomerCheckin — khách nhập OTP môi giới hiển thị tại chỗ → CHECKED_IN.
func (s *depositService) CustomerCheckin(depositID, customerID uint64, otp string) (*dto.DepositResponse, error) {
	deposit, err := s.requireSide(depositID, customerID, model.RoleCustomer)
	if err != nil {
		return nil, err
	}
	if deposit.Status != model.DepositStatusBrokerConfirmed {
		return nil, errors.New("đơn không ở trạng thái chờ check-in")
	}
	if deposit.OTPHash == "" || deposit.OTPExpiresAt == nil {
		return nil, errors.New("môi giới chưa sinh mã OTP")
	}
	if time.Now().After(*deposit.OTPExpiresAt) {
		return nil, errors.New("mã OTP đã hết hạn, đề nghị môi giới sinh lại")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(deposit.OTPHash), []byte(strings.TrimSpace(otp))); err != nil {
		return nil, errors.New("mã OTP không đúng")
	}

	customerChecked := true
	deadline := time.Now().Add(time.Duration(s.cfg.ReportWindowHours) * time.Hour)
	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"status":           model.DepositStatusCheckedIn,
		"customer_checkin": customerChecked,
		"otp_hash":         "",
		"otp_expires_at":   nil,
		"report_deadline":  deadline,
	}); err != nil {
		return nil, err
	}

	deposit.Status = model.DepositStatusCheckedIn
	s.notifyBoth(deposit, "checked_in",
		fmt.Sprintf("Check-in xem nhà #%d thành công", deposit.ID),
		fmt.Sprintf("Xác nhận bạn đã gặp môi giới tại buổi xem nhà. Vui lòng báo cáo kết quả (có mua / không mua) trong %d giờ tới.", s.cfg.ReportWindowHours),
		fmt.Sprintf("Khách đã check-in tại buổi xem nhà #%d. Vui lòng báo cáo kết quả trong %d giờ tới.", deposit.ID, s.cfg.ReportWindowHours),
	)

	return s.GetDepositDetail(deposit.ID, customerID)
}

// ══════════════════════════════════════════════════════════
// 2.4 Xử lý tiền theo kết quả
// ══════════════════════════════════════════════════════════

// SubmitReport — 1 bên báo cáo kết quả. Khi 2 bên đồng thuận → tất toán tiền,
// mâu thuẫn → mở tranh chấp cho admin (plan mục 2.5).
//
// Mỗi bên chỉ báo cáo được MỘT LẦN (không cho sửa) để bên khai sau không thể
// xem câu trả lời của bên kia rồi đổi cho khớp.
//
// Nguyên tắc "ai claim quyền lợi thì phải chứng minh":
// - Báo mua/không mua (đã check-in) bắt buộc kèm ảnh bằng chứng.
// - Riêng khách khai BOUGHT phải kèm tài liệu mua bán thật, vì khách là bên
//   duy nhất hưởng lợi từ việc khai BOUGHT (hoàn 100% thay vì mất phí môi giới).
func (s *depositService) SubmitReport(depositID, userID uint64, req dto.ReportResultRequest) (*dto.DepositResponse, error) {
	// Phía báo cáo suy ra từ chính bản ghi đơn, không nhận từ client
	deposit, side, err := s.getOwnedDeposit(depositID, userID)
	if err != nil {
		return nil, err
	}

	isBroker := side == model.RoleBroker
	var allowed []string

	switch deposit.Status {
	case model.DepositStatusCheckedIn:
		// Đã check-in → báo cáo kết quả mua/không mua
		allowed = []string{model.ReportBought, model.ReportNotBuy}
	case model.DepositStatusBrokerConfirmed:
		// Không check-in được → điểm danh có mặt / không đến
		allowed = []string{model.ReportAttended, model.ReportNoShow}
	default:
		return nil, errors.New("đơn đặt cọc không ở trạng thái có thể báo cáo")
	}

	// Khoá báo cáo: mỗi bên chỉ gửi 1 lần
	if isBroker && deposit.BrokerReport != "" {
		return nil, errors.New("bạn đã báo cáo kết quả rồi, không thể sửa lại")
	}
	if !isBroker && deposit.CustomerReport != "" {
		return nil, errors.New("bạn đã báo cáo kết quả rồi, không thể sửa lại")
	}

	report := strings.ToUpper(strings.TrimSpace(req.Report))
	if !containsString(allowed, report) {
		return nil, fmt.Errorf("giá trị báo cáo không hợp lệ, chỉ nhận: %s", strings.Join(allowed, ", "))
	}

	evidenceURLs := req.EvidenceURLs
	purchaseProof := strings.ToUpper(strings.TrimSpace(req.PurchaseProof))

	// Báo cáo mua/không mua phải có bằng chứng, vì tiền chia theo đúng câu trả lời.
	if deposit.Status == model.DepositStatusCheckedIn && len(evidenceURLs) == 0 {
		return nil, errors.New("vui lòng gửi kèm ít nhất 1 ảnh bằng chứng cho kết quả mua/không mua")
	}

	// Khách khai đã mua → phải xuất trình tài liệu mua bán, không nhận ảnh bất kỳ
	if !isBroker && report == model.ReportBought {
		if !containsString(model.PurchaseProofTypes, purchaseProof) {
			return nil, errors.New("khai đã mua nhà phải chọn loại tài liệu chứng minh: " +
				strings.Join(model.PurchaseProofTypes, " / "))
		}
		if len(evidenceURLs) == 0 {
			return nil, errors.New("vui lòng upload tài liệu chứng minh đã mua nhà")
		}
	}

	now := time.Now()
	fields := map[string]interface{}{}
	if isBroker {
		fields["broker_report"] = report
		fields["broker_reported_at"] = now
		fields["broker_report_evidence"] = encodeEvidenceURLs(evidenceURLs)
	} else {
		fields["customer_report"] = report
		fields["customer_reported_at"] = now
		fields["customer_report_evidence"] = encodeEvidenceURLs(evidenceURLs)
		fields["customer_purchase_proof"] = purchaseProof
	}
	if err := s.depositRepo.UpdateFields(deposit.ID, fields); err != nil {
		return nil, err
	}

	if isBroker {
		deposit.BrokerReport = report
		deposit.BrokerReportEvidence = encodeEvidenceURLs(evidenceURLs)
	} else {
		deposit.CustomerReport = report
		deposit.CustomerReportEvidence = encodeEvidenceURLs(evidenceURLs)
		deposit.CustomerPurchaseProof = purchaseProof
	}

	if err := s.evaluateReports(deposit); err != nil {
		return nil, err
	}

	return s.GetDepositDetail(deposit.ID, userID)
}

// evaluateReports áp bảng quyết định của plan mục 3 khi đã có đủ báo cáo 2 bên.
func (s *depositService) evaluateReports(deposit *model.Deposit) error {
	switch deposit.Status {
	case model.DepositStatusCheckedIn:
		// Chưa đủ 2 báo cáo → chờ bên còn lại
		if deposit.BrokerReport == "" || deposit.CustomerReport == "" {
			return nil
		}
		// Mâu thuẫn về kết quả mua/không mua → chuyển admin
		if deposit.BrokerReport != deposit.CustomerReport {
			return s.openSystemDispute(deposit, "Hai bên báo cáo kết quả mua/không mua khác nhau")
		}
		if deposit.BrokerReport == model.ReportBought {
			// Khách mua nhà → phải chờ admin duyệt tài liệu mua bán trước khi
			// tất toán và trừ tồn kho dự án (tránh tài liệu giả làm mất căn).
			return s.awaitPurchaseApproval(deposit)
		}
		return s.settleWithNotify(deposit, model.DepositStatusVisitedNotBuy, "Khách đến nhưng không mua")

	case model.DepositStatusBrokerConfirmed:
		if deposit.BrokerReport == "" || deposit.CustomerReport == "" {
			return nil
		}
		brokerPresent := deposit.BrokerReport == model.ReportAttended
		customerPresent := deposit.CustomerReport == model.ReportAttended

		switch {
		case brokerPresent && customerPresent:
			// OTP không hoạt động nhưng cả 2 đều claim có mặt → admin xác minh (plan mục 2.5)
			return s.openSystemDispute(deposit, "Hai bên đều báo có mặt nhưng chưa check-in được bằng OTP")
		case brokerPresent && !customerPresent:
			// Cả 2 thống nhất khách không đến
			return s.settleWithNotify(deposit, model.DepositStatusNoShowCustomer, "Khách không đến, môi giới có mặt")
		case !brokerPresent && customerPresent:
			// Cả 2 thống nhất môi giới không đến
			return s.settleWithNotify(deposit, model.DepositStatusNoShowBroker, "Môi giới không đến, khách có mặt")
		default:
			// Cả 2 đều nói mình vắng mặt → không xác định được bên vi phạm
			return s.openSystemDispute(deposit, "Cả hai bên đều báo không có mặt tại buổi xem nhà")
		}
	}
	return nil
}

// awaitPurchaseApproval — 2 bên đã khai khách MUA NHÀ: giữ tiền ở escrow và
// chuyển admin duyệt tài liệu mua bán. Chỉ khi admin duyệt mới tất toán
// (hoàn 100%) và trừ tồn kho dự án.
func (s *depositService) awaitPurchaseApproval(deposit *model.Deposit) error {
	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"status": model.DepositStatusPendingPurchase,
	}); err != nil {
		return err
	}
	deposit.Status = model.DepositStatusPendingPurchase

	s.notifyBoth(deposit, "purchase_pending_approval",
		fmt.Sprintf("Đơn đặt cọc #%d đang chờ duyệt tài liệu mua nhà", deposit.ID),
		"Hai bên đã xác nhận bạn mua nhà. Tài liệu mua bán đang được admin kiểm tra, tiền cọc được giữ nguyên tại platform.",
		fmt.Sprintf("Hai bên đã xác nhận khách mua nhà ở đơn #%d. Tài liệu đang chờ admin duyệt trước khi tất toán.", deposit.ID),
	)
	return nil
}

// ══════════════════════════════════════════════════════════
// Admin duyệt tài liệu mua nhà
// ══════════════════════════════════════════════════════════

// DecidePurchase — admin duyệt hoặc từ chối tài liệu mua nhà.
// Duyệt: trừ 1 căn tồn kho dự án rồi tất toán hoàn 100% cho khách.
// Từ chối: chuyển sang tranh chấp để admin xử lý tiếp bằng luồng dispute.
func (s *depositService) DecidePurchase(depositID, adminID uint64, req dto.PurchaseDecisionRequest) (*dto.DepositResponse, error) {
	deposit, err := s.depositRepo.GetByID(depositID)
	if err != nil {
		return nil, errors.New("không tìm thấy đơn đặt cọc")
	}
	if deposit.Status != model.DepositStatusPendingPurchase {
		return nil, errors.New("đơn đặt cọc không ở trạng thái chờ duyệt tài liệu mua nhà")
	}

	now := time.Now()
	note := strings.TrimSpace(req.Note)

	if !req.Approved {
		if note == "" {
			return nil, errors.New("vui lòng nhập lý do từ chối tài liệu")
		}
		if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
			"purchase_decision_by":   adminID,
			"purchase_decision_at":   now,
			"purchase_decision_note": note,
		}); err != nil {
			return nil, err
		}
		// Tài liệu không hợp lệ → chuyển tranh chấp, tiền vẫn freeze
		if err := s.openSystemDispute(deposit, "Admin từ chối tài liệu mua nhà: "+note); err != nil {
			return nil, err
		}
		return s.GetDepositForAdmin(deposit.ID)
	}

	// Duyệt: trừ tồn kho dự án trước, hết căn thì không cho duyệt
	if deposit.ProjectID != nil {
		increased, err := s.realEstateRepo.IncrementProjectSoldUnits(*deposit.ProjectID)
		if err != nil {
			return nil, err
		}
		if !increased {
			return nil, errors.New("dự án đã hết căn, không thể duyệt — cần xử lý thủ công")
		}
	}

	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"purchase_decision_by":   adminID,
		"purchase_decision_at":   now,
		"purchase_decision_note": note,
	}); err != nil {
		return nil, err
	}

	plan := s.planForStatus(model.DepositStatusVisitedBought, deposit)
	plan.Note = "Admin đã duyệt tài liệu mua nhà"
	if note != "" {
		plan.Note += ": " + note
	}
	if err := s.settle(deposit, plan); err != nil {
		return nil, err
	}
	s.notifySettlement(deposit, model.DepositStatusVisitedBought)

	return s.GetDepositForAdmin(deposit.ID)
}

// ══════════════════════════════════════════════════════════
// Đánh giá môi giới
// ══════════════════════════════════════════════════════════

func (s *depositService) RateBroker(depositID, customerID uint64, req dto.RateBrokerRequest) error {
	deposit, err := s.requireSide(depositID, customerID, model.RoleCustomer)
	if err != nil {
		return err
	}
	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("điểm đánh giá phải từ 1 đến 5")
	}
	// Chỉ đánh giá được sau khi buổi xem đã có kết luận
	if !isVisitedStatus(deposit.Status) {
		return errors.New("chỉ đánh giá được sau khi buổi xem nhà kết thúc")
	}
	if _, err := s.ratingRepo.GetByDeposit(depositID); err == nil {
		return errors.New("bạn đã đánh giá buổi xem này rồi")
	}

	if err := s.ratingRepo.Create(&model.BrokerRating{
		DepositID:  deposit.ID,
		BrokerID:   deposit.BrokerID,
		CustomerID: customerID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		CreatedAt:  time.Now(),
	}); err != nil {
		return err
	}

	// Cập nhật lại điểm trung bình + tổng số đánh giá của môi giới
	avg, total, err := s.ratingRepo.RecomputedRating(deposit.BrokerID)
	if err != nil {
		return err
	}
	return s.userRepo.UpdateFields(deposit.BrokerID, map[string]interface{}{
		"rating_avg":    math.Round(avg*100) / 100,
		"total_reviews": total,
	})
}

// ══════════════════════════════════════════════════════════
// Truy vấn chi tiết / danh sách
// ══════════════════════════════════════════════════════════

func (s *depositService) ListCustomerDeposits(customerID uint64, status string, page, size int) ([]dto.DepositResponse, int64, error) {
	offset, limit := paginate(page, size)
	items, total, err := s.depositRepo.ListByCustomer(customerID, status, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return s.mapList(items), total, nil
}

func (s *depositService) ListBrokerDeposits(brokerID uint64, status string, page, size int) ([]dto.DepositResponse, int64, error) {
	offset, limit := paginate(page, size)
	items, total, err := s.depositRepo.ListByBroker(brokerID, status, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return s.mapList(items), total, nil
}

func (s *depositService) ListAllDeposits(status string, page, size int) ([]dto.DepositResponse, int64, error) {
	offset, limit := paginate(page, size)
	items, total, err := s.depositRepo.ListAll(status, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return s.mapList(items), total, nil
}

// GetDepositDetail trả chi tiết 1 deposit cho KHÁCH hoặc MÔI GIỚI của đơn.
// Xác định vai trò theo chính bản ghi đơn (so user id), KHÔNG dựa vào role của user —
// vì một user có thể giữ nhiều role (vừa là khách của đơn này, vừa là môi giới đơn khác).
func (s *depositService) GetDepositDetail(depositID, requesterID uint64) (*dto.DepositResponse, error) {
	deposit, _, err := s.getOwnedDeposit(depositID, requesterID)
	if err != nil {
		return nil, err
	}
	return s.buildDepositDetail(deposit)
}

// GetDepositForAdmin trả chi tiết đơn cho admin (bỏ qua kiểm tra sở hữu).
// Route gọi hàm này đã được chặn bằng permission admin.deposit.view.
func (s *depositService) GetDepositForAdmin(depositID uint64) (*dto.DepositResponse, error) {
	deposit, err := s.depositRepo.GetByID(depositID)
	if err != nil {
		return nil, errors.New("không tìm thấy đơn đặt cọc")
	}
	return s.buildDepositDetail(deposit)
}

// buildDepositDetail dựng response chi tiết kèm tranh chấp đang mở và đánh giá.
func (s *depositService) buildDepositDetail(deposit *model.Deposit) (*dto.DepositResponse, error) {
	resp := s.toDepositResponse(deposit)

	// Kèm tranh chấp đang mở (nếu có)
	if dispute, err := s.disputeRepo.GetActiveByDeposit(deposit.ID); err == nil {
		d := toDisputeResponse(dispute, nil)
		resp.Dispute = &d
		resp.HasDispute = true
	}
	// Kèm đánh giá (nếu có)
	if rating, err := s.ratingRepo.GetByDeposit(deposit.ID); err == nil {
		r := toRatingResponse(rating)
		resp.Rating = &r
		resp.HasRating = true
	}

	return &resp, nil
}

// getOwnedDeposit lấy deposit và xác định user đang ở phía nào của đơn.
// Trả về side = CUSTOMER hoặc BROKER. Không nhận role từ client để tránh việc
// user nhiều role bị nhận sai phía (và không thể giả mạo vai trò).
func (s *depositService) getOwnedDeposit(depositID, userID uint64) (*model.Deposit, string, error) {
	deposit, err := s.depositRepo.GetByID(depositID)
	if err != nil {
		return nil, "", errors.New("không tìm thấy đơn đặt cọc")
	}

	if deposit.CustomerID == userID {
		return deposit, model.RoleCustomer, nil
	}
	if deposit.BrokerID == userID {
		return deposit, model.RoleBroker, nil
	}
	return nil, "", errors.New("bạn không phải khách hàng hoặc môi giới của đơn đặt cọc này")
}

// requireSide kiểm tra user có đúng vai trò yêu cầu ở đơn hay không (VD môi giới để xác nhận lịch).
func (s *depositService) requireSide(depositID, userID uint64, wantSide string) (*model.Deposit, error) {
	deposit, side, err := s.getOwnedDeposit(depositID, userID)
	if err != nil {
		return nil, err
	}
	if side != wantSide {
		if wantSide == model.RoleBroker {
			return nil, errors.New("bạn không phụ trách đơn đặt cọc này")
		}
		return nil, errors.New("đơn đặt cọc không thuộc tài khoản của bạn")
	}
	return deposit, nil
}

func (s *depositService) mapList(items []model.Deposit) []dto.DepositResponse {
	result := make([]dto.DepositResponse, 0, len(items))
	for i := range items {
		result = append(result, s.toDepositResponse(&items[i]))
	}
	return result
}

// toDepositResponse chuyển model → DTO phẳng cho FE (không lộ otp_hash).
func (s *depositService) toDepositResponse(deposit *model.Deposit) dto.DepositResponse {
	resp := dto.DepositResponse{
		ID:            deposit.ID,
		Status:        deposit.Status,
		CustomerID:    deposit.CustomerID,
		BrokerID:      deposit.BrokerID,
		RealEstateID:  deposit.RealEstateID,
		ProjectID:     deposit.ProjectID,
		Amount:        deposit.Amount,
		BrokerFee:     deposit.BrokerFee,
		RefundAmount:  deposit.RefundAmount,
		PenaltyAmount: deposit.PenaltyAmount,
		ViewingDate:   formatDate(deposit.ViewingDate),
		ViewingStart:  deposit.ViewingStart,
		ViewingEnd:    deposit.ViewingEnd,
		PaymentMethod: deposit.PaymentMethod,
		PaymentRef:    deposit.PaymentRef,
		BrokerReport:  deposit.BrokerReport,
		CustomerReport: deposit.CustomerReport,
		BrokerReportEvidence:   decodeEvidenceURLs(deposit.BrokerReportEvidence),
		CustomerReportEvidence: decodeEvidenceURLs(deposit.CustomerReportEvidence),
		CustomerPurchaseProof:  deposit.CustomerPurchaseProof,
		RejectReason:  deposit.RejectReason,
		CreatedAt:     deposit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     deposit.UpdatedAt.Format(time.RFC3339),
	}

	if deposit.Customer != nil {
		resp.CustomerName = deposit.Customer.Name
		resp.CustomerPhone = deposit.Customer.Phone
		resp.CustomerEmail = deposit.Customer.Email
	}
	if deposit.Broker != nil {
		resp.BrokerName = deposit.Broker.Name
		resp.BrokerPhone = deposit.Broker.Phone
		resp.BrokerEmail = deposit.Broker.Email
	}
	if deposit.RealEstate != nil {
		resp.RealEstateTitle = deposit.RealEstate.Title
		resp.RealEstateAddress = deposit.RealEstate.Address
		resp.RealEstateSlug = deposit.RealEstate.Slug
		if len(deposit.RealEstate.Images) > 0 {
			resp.RealEstateThumb = deposit.RealEstate.Images[0].URL
		}
	}
	if deposit.PaidAt != nil {
		resp.PaidAt = deposit.PaidAt.Format(time.RFC3339)
	}
	if deposit.BrokerConfirmedAt != nil {
		resp.BrokerConfirmedAt = deposit.BrokerConfirmedAt.Format(time.RFC3339)
	}
	if deposit.ReportDeadline != nil {
		resp.ReportDeadline = deposit.ReportDeadline.Format(time.RFC3339)
	}
	if deposit.BrokerCheckin != nil {
		resp.BrokerCheckin = *deposit.BrokerCheckin
	}
	if deposit.CustomerCheckin != nil {
		resp.CustomerCheckin = *deposit.CustomerCheckin
	}

	now := time.Now()
	resp.CanConfirm = deposit.Status == model.DepositStatusPending
	resp.CanCheckin = deposit.Status == model.DepositStatusBrokerConfirmed && s.withinCheckinWindow(deposit, now)
	resp.CanReport = deposit.Status == model.DepositStatusCheckedIn || deposit.Status == model.DepositStatusBrokerConfirmed
	resp.CanApprovePurchase = deposit.Status == model.DepositStatusPendingPurchase
	if deposit.PurchaseDecisionBy != nil {
		resp.PurchaseDecisionBy = deposit.PurchaseDecisionBy
	}
	if deposit.PurchaseDecisionAt != nil {
		resp.PurchaseDecisionAt = deposit.PurchaseDecisionAt.Format(time.RFC3339)
	}
	resp.PurchaseDecisionNote = deposit.PurchaseDecisionNote

	return resp
}

// withinCheckinWindow — cho phép check-in từ 30 phút trước giờ hẹn tới hết giờ kết thúc + ân hạn 2h.
func (s *depositService) withinCheckinWindow(deposit *model.Deposit, now time.Time) bool {
	startAt, endAt, err := viewingRange(deposit)
	if err != nil {
		return false
	}
	grace := time.Duration(s.cfg.CheckinGraceHours) * time.Hour
	return now.After(startAt.Add(-otpEarlyMinutes*time.Minute)) && now.Before(endAt.Add(grace))
}

// ══════════════════════════════════════════════════════════
// Tất toán escrow (plan mục 2.4)
// ══════════════════════════════════════════════════════════

// settlementPlan — kế hoạch chia tiền khi tất toán 1 deposit.
type settlementPlan struct {
	Status   string
	Refund   float64 // hoàn cho khách
	Transfer float64 // chuyển cho môi giới
	Penalty  float64 // phạt trừ vào bảo lãnh môi giới
	Note     string
}

// planForStatus dựng kế hoạch chia tiền chuẩn theo bảng mục 2.4.
func (s *depositService) planForStatus(status string, deposit *model.Deposit) settlementPlan {
	switch status {
	case model.DepositStatusBrokerRejected:
		return settlementPlan{Status: status, Refund: deposit.Amount, Note: "Môi giới từ chối / quá hạn xác nhận → hoàn 100%"}
	case model.DepositStatusVisitedBought:
		return settlementPlan{Status: status, Refund: deposit.Amount, Note: "Khách đến và mua nhà → hoàn 100% tiền cọc"}
	case model.DepositStatusVisitedNotBuy:
		return settlementPlan{
			Status:   status,
			Refund:   deposit.Amount - deposit.BrokerFee,
			Transfer: deposit.BrokerFee,
			Note:     "Khách đến nhưng không mua → hoàn (cọc - phí), môi giới nhận phí",
		}
	case model.DepositStatusNoShowCustomer:
		return settlementPlan{Status: status, Transfer: deposit.Amount, Note: "Khách không đến → môi giới nhận toàn bộ cọc"}
	case model.DepositStatusNoShowBroker:
		return settlementPlan{
			Status:   status,
			Refund:   deposit.Amount,
			Penalty:  deposit.BrokerFee,
			Note:     "Môi giới không đến → hoàn 100% khách + phạt môi giới",
		}
	}
	return settlementPlan{Status: status, Note: "Kết thúc đơn đặt cọc"}
}

// settleWithNotify tất toán theo status chuẩn rồi gửi thông báo cho 2 bên.
func (s *depositService) settleWithNotify(deposit *model.Deposit, status, note string) error {
	plan := s.planForStatus(status, deposit)
	if note != "" {
		plan.Note = note
	}
	if err := s.settle(deposit, plan); err != nil {
		return err
	}
	s.notifySettlement(deposit, status)
	return nil
}

// settle ghi nhận toàn bộ dòng tiền + đổi trạng thái trong 1 transaction.
func (s *depositService) settle(deposit *model.Deposit, plan settlementPlan) error {
	now := time.Now()
	refund := round2(plan.Refund)
	transfer := round2(plan.Transfer)
	penalty := round2(plan.Penalty)

	// Phạt không được vượt quá tiền bảo lãnh còn lại của môi giới
	if penalty > 0 {
		broker, err := s.userRepo.FindByID(deposit.BrokerID)
		if err == nil && broker.GuaranteeDeposit < penalty {
			penalty = broker.GuaranteeDeposit
		}
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if refund > 0 {
			txType := model.TransactionRefundPartial
			if refund >= deposit.Amount {
				txType = model.TransactionRefundFull
			}
			if err := tx.Create(&model.Transaction{
				DepositID: deposit.ID, Type: txType, Amount: refund, Note: plan.Note, CreatedAt: now,
			}).Error; err != nil {
				return err
			}
		}
		if transfer > 0 {
			if err := tx.Create(&model.Transaction{
				DepositID: deposit.ID, Type: model.TransactionTransferToBroker, Amount: transfer,
				Note: plan.Note, CreatedAt: now,
			}).Error; err != nil {
				return err
			}
		}
		if penalty > 0 {
			if err := tx.Create(&model.Transaction{
				DepositID: deposit.ID, Type: model.TransactionPenaltyBroker, Amount: penalty,
				Note: "Phạt môi giới không đến, trừ vào tiền bảo lãnh", CreatedAt: now,
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", deposit.BrokerID).
				UpdateColumn("guarantee_deposit", gorm.Expr("GREATEST(guarantee_deposit - ?, 0)", penalty)).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.Deposit{}).Where("id = ?", deposit.ID).Updates(map[string]interface{}{
			"status":         plan.Status,
			"refund_amount":  refund,
			"penalty_amount": penalty,
			"updated_at":     now,
		}).Error
	})
	if err != nil {
		return err
	}

	deposit.Status = plan.Status
	deposit.RefundAmount = &refund
	deposit.PenaltyAmount = &penalty
	return nil
}

// ══════════════════════════════════════════════════════════
// Thông báo (email + log)
// ══════════════════════════════════════════════════════════

// notifyBoth gửi email cho cả khách và môi giới kèm ghi log.
func (s *depositService) notifyBoth(deposit *model.Deposit, trigger, subject, customerBody, brokerBody string) {
	customerEmail := ""
	if deposit.Customer != nil {
		customerEmail = deposit.Customer.Email
	}
	brokerEmail := ""
	if deposit.Broker != nil {
		brokerEmail = deposit.Broker.Email
	}

	s.sendMail(deposit.ID, trigger, customerEmail, subject, customerBody)
	s.sendMail(deposit.ID, trigger, brokerEmail, subject, brokerBody)
}

func (s *depositService) sendMail(depositID uint64, trigger, to, subject, body string) {
	if to == "" {
		return
	}

	logItem := &model.NotificationLog{
		DepositID: &depositID,
		Channel:   "EMAIL",
		Trigger:   trigger,
		Recipient: to,
		Subject:   subject,
		Content:   body,
		Status:    "SENT",
		CreatedAt: time.Now(),
	}

	if err := s.mailer.Send(mailer.Message{To: to, Subject: subject, Body: body}); err != nil {
		logItem.Status = "FAILED"
		logItem.Error = err.Error()
		log.Printf("⚠️ [Deposit] gửi mail thất bại cho %s: %v", to, err)
	}

	if err := s.notifyLogRepo.Create(logItem); err != nil {
		log.Printf("⚠️ [Deposit] ghi log thông báo thất bại: %v", err)
	}
}

// notifySettlement thông báo kết quả tất toán tiền cho 2 bên.
func (s *depositService) notifySettlement(deposit *model.Deposit, status string) {
	refund := 0.0
	if deposit.RefundAmount != nil {
		refund = *deposit.RefundAmount
	}
	transfer := 0.0
	switch status {
	case model.DepositStatusVisitedNotBuy:
		transfer = deposit.BrokerFee
	case model.DepositStatusNoShowCustomer:
		transfer = deposit.Amount
	}

	s.notifyBoth(deposit, "settled_"+strings.ToLower(status),
		fmt.Sprintf("Kết quả đơn đặt cọc #%d", deposit.ID),
		fmt.Sprintf("Buổi xem nhà đã kết thúc với kết quả: %s.\nSố tiền hoàn cho bạn: %s.", statusLabel(status), formatMoney(refund)),
		fmt.Sprintf("Đơn #%d kết thúc với kết quả: %s.\nSố tiền bạn nhận được: %s.", deposit.ID, statusLabel(status), formatMoney(transfer)),
	)
}

// ══════════════════════════════════════════════════════════
// Helper
// ══════════════════════════════════════════════════════════

func parseViewingSlot(dateStr, startStr, endStr string) (time.Time, string, string, error) {
	viewingDate, err := time.ParseInLocation(dateLayout, strings.TrimSpace(dateStr), time.Local)
	if err != nil {
		return time.Time{}, "", "", errors.New("ngày xem nhà không hợp lệ (định dạng YYYY-MM-DD)")
	}

	todayStr := time.Now().Format(dateLayout)
	if viewingDate.Format(dateLayout) < todayStr {
		return time.Time{}, "", "", errors.New("ngày xem nhà không được ở quá khứ")
	}

	start := strings.TrimSpace(startStr)
	end := strings.TrimSpace(endStr)
	if _, err := time.Parse(timeLayout, start); err != nil {
		return time.Time{}, "", "", errors.New("giờ bắt đầu không hợp lệ (định dạng HH:MM)")
	}
	if _, err := time.Parse(timeLayout, end); err != nil {
		return time.Time{}, "", "", errors.New("giờ kết thúc không hợp lệ (định dạng HH:MM)")
	}
	if start >= end {
		return time.Time{}, "", "", errors.New("giờ kết thúc phải sau giờ bắt đầu")
	}

	return viewingDate, start, end, nil
}

// viewingRange ghép viewing_date + giờ bắt đầu/kết thúc thành mốc thời gian đầy đủ.
func viewingRange(deposit *model.Deposit) (time.Time, time.Time, error) {
	dateStr := deposit.ViewingDate.Format(dateLayout)
	startAt, err := time.ParseInLocation(dateLayout+" "+timeLayout, dateStr+" "+deposit.ViewingStart, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endAt, err := time.ParseInLocation(dateLayout+" "+timeLayout, dateStr+" "+deposit.ViewingEnd, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return startAt, endAt, nil
}

func buildPaymentRef() string {
	raw := strings.ReplaceAll(uuid.New().String(), "-", "")
	return fmt.Sprintf("DEP%d%s", time.Now().Unix(), strings.ToUpper(raw[:8]))
}

// buildNumericOTP sinh OTP 6 số bằng nguồn ngẫu nhiên an toàn (chống đoán mã).
func buildNumericOTP() string {
	max := big.NewInt(1_000_000)
	value, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback hiếm gặp: dùng thời gian nano để không chặn luồng nghiệp vụ
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1_000_000)
	}
	return fmt.Sprintf("%06d", value.Int64())
}

func paginate(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return (page - 1) * size, size
}

func formatDate(t time.Time) string { return t.Format(dateLayout) }

func formatMoney(amount float64) string {
	return fmt.Sprintf("%.0f VNĐ", amount)
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func containsString(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func isVisitedStatus(status string) bool {
	switch status {
	case model.DepositStatusVisitedBought, model.DepositStatusVisitedNotBuy, model.DepositStatusNoShowCustomer, model.DepositStatusNoShowBroker:
		return true
	}
	return false
}

func estateAddress(deposit *model.Deposit) string {
	if deposit.RealEstate == nil {
		return ""
	}
	parts := []string{deposit.RealEstate.Address, deposit.RealEstate.District, deposit.RealEstate.City}
	result := ""
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		if result != "" {
			result += ", "
		}
		result += part
	}
	return result
}

// statusLabel — nhãn tiếng Việt cho từng trạng thái cuối.
func statusLabel(status string) string {
	switch status {
	case model.DepositStatusBrokerRejected:
		return "Môi giới từ chối (hoàn 100%)"
	case model.DepositStatusVisitedBought:
		return "Khách đến và mua nhà (hoàn 100%)"
	case model.DepositStatusVisitedNotBuy:
		return "Khách đến nhưng không mua (hoàn cọc - phí)"
	case model.DepositStatusNoShowCustomer:
		return "Khách không đến (mất cọc)"
	case model.DepositStatusNoShowBroker:
		return "Môi giới không đến (hoàn 100% + phạt)"
	case model.DepositStatusDispute:
		return "Đang tranh chấp, chờ admin xử lý"
	case model.DepositStatusPendingPurchase:
		return "Chờ admin duyệt tài liệu mua nhà"
	case model.DepositStatusRefunded:
		return "Đã hoàn tiền cho khách"
	case model.DepositStatusCompleted:
		return "Đã tất toán"
	}
	return status
}

func toRatingResponse(rating *model.BrokerRating) dto.BrokerRatingResponse {
	return dto.BrokerRatingResponse{
		ID:         rating.ID,
		DepositID:  rating.DepositID,
		BrokerID:   rating.BrokerID,
		CustomerID: rating.CustomerID,
		Rating:     rating.Rating,
		Comment:    rating.Comment,
		CreatedAt:  rating.CreatedAt.Format(time.RFC3339),
	}
}

// decodeEvidenceURLs đọc mảng URL bằng chứng lưu dạng JSON string.
func decodeEvidenceURLs(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}
	var urls []string
	if err := json.Unmarshal([]byte(raw), &urls); err != nil {
		return []string{}
	}
	return urls
}

func encodeEvidenceURLs(urls []string) string {
	if urls == nil {
		urls = []string{}
	}
	raw, err := json.Marshal(urls)
	if err != nil {
		return "[]"
	}
	return string(raw)
}
