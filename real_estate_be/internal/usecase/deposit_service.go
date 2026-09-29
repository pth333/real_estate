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

	// Khung giờ xem nhà: mỗi khung cố định 1 tiếng, chỉ trong giờ làm việc
	viewingSlotFirstHour = 8
	viewingSlotLastHour  = 18
	viewingSlotHours     = 1

	// Bán kính coi là 2 bên đã gặp nhau / đã ở đúng bất động sản (mét).
	// Rộng rãi vì GPS trong nhà/chung cư lệch vài chục tới vài trăm mét.
	checkinMatchRadiusMeters = 300
	// Sai số GPS tối đa còn dùng được làm bằng chứng (mét)
	checkinMaxAccuracyMeters = 150
)

// IDepositService — toàn bộ nghiệp vụ đặt lịch escrow (plan mục 2 → 8)
type IDepositService interface {
	// ── Khách hàng ──
	GetBookingOptions(realEstateID uint64) (*dto.BookingOptionsResponse, error)
	CreateDeposit(customerID uint64, req dto.CreateDepositRequest, clientIP string) (*dto.CreateDepositResponse, error)
	HandlePaymentCallback(params map[string]string) (*dto.PaymentCallbackResponse, error)
	ListCustomerDeposits(customerID uint64, status string, page, size int) ([]dto.DepositResponse, int64, error)
	CustomerCheckin(depositID, customerID uint64, otp string, location dto.CheckinLocationRequest) (*dto.DepositResponse, error)
	// CustomerArrived — khách báo "tôi đã tới" bằng vị trí khi không nhập được OTP
	CustomerArrived(depositID, customerID uint64, location dto.CheckinLocationRequest) (*dto.DepositResponse, error)
	// RefundFeeOnPurchaseDeposit — luồng đặt cọc mua BĐS gọi khi khách đặt cọc thành công:
	// hoàn 100% phí môi giới của buổi xem nếu còn trong thời gian giữ phí
	RefundFeeOnPurchaseDeposit(depositID uint64) error
	RateBroker(depositID, customerID uint64, req dto.RateBrokerRequest) error

	// ── Môi giới ──
	ListBrokerDeposits(brokerID uint64, status string, page, size int) ([]dto.DepositResponse, int64, error)
	ConfirmDeposit(depositID, brokerID uint64) (*dto.DepositResponse, error)
	RejectDeposit(depositID, brokerID uint64, reason string) (*dto.DepositResponse, error)
	GenerateCheckinOTP(depositID, brokerID uint64, location dto.CheckinLocationRequest) (*dto.CheckinOTPResponse, error)

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
	if cfg.DefaultBrokerFee <= 0 {
		cfg.DefaultBrokerFee = 200_000
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
	if cfg.RefundWindowDays <= 0 {
		cfg.RefundWindowDays = 4
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
// 2.1 Luồng đặt lịch
// ══════════════════════════════════════════════════════════

// GetBookingOptions trả phí môi giới hệ thống áp dụng cho 1 BĐS.
// FE gọi API này để hiển thị, khách không được sửa số tiền.
func (s *depositService) GetBookingOptions(realEstateID uint64) (*dto.BookingOptionsResponse, error) {
	estate, err := s.realEstateRepo.GetModelByID(realEstateID)
	if err != nil {
		return nil, errors.New("không tìm thấy bất động sản")
	}
	return s.resolveBookingOptions(estate)
}

// ensureProjectHasStock chặn đặt lịch khi dự án đã bán hết căn.
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

// resolveBookingOptions tra bảng deposit_policies theo giá BĐS để lấy PHÍ MÔI GIỚI.
// BĐS chưa có giá (price <= 0) → dùng mức mặc định trong config, không tra bảng
// (nếu tra sẽ khớp nhầm khoảng "Dưới 1 tỷ" và báo sai phân khúc cho khách).
func (s *depositService) resolveBookingOptions(estate *model.RealEstate) (*dto.BookingOptionsResponse, error) {
	response := &dto.BookingOptionsResponse{
		RealEstateID:    estate.ID,
		RealEstateTitle: estate.Title,
		PriceVND:        estate.PriceVND,
	}

	if estate.PriceVND > 0 {
		if policy, err := s.policyRepo.GetByPrice(estate.PriceVND); err == nil {
			response.BrokerFee = policy.BrokerFee
			response.PolicyLabel = policy.Label
			return response, nil
		}
	}

	// Fallback: BĐS không có giá hoặc chưa cấu hình chính sách cho phân khúc đó
	response.BrokerFee = s.cfg.DefaultBrokerFee
	response.PolicyLabel = "Phí môi giới mặc định"
	response.IsFallback = true

	if response.BrokerFee <= 0 {
		return nil, errors.New("chưa cấu hình phí môi giới hợp lệ cho bất động sản này")
	}
	return response, nil
}

// CreateDeposit — khách điền form đặt lịch, hệ thống giữ chỗ khung giờ và trả URL thanh toán.
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

	contactName, contactPhone, err := normalizeContactInfo(req.ContactName, req.ContactPhone)
	if err != nil {
		return nil, err
	}

	// Phí môi giới do hệ thống tra theo giá BĐS, KHÔNG nhận từ client
	options, err := s.resolveBookingOptions(estate)
	if err != nil {
		return nil, err
	}
	// Khách chỉ trả phí môi giới cho buổi xem, không còn phí môi giới
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

	// BĐS thuộc dự án đã hết căn thì không nhận đặt lịch nữa
	if err := s.ensureProjectHasStock(estate.ProjectID); err != nil {
		return nil, err
	}

	deposit := &model.Deposit{
		CustomerID:    customerID,
		RealEstateID:  estate.ID,
		BrokerID:      *estate.UserID,
		ProjectID:     estate.ProjectID,
		Amount:        brokerFee,
		BrokerFee:     brokerFee,
		ViewingDate:   viewingDate,
		ViewingStart:  start,
		ViewingEnd:    end,
		ContactName:   contactName,
		ContactPhone:  contactPhone,
		Status:        model.DepositStatusAwaitingPayment,
		PaymentMethod: method,
		PaymentRef:    buildPaymentRef(),
	}
	if err := s.depositRepo.Create(deposit); err != nil {
		return nil, err
	}

	// Sinh URL thanh toán cho PHÍ MÔI GIỚI; tiền vào tài khoản platform, KHÔNG vào môi giới ngay
	paymentURL, err := s.gateway.CreatePaymentURL(payment.CreatePaymentRequest{
		OrderRef:  deposit.PaymentRef,
		Amount:    brokerFee,
		OrderInfo: fmt.Sprintf("Phi moi gioi xem nha don %d", deposit.ID),
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
		return nil, errors.New("không tìm thấy đơn đặt lịch theo mã giao dịch")
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
			Message:   "Đơn đặt lịch đã quá hạn thanh toán và bị huỷ",
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

	// Ghi nhận dòng tiền phí môi giới vào escrow
	if err := s.transactionRepo.Create(&model.Transaction{
		DepositID: deposit.ID,
		Type:      model.TransactionDeposit,
		Amount:    deposit.BrokerFee,
		Note:      fmt.Sprintf("Khách thanh toán phí môi giới qua %s, mã giao dịch %s", deposit.PaymentMethod, result.TransactionNo),
		CreatedAt: now,
	}); err != nil {
		return nil, err
	}

	deposit.Status = model.DepositStatusPending
	deposit.PaidAt = &now

	// Thông báo môi giới: có 24h để xác nhận lịch
	contactName, contactPhone := contactOf(deposit)
	s.notifyBoth(deposit, "deposit_paid",
		fmt.Sprintf("Khách đã đặt lịch xem nhà #%d", deposit.ID),
		fmt.Sprintf("Bạn đã thanh toán phí môi giới %s cho lịch xem nhà ngày %s (%s - %s). Tiền đang được platform giữ, sẽ hoàn nếu môi giới từ chối.",
			formatMoney(deposit.BrokerFee), formatDate(deposit.ViewingDate), deposit.ViewingStart, deposit.ViewingEnd),
		fmt.Sprintf("Khách đã trả phí môi giới %s cho BĐS #%d. Vui lòng xác nhận hoặc từ chối trong %d giờ.\nLiên hệ khách: %s - %s",
			formatMoney(deposit.BrokerFee), deposit.RealEstateID, s.cfg.BrokerConfirmHours, contactName, contactPhone),
	)

	return &dto.PaymentCallbackResponse{
		DepositID: deposit.ID,
		Status:    deposit.Status,
		Success:   true,
		Message:   "Thanh toán thành công, phí môi giới đang được platform giữ",
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

	// Từ chối → hoàn 100% phí môi giới cho khách
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
		fmt.Sprintf("Bạn đã từ chối đơn #%d. Phí môi giới đã được hoàn cho khách.", deposit.ID),
	)

	return s.GetDepositDetail(deposit.ID, brokerID)
}

// ══════════════════════════════════════════════════════════
// 2.3 OTP check-in chống gian lận
// ══════════════════════════════════════════════════════════

// GenerateCheckinOTP — môi giới xác nhận đã tới nơi kèm vị trí, hệ thống sinh OTP 6 số
// (hiệu lực 10 phút, dùng 1 lần) cho khách nhập.
//
// Vị trí là bằng chứng cho "môi giới có mặt": chỉ khi toạ độ nằm trong bán kính BĐS thì
// broker_checkin mới tính là đã xác nhận. Không có vị trí thì vẫn sinh OTP (khách vẫn
// check-in được) nhưng môi giới chưa có căn cứ đòi phí, FE nhận cảnh báo để bấm lại.
func (s *depositService) GenerateCheckinOTP(depositID, brokerID uint64, location dto.CheckinLocationRequest) (*dto.CheckinOTPResponse, error) {
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
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.cfg.OTPValidMinutes) * time.Minute)

	fields := map[string]interface{}{
		"otp_hash":       string(hash),
		"otp_expires_at": expiresAt,
	}

	warning := ""
	if usableCheckinLocation(location.CheckinLocation) && atEstate(deposit.RealEstate, location.CheckinLocation) {
		saveCheckinLocation(deposit, fields, true, location.CheckinLocation, now)
		brokerChecked := true
		fields["broker_checkin"] = brokerChecked
		deposit.BrokerCheckin = &brokerChecked
	} else {
		fields["broker_checkin"] = false
		warning = "Chưa ghi nhận được vị trí của bạn tại bất động sản. Hãy bật định vị và bấm lại để có bằng chứng bạn đã tới."
	}

	if err := s.depositRepo.UpdateFields(deposit.ID, fields); err != nil {
		return nil, err
	}
	if err := s.refreshCheckinMatch(deposit); err != nil {
		return nil, err
	}
	// 2 bên đã ở gần nhau thì coi như đã gặp mặt, không cần khách nhập OTP nữa
	if err := s.autoCheckinIfMatched(deposit); err != nil {
		return nil, err
	}

	return &dto.CheckinOTPResponse{
		OTP:             otp,
		ExpiresAt:       expiresAt.Format(time.RFC3339),
		LocationWarning: warning,
	}, nil
}

// saveCheckinLocation lưu toạ độ (đã làm tròn ~100m) của 1 bên vào đơn và map update.
// Trả về false nếu toạ độ gửi lên không dùng được làm bằng chứng.
func saveCheckinLocation(
	deposit *model.Deposit, fields map[string]interface{}, isBroker bool,
	location dto.CheckinLocation, now time.Time,
) bool {
	if !usableCheckinLocation(location) {
		return false
	}

	lat, lng := roundCoord(location.Latitude), roundCoord(location.Longitude)
	var accuracy *float64
	if location.Accuracy > 0 {
		value := location.Accuracy
		accuracy = &value
	}

	if isBroker {
		fields["broker_checkin_lat"] = lat
		fields["broker_checkin_lng"] = lng
		fields["broker_checkin_at"] = now
		if accuracy != nil {
			fields["broker_checkin_accuracy"] = *accuracy
		}
		deposit.BrokerCheckinLat, deposit.BrokerCheckinLng = &lat, &lng
		deposit.BrokerCheckinAcc, deposit.BrokerCheckinAt = accuracy, &now
		return true
	}

	fields["customer_checkin_lat"] = lat
	fields["customer_checkin_lng"] = lng
	fields["customer_checkin_at"] = now
	if accuracy != nil {
		fields["customer_checkin_accuracy"] = *accuracy
	}
	deposit.CustomerCheckinLat, deposit.CustomerCheckinLng = &lat, &lng
	deposit.CustomerCheckinAcc, deposit.CustomerCheckinAt = accuracy, &now
	return true
}

// refreshCheckinMatch cập nhật khoảng cách 2 bên + cờ "đã gặp nhau" sau mỗi lần có vị trí mới.
func (s *depositService) refreshCheckinMatch(deposit *model.Deposit) error {
	var distance *float64
	matched := false
	if deposit.BrokerCheckinLat != nil && deposit.BrokerCheckinLng != nil &&
		deposit.CustomerCheckinLat != nil && deposit.CustomerCheckinLng != nil {
		value := round2(distanceMeters(
			*deposit.BrokerCheckinLat, *deposit.BrokerCheckinLng,
			*deposit.CustomerCheckinLat, *deposit.CustomerCheckinLng,
		))
		distance = &value
		matched = value <= checkinMatchRadiusMeters
	}

	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"checkin_distance_meters": distance,
		"checkin_matched":         matched,
	}); err != nil {
		return err
	}

	deposit.CheckinDistanceMeters = distance
	deposit.CheckinMatched = matched
	return nil
}

// autoCheckinIfMatched — 2 bên đã ở gần nhau ⇒ coi như buổi xem đã diễn ra:
// chuyển CHECKED_IN luôn, không bắt khách nhập OTP nữa (OTP chỉ còn là đường dự phòng).
func (s *depositService) autoCheckinIfMatched(deposit *model.Deposit) error {
	if !deposit.CheckinMatched || deposit.Status != model.DepositStatusBrokerConfirmed {
		return nil
	}

	customerChecked := true
	// Hạn giữ phí: hết buổi xem + 4 ngày để khách kịp đặt cọc mua BĐS
	deadline := s.feeHoldDeadline(deposit)
	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"status":           model.DepositStatusCheckedIn,
		"customer_checkin": customerChecked,
		"otp_hash":         "",
		"otp_expires_at":   nil,
		"report_deadline":  deadline,
	}); err != nil {
		return err
	}

	deposit.Status = model.DepositStatusCheckedIn
	deposit.CustomerCheckin = &customerChecked
	deposit.OTPHash = ""
	deposit.OTPExpiresAt = nil
	s.notifyBoth(deposit, "checked_in_location",
		fmt.Sprintf("Đã xác nhận buổi xem nhà #%d", deposit.ID),
		fmt.Sprintf("Hệ thống ghi nhận bạn và môi giới đã ở cùng địa điểm tại buổi xem nhà. Phí môi giới được giữ trong %d ngày: nếu bạn đặt cọc mua bất động sản trong thời gian này, phí sẽ được hoàn 100%% cho bạn.", s.cfg.RefundWindowDays),
		fmt.Sprintf("Hệ thống ghi nhận bạn và khách đã ở cùng địa điểm tại buổi xem nhà #%d. Phí môi giới sẽ được chuyển cho bạn sau %d ngày nếu khách không đặt cọc mua bất động sản.", deposit.ID, s.cfg.RefundWindowDays),
	)
	return nil
}

// CustomerArrived — khách báo "tôi đã tới" bằng vị trí khi không nhập được OTP
// (mất mạng, hết pin, mã hết hạn...). Không cần OTP nhưng BẮT BUỘC có vị trí dùng được:
// - 2 bên ở gần nhau ⇒ tự động CHECKED_IN (đã gặp mặt thật).
// - Chỉ khách ở gần BĐS ⇒ ghi nhận khách có mặt, chờ đối chiếu với log của môi giới.
func (s *depositService) CustomerArrived(depositID, customerID uint64, location dto.CheckinLocationRequest) (*dto.DepositResponse, error) {
	deposit, err := s.requireSide(depositID, customerID, model.RoleCustomer)
	if err != nil {
		return nil, err
	}
	if deposit.Status != model.DepositStatusBrokerConfirmed {
		return nil, errors.New("đơn không ở trạng thái chờ check-in")
	}
	if !s.withinCheckinWindow(deposit, time.Now()) {
		return nil, errors.New("chưa tới thời gian check-in của buổi xem nhà")
	}
	if !usableCheckinLocation(location.CheckinLocation) {
		return nil, errors.New("chưa lấy được vị trí của bạn, vui lòng bật định vị và thử lại")
	}

	now := time.Now()
	fields := map[string]interface{}{}
	saveCheckinLocation(deposit, fields, false, location.CheckinLocation, now)

	// Chỉ tính khách có mặt khi ở gần BĐS; ở xa thì vẫn lưu vị trí nhưng chưa phải bằng chứng
	if atEstate(deposit.RealEstate, location.CheckinLocation) {
		customerChecked := true
		fields["customer_checkin"] = customerChecked
		deposit.CustomerCheckin = &customerChecked
	}
	if err := s.depositRepo.UpdateFields(deposit.ID, fields); err != nil {
		return nil, err
	}
	if err := s.refreshCheckinMatch(deposit); err != nil {
		return nil, err
	}
	if err := s.autoCheckinIfMatched(deposit); err != nil {
		return nil, err
	}

	return s.GetDepositDetail(deposit.ID, customerID)
}

// CustomerCheckin — khách nhập OTP môi giới hiển thị tại chỗ → CHECKED_IN.
func (s *depositService) CustomerCheckin(depositID, customerID uint64, otp string, location dto.CheckinLocationRequest) (*dto.DepositResponse, error) {
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

	now := time.Now()
	fields := map[string]interface{}{}
	// Vị trí của khách là bằng chứng bổ trợ để đối chiếu với vị trí môi giới
	saveCheckinLocation(deposit, fields, false, location.CheckinLocation, now)

	customerChecked := true
	// Hạn giữ phí: hết buổi xem + 4 ngày để khách kịp đặt cọc mua BĐS
	deadline := s.feeHoldDeadline(deposit)
	fields["status"] = model.DepositStatusCheckedIn
	fields["customer_checkin"] = customerChecked
	fields["otp_hash"] = ""
	fields["otp_expires_at"] = nil
	fields["report_deadline"] = deadline
	if err := s.depositRepo.UpdateFields(deposit.ID, fields); err != nil {
		return nil, err
	}
	if err := s.refreshCheckinMatch(deposit); err != nil {
		return nil, err
	}

	deposit.Status = model.DepositStatusCheckedIn
	deposit.CustomerCheckin = &customerChecked
	s.notifyBoth(deposit, "checked_in",
		fmt.Sprintf("Check-in xem nhà #%d thành công", deposit.ID),
		fmt.Sprintf("Xác nhận bạn đã gặp môi giới tại buổi xem nhà. Phí môi giới được giữ trong %d ngày: nếu bạn đặt cọc mua bất động sản trong thời gian này, phí sẽ được hoàn 100%% cho bạn.", s.cfg.RefundWindowDays),
		fmt.Sprintf("Khách đã check-in tại buổi xem nhà #%d. Phí môi giới sẽ được chuyển cho bạn sau %d ngày nếu khách không đặt cọc mua bất động sản.", deposit.ID, s.cfg.RefundWindowDays),
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
// Báo cáo chỉ dùng cho ĐIỂM DANH khi không check-in được (có mặt / không đến) — phần
// mua/không mua đã bỏ: phí hoàn hay không do KHÁCH ĐẶT CỌC MUA BĐS quyết định, không do lời khai.
//
// Mỗi bên chỉ báo cáo được MỘT LẦN (không cho sửa) để bên khai sau không thể
// xem câu trả lời của bên kia rồi đổi cho khớp.
func (s *depositService) SubmitReport(depositID, userID uint64, req dto.ReportResultRequest) (*dto.DepositResponse, error) {
	// Phía báo cáo suy ra từ chính bản ghi đơn, không nhận từ client
	deposit, side, err := s.getOwnedDeposit(depositID, userID)
	if err != nil {
		return nil, err
	}

	isBroker := side == model.RoleBroker
	var allowed []string

	switch deposit.Status {
	case model.DepositStatusBrokerConfirmed:
		// Không check-in được → điểm danh có mặt / không đến
		allowed = []string{model.ReportAttended, model.ReportNoShow}
	case model.DepositStatusCheckedIn:
		return nil, errors.New("buổi xem đã check-in xong, không cần báo cáo thêm")
	default:
		return nil, errors.New("đơn đặt lịch không ở trạng thái có thể báo cáo")
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

// evaluateReports áp bảng quyết định khi đã có đủ báo cáo ĐIỂM DANH của 2 bên.
// (Nhánh báo cáo mua/không mua đã bỏ — xem settleOverdueByCheckinLog và releaseFeeAfterHold.)
func (s *depositService) evaluateReports(deposit *model.Deposit) error {
	switch deposit.Status {
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
		fmt.Sprintf("Đơn đặt lịch #%d đang chờ duyệt tài liệu mua nhà", deposit.ID),
		"Hai bên đã xác nhận bạn mua nhà. Tài liệu mua bán đang được admin kiểm tra, phí môi giới được giữ nguyên tại platform.",
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
		return nil, errors.New("không tìm thấy đơn đặt lịch")
	}
	if deposit.Status != model.DepositStatusPendingPurchase {
		return nil, errors.New("đơn đặt lịch không ở trạng thái chờ duyệt tài liệu mua nhà")
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
	// Chỉ đánh giá được sau khi buổi xem đã diễn ra hoặc đơn đã tất toán
	if !isRateableStatus(deposit.Status) {
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
		return nil, errors.New("không tìm thấy đơn đặt lịch")
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
		return nil, "", errors.New("không tìm thấy đơn đặt lịch")
	}

	if deposit.CustomerID == userID {
		return deposit, model.RoleCustomer, nil
	}
	if deposit.BrokerID == userID {
		return deposit, model.RoleBroker, nil
	}
	return nil, "", errors.New("bạn không phải khách hàng hoặc môi giới của đơn đặt lịch này")
}

// requireSide kiểm tra user có đúng vai trò yêu cầu ở đơn hay không (VD môi giới để xác nhận lịch).
func (s *depositService) requireSide(depositID, userID uint64, wantSide string) (*model.Deposit, error) {
	deposit, side, err := s.getOwnedDeposit(depositID, userID)
	if err != nil {
		return nil, err
	}
	if side != wantSide {
		if wantSide == model.RoleBroker {
			return nil, errors.New("bạn không phụ trách đơn đặt lịch này")
		}
		return nil, errors.New("đơn đặt lịch không thuộc tài khoản của bạn")
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
		PurchaseDepositAt:      formatOptionalTime(deposit.PurchaseDepositAt),
		RejectReason:  deposit.RejectReason,
		// Bằng chứng vị trí check-in: FE hiển thị cho 2 bên và admin đối chiếu
		BrokerCheckinAt:       formatOptionalTime(deposit.BrokerCheckinAt),
		CustomerCheckinAt:     formatOptionalTime(deposit.CustomerCheckinAt),
		BrokerCheckinAcc:      deposit.BrokerCheckinAcc,
		CustomerCheckinAcc:    deposit.CustomerCheckinAcc,
		CheckinDistanceMeters: deposit.CheckinDistanceMeters,
		CheckinMatched:        deposit.CheckinMatched,
		CreatedAt:     deposit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     deposit.UpdatedAt.Format(time.RFC3339),
	}

	if deposit.Customer != nil {
		resp.CustomerName = deposit.Customer.Name
		resp.CustomerPhone = deposit.Customer.Phone
		resp.CustomerEmail = deposit.Customer.Email
	}
	// Ưu tiên liên hệ khách để lại lúc đặt lịch để môi giới gọi đúng số của buổi xem
	if name, phone := contactOf(deposit); name != "" || phone != "" {
		resp.CustomerName, resp.CustomerPhone = name, phone
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
	resp.CanReport = deposit.Status == model.DepositStatusBrokerConfirmed
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
// Khách chỉ trả PHÍ MÔI GIỚI (Amount = BrokerFee) nên không còn phần phí môi giới hoàn lại.
func (s *depositService) planForStatus(status string, deposit *model.Deposit) settlementPlan {
	switch status {
	case model.DepositStatusBrokerRejected:
		return settlementPlan{Status: status, Refund: deposit.Amount, Note: "Môi giới từ chối / quá hạn xác nhận → hoàn 100% phí"}
	case model.DepositStatusVisitedBought:
		return settlementPlan{Status: status, Refund: deposit.Amount, Note: "Khách đến và mua nhà → hoàn 100% phí môi giới"}
	case model.DepositStatusVisitedNotBuy:
		return settlementPlan{
			Status:   status,
			Refund:   deposit.Amount - deposit.BrokerFee,
			Transfer: deposit.BrokerFee,
			Note:     "Khách đến nhưng không mua → môi giới nhận phí, khách không được hoàn",
		}
	case model.DepositStatusNoShowCustomer:
		return settlementPlan{Status: status, Transfer: deposit.Amount, Note: "Khách không đến → môi giới nhận toàn bộ phí"}
	case model.DepositStatusNoShowBroker:
		return settlementPlan{
			Status:   status,
			Refund:   deposit.Amount,
			Penalty:  deposit.BrokerFee,
			Note:     "Môi giới không đến → hoàn 100% phí cho khách + phạt môi giới",
		}
	case model.DepositStatusRefunded:
		// Dùng cho ca không bên nào xác nhận buổi xem: hoàn lại toàn bộ phí cho khách
		return settlementPlan{Status: status, Refund: deposit.Amount, Note: "Hoàn 100% phí cho khách"}
	case model.DepositStatusCompleted:
		// Hết thời gian giữ phí mà khách không đặt cọc mua BĐS → phí thuộc về môi giới
		return settlementPlan{Status: status, Transfer: deposit.Amount, Note: "Hết thời gian giữ phí → môi giới nhận phí"}
	}
	return settlementPlan{Status: status, Note: "Kết thúc đơn đặt lịch"}
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
		fmt.Sprintf("Kết quả đơn đặt lịch xem nhà #%d", deposit.ID),
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

	now := time.Now()
	todayStr := now.Format(dateLayout)
	if viewingDate.Format(dateLayout) < todayStr {
		return time.Time{}, "", "", errors.New("ngày xem nhà không được ở quá khứ")
	}

	start := strings.TrimSpace(startStr)
	end := strings.TrimSpace(endStr)
	startAt, err := time.Parse(timeLayout, start)
	if err != nil {
		return time.Time{}, "", "", errors.New("giờ bắt đầu không hợp lệ (định dạng HH:MM)")
	}
	endAt, err := time.Parse(timeLayout, end)
	if err != nil {
		return time.Time{}, "", "", errors.New("giờ kết thúc không hợp lệ (định dạng HH:MM)")
	}
	if !startAt.Before(endAt) {
		return time.Time{}, "", "", errors.New("giờ kết thúc phải sau giờ bắt đầu")
	}

	// Mỗi buổi xem chỉ có 1 khung cố định 1 tiếng: bắt đầu ở mốc giờ chẵn trong giờ làm việc
	if endAt.Sub(startAt) != time.Duration(viewingSlotHours)*time.Hour || startAt.Minute() != 0 {
		return time.Time{}, "", "", fmt.Errorf("mỗi khung xem nhà kéo dài đúng %d tiếng, bắt đầu từ mốc giờ chẵn", viewingSlotHours)
	}
	if startAt.Hour() < viewingSlotFirstHour || endAt.Hour() > viewingSlotLastHour {
		return time.Time{}, "", "", fmt.Errorf(
			"khung giờ xem nhà chỉ trong khoảng %02d:00 - %02d:00", viewingSlotFirstHour, viewingSlotLastHour)
	}

	// Đặt trong ngày hôm nay thì khung phải còn ở tương lai
	if viewingDate.Format(dateLayout) == todayStr && !startAt.After(now) {
		return time.Time{}, "", "", errors.New("khung giờ xem nhà hôm nay phải sau thời điểm hiện tại")
	}

	// Chuẩn hoá lại "HH:MM" để cột viewing_start/viewing_end luôn cùng định dạng
	return viewingDate, startAt.Format(timeLayout), endAt.Format(timeLayout), nil
}

// normalizeContactInfo kiểm tra họ tên + SĐT khách để lại cho buổi xem.
// Trả về giá trị đã chuẩn hoá để lưu snapshot lên deposit.
func normalizeContactInfo(nameStr, phoneStr string) (string, string, error) {
	name := strings.TrimSpace(nameStr)
	if length := len([]rune(name)); length < 2 || length > 100 {
		return "", "", errors.New("vui lòng nhập họ tên người liên hệ (2 - 100 ký tự)")
	}

	// Khách hay gõ SĐT kèm khoảng trắng/dấu chấm/gạch nối, hoặc dạng +84
	phone := strings.NewReplacer(" ", "", ".", "", "-", "", "(", "", ")", "").Replace(strings.TrimSpace(phoneStr))
	if strings.HasPrefix(phone, "+84") {
		phone = "0" + strings.TrimPrefix(phone, "+84")
	}
	if !isVietnamesePhone(phone) {
		return "", "", errors.New("số điện thoại liên hệ không hợp lệ (VD 0901234567)")
	}

	return name, phone, nil
}

// isVietnamesePhone — SĐT di động VN: 10 số, bắt đầu bằng 0
func isVietnamesePhone(phone string) bool {
	if len(phone) != 10 || phone[0] != '0' {
		return false
	}
	for _, digit := range phone {
		if digit < '0' || digit > '9' {
			return false
		}
	}
	return true
}

// contactOf trả về tên + SĐT liên hệ của buổi xem:
// ưu tiên thông tin khách nhập lúc đặt lịch, thiếu thì lấy theo hồ sơ user.
func contactOf(deposit *model.Deposit) (string, string) {
	name, phone := "", ""
	if deposit.Customer != nil {
		name, phone = deposit.Customer.Name, deposit.Customer.Phone
	}
	if deposit.ContactName != "" {
		name = deposit.ContactName
	}
	if deposit.ContactPhone != "" {
		phone = deposit.ContactPhone
	}
	return name, phone
}

// feeHoldDeadline — mốc hết hạn GIỮ PHÍ: hết giờ buổi xem + số ngày giữ phí.
// Trong khoảng này khách đặt cọc mua BĐS thì hoàn 100% phí; quá hạn thì phí thuộc môi giới.
func (s *depositService) feeHoldDeadline(deposit *model.Deposit) time.Time {
	_, endAt, err := viewingRange(deposit)
	if err != nil {
		// Không ghép được mốc thời gian thì tính từ lúc gọi để đơn không bị treo vô hạn
		endAt = time.Now()
	}
	return endAt.AddDate(0, 0, s.cfg.RefundWindowDays)
}

// roundCoord làm tròn toạ độ về ~3 chữ số thập phân (~100m).
// Chỉ cần độ chính xác cấp phường để đối chiếu, không lưu vị trí chính xác của người dùng.
func roundCoord(value float64) float64 {
	return math.Round(value*1000) / 1000
}

// distanceMeters — khoảng cách Haversine giữa 2 toạ độ (mét).
func distanceMeters(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusMeters = 6371000

	rad := math.Pi / 180
	dLat := (lat2 - lat1) * rad
	dLng := (lng2 - lng1) * rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLng/2)*math.Sin(dLng/2)

	return earthRadiusMeters * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// usableCheckinLocation — toạ độ gửi lên có dùng được làm bằng chứng không.
// accuracy = 0 nghĩa là client không báo sai số, vẫn chấp nhận.
func usableCheckinLocation(location dto.CheckinLocation) bool {
	if location.Latitude == 0 && location.Longitude == 0 {
		return false
	}
	if location.Latitude < -90 || location.Latitude > 90 {
		return false
	}
	if location.Longitude < -180 || location.Longitude > 180 {
		return false
	}
	if location.Accuracy < 0 {
		return false
	}
	return location.Accuracy == 0 || location.Accuracy <= checkinMaxAccuracyMeters
}

// atEstate — toạ độ có nằm trong bán kính bất động sản không.
// BĐS chưa có toạ độ ⇒ không kiểm chứng được, trả true để không chặn oan người đã tới thật.
func atEstate(estate *model.RealEstate, location dto.CheckinLocation) bool {
	if estate == nil || estate.Latitude == nil || estate.Longitude == nil {
		return true
	}
	return distanceMeters(*estate.Latitude, *estate.Longitude, location.Latitude, location.Longitude) <= checkinMatchRadiusMeters
}

// RefundFeeOnPurchaseDeposit — luồng ĐẶT CỌC MUA bất động sản (làm sau) gọi hàm này khi
// khách đặt cọc thành công: đánh dấu đã đặt cọc và hoàn 100% phí môi giới của buổi xem.
// Chỉ hoàn khi đơn còn trong thời gian giữ phí và chưa tất toán.
func (s *depositService) RefundFeeOnPurchaseDeposit(depositID uint64) error {
	deposit, err := s.depositRepo.GetByID(depositID)
	if err != nil {
		return errors.New("không tìm thấy đơn đặt lịch")
	}

	now := time.Now()
	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"purchase_deposit_at": now,
	}); err != nil {
		return err
	}
	deposit.PurchaseDepositAt = &now

	// Hết thời gian giữ phí thì phí đã thuộc môi giới → không hoàn nữa
	if deposit.ReportDeadline != nil && now.After(*deposit.ReportDeadline) {
		return errors.New("đã quá thời gian giữ phí, không hoàn phí môi giới được nữa")
	}
	// Đơn đã tất toán trước đó thì chỉ ghi nhận dấu hiệu, không hoàn tiền lần nữa
	if deposit.Status != model.DepositStatusCheckedIn {
		return nil
	}

	return s.settleWithNotify(deposit, model.DepositStatusRefunded,
		"Khách đặt cọc mua bất động sản → hoàn 100% phí môi giới")
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

// isRateableStatus — đánh giá môi giới được sau khi buổi xem đã diễn ra (đã check-in)
// hoặc sau khi đơn đã tất toán. Các trạng thái VISITED_* giữ lại cho dữ liệu đơn cũ.
func isRateableStatus(status string) bool {
	switch status {
	case model.DepositStatusCheckedIn, model.DepositStatusCompleted, model.DepositStatusRefunded,
		model.DepositStatusVisitedBought, model.DepositStatusVisitedNotBuy,
		model.DepositStatusNoShowCustomer, model.DepositStatusNoShowBroker:
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
		return "Khách đến nhưng không mua (môi giới nhận phí)"
	case model.DepositStatusNoShowCustomer:
		return "Khách không đến (môi giới nhận phí)"
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
