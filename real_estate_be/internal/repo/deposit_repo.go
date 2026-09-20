package repo

import (
	"strings"
	"time"

	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type IDepositRepository interface {
	Create(deposit *model.Deposit) error
	Save(deposit *model.Deposit) error
	GetByID(id uint64) (*model.Deposit, error)
	GetByPaymentRef(paymentRef string) (*model.Deposit, error)
	UpdateFields(id uint64, fields map[string]interface{}) error

	// CountOverlapSlot đếm số deposit đang giữ khung giờ trùng với khung giờ mới
	CountOverlapSlot(realEstateID uint64, viewingDate time.Time, start, end string) (int64, error)

	ListByCustomer(customerID uint64, status string, offset, limit int) ([]model.Deposit, int64, error)
	ListByBroker(brokerID uint64, status string, offset, limit int) ([]model.Deposit, int64, error)
	ListAll(status string, offset, limit int) ([]model.Deposit, int64, error)

	// ── Dùng cho cron jobs ──
	ListPendingBefore(before time.Time) ([]model.Deposit, error)
	ListReminderDue(viewingDate time.Time) ([]model.Deposit, error)
	ListReportDeadlinePassed(now time.Time) ([]model.Deposit, error)
	ListUnpaidBefore(before time.Time) ([]model.Deposit, error)
	// ListConfirmedWithoutDeadline — đơn đã xác nhận nhưng chưa mở cửa sổ báo cáo kết quả
	ListConfirmedWithoutDeadline() ([]model.Deposit, error)

	CountByStatus() (map[string]int64, error)
	ListCheckedInBefore(before time.Time) ([]model.Deposit, error)
}

type depositRepo struct {
	db *gorm.DB
}

func NewDepositRepository(db *gorm.DB) IDepositRepository {
	return &depositRepo{db: db}
}

// withRelations nạp sẵn khách, môi giới, BĐS và ảnh của BĐS để dựng response.
func (r *depositRepo) withRelations() *gorm.DB {
	return r.db.Preload("Customer").Preload("Broker").Preload("RealEstate").Preload("RealEstate.Images")
}

// normalizeJSONFields — MariaDB/MySQL không chấp nhận chuỗi rỗng ở cột JSON
// (CHECK json_valid), nên field bằng chứng trống phải ghi thành mảng rỗng "[]".
func normalizeJSONFields(deposit *model.Deposit) {
	if strings.TrimSpace(deposit.BrokerReportEvidence) == "" {
		deposit.BrokerReportEvidence = "[]"
	}
	if strings.TrimSpace(deposit.CustomerReportEvidence) == "" {
		deposit.CustomerReportEvidence = "[]"
	}
}

func (r *depositRepo) Create(deposit *model.Deposit) error {
	normalizeJSONFields(deposit)
	return r.db.Create(deposit).Error
}

func (r *depositRepo) Save(deposit *model.Deposit) error {
	normalizeJSONFields(deposit)
	return r.db.Save(deposit).Error
}

func (r *depositRepo) GetByID(id uint64) (*model.Deposit, error) {
	var deposit model.Deposit
	if err := r.withRelations().First(&deposit, id).Error; err != nil {
		return nil, err
	}
	return &deposit, nil
}

func (r *depositRepo) GetByPaymentRef(paymentRef string) (*model.Deposit, error) {
	var deposit model.Deposit
	if err := r.withRelations().Where("payment_ref = ?", paymentRef).First(&deposit).Error; err != nil {
		return nil, err
	}
	return &deposit, nil
}

// UpdateFields cập nhật một phần cột — dùng khi chỉ đổi trạng thái/otp để không ghi đè dữ liệu khác.
func (r *depositRepo) UpdateFields(id uint64, fields map[string]interface{}) error {
	return r.db.Model(&model.Deposit{}).Where("id = ?", id).Updates(fields).Error
}

// CountOverlapSlot kiểm tra trùng khung giờ của cùng 1 BĐS theo plan mục 5.
func (r *depositRepo) CountOverlapSlot(realEstateID uint64, viewingDate time.Time, start, end string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Deposit{}).
		Where("real_estate_id = ?", realEstateID).
		Where("viewing_date = ?", viewingDate.Format("2006-01-02")).
		Where("status NOT IN ?", model.DepositFinishedStatuses).
		Where("viewing_start < ?", end).
		Where("viewing_end > ?", start).
		Count(&count).Error
	return count, err
}

func (r *depositRepo) ListByCustomer(customerID uint64, status string, offset, limit int) ([]model.Deposit, int64, error) {
	query := r.db.Model(&model.Deposit{}).Where("customer_id = ?", customerID)
	return r.listWithQuery(query, status, offset, limit)
}

func (r *depositRepo) ListByBroker(brokerID uint64, status string, offset, limit int) ([]model.Deposit, int64, error) {
	query := r.db.Model(&model.Deposit{}).Where("broker_id = ?", brokerID)
	return r.listWithQuery(query, status, offset, limit)
}

func (r *depositRepo) ListAll(status string, offset, limit int) ([]model.Deposit, int64, error) {
	return r.listWithQuery(r.db.Model(&model.Deposit{}), status, offset, limit)
}

func (r *depositRepo) listWithQuery(query *gorm.DB, status string, offset, limit int) ([]model.Deposit, int64, error) {
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var deposits []model.Deposit
	err := query.
		Preload("Customer").
		Preload("Broker").
		Preload("RealEstate").
		Preload("RealEstate.Images").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&deposits).Error

	return deposits, total, err
}

// ListPendingBefore — deposit đã thanh toán nhưng quá hạn môi giới xác nhận.
func (r *depositRepo) ListPendingBefore(before time.Time) ([]model.Deposit, error) {
	var deposits []model.Deposit
	err := r.withRelations().
		Where("status = ?", model.DepositStatusPending).
		Where("paid_at IS NOT NULL AND paid_at < ?", before).
		Find(&deposits).Error
	return deposits, err
}

// ListReminderDue — deposit đã xác nhận và ngày xem là ngày truyền vào, chưa gửi mail nhắc.
func (r *depositRepo) ListReminderDue(viewingDate time.Time) ([]model.Deposit, error) {
	var deposits []model.Deposit
	err := r.withRelations().
		Where("status = ?", model.DepositStatusBrokerConfirmed).
		Where("viewing_date = ?", viewingDate.Format("2006-01-02")).
		Where("reminder_sent_at IS NULL").
		Find(&deposits).Error
	return deposits, err
}

// ListReportDeadlinePassed — quá hạn 2 bên báo cáo mà vẫn chưa check-in → đẩy sang DISPUTE.
func (r *depositRepo) ListReportDeadlinePassed(now time.Time) ([]model.Deposit, error) {
	var deposits []model.Deposit
	err := r.withRelations().
		Where("status = ?", model.DepositStatusBrokerConfirmed).
		Where("report_deadline IS NOT NULL AND report_deadline < ?", now).
		Find(&deposits).Error
	return deposits, err
}

// ListUnpaidBefore — deposit chưa thanh toán quá hạn → tự huỷ để trả lại khung giờ.
func (r *depositRepo) ListUnpaidBefore(before time.Time) ([]model.Deposit, error) {
	var deposits []model.Deposit
	err := r.db.
		Where("status = ?", model.DepositStatusAwaitingPayment).
		Where("created_at < ?", before).
		Find(&deposits).Error
	return deposits, err
}

func (r *depositRepo) CountByStatus() (map[string]int64, error) {
	type row struct {
		Status string
		Total  int64
	}
	var rows []row
	err := r.db.Model(&model.Deposit{}).
		Select("status, COUNT(*) AS total").
		Group("status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64, len(rows))
	for _, item := range rows {
		counts[item.Status] = item.Total
	}
	return counts, nil
}

// ListConfirmedWithoutDeadline — đơn môi giới đã xác nhận nhưng chưa mở cửa sổ báo cáo.
// Service tự so mốc "viewing_start + ân hạn" trong Go để không phụ thuộc hàm ngày giờ của DB.
func (r *depositRepo) ListConfirmedWithoutDeadline() ([]model.Deposit, error) {
	var deposits []model.Deposit
	err := r.withRelations().
		Where("status = ?", model.DepositStatusBrokerConfirmed).
		Where("report_deadline IS NULL").
		Find(&deposits).Error
	return deposits, err
}

// ListCheckedInBefore — đã check-in nhưng quá hạn mà 2 bên chưa báo cáo đủ → chuyển admin xử lý.
func (r *depositRepo) ListCheckedInBefore(before time.Time) ([]model.Deposit, error) {
	var deposits []model.Deposit
	err := r.withRelations().
		Where("status = ?", model.DepositStatusCheckedIn).
		Where("report_deadline IS NOT NULL AND report_deadline < ?", before).
		Find(&deposits).Error
	return deposits, err
}
