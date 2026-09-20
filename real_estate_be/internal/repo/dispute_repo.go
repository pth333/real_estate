package repo

import (
	"strings"

	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type IDisputeRepository interface {
	Create(dispute *model.Dispute) error
	Save(dispute *model.Dispute) error
	GetByID(id uint64) (*model.Dispute, error)
	GetActiveByDeposit(depositID uint64) (*model.Dispute, error)
	List(status string, offset, limit int) ([]model.Dispute, int64, error)
	CountOpen() (int64, error)
}

type disputeRepo struct {
	db *gorm.DB
}

func NewDisputeRepository(db *gorm.DB) IDisputeRepository {
	return &disputeRepo{db: db}
}

func (r *disputeRepo) Create(dispute *model.Dispute) error {
	// Cột JSON không nhận chuỗi rỗng → mặc định là mảng rỗng
	if strings.TrimSpace(dispute.EvidenceURLs) == "" {
		dispute.EvidenceURLs = "[]"
	}
	return r.db.Create(dispute).Error
}

func (r *disputeRepo) Save(dispute *model.Dispute) error {
	return r.db.Save(dispute).Error
}

func (r *disputeRepo) GetByID(id uint64) (*model.Dispute, error) {
	var dispute model.Dispute
	err := r.db.
		Preload("Deposit").
		Preload("Deposit.Customer").
		Preload("Deposit.Broker").
		Preload("Deposit.RealEstate").
		First(&dispute, id).Error
	if err != nil {
		return nil, err
	}
	return &dispute, nil
}

// GetActiveByDeposit lấy dispute chưa xử lý xong của 1 deposit.
// Dùng Find thay vì First để trường hợp không có dispute không bị log lỗi.
func (r *disputeRepo) GetActiveByDeposit(depositID uint64) (*model.Dispute, error) {
	var disputes []model.Dispute
	err := r.db.
		Where("deposit_id = ?", depositID).
		Where("status <> ?", model.DisputeStatusResolved).
		Order("created_at DESC").
		Limit(1).
		Find(&disputes).Error
	if err != nil {
		return nil, err
	}
	if len(disputes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &disputes[0], nil
}

func (r *disputeRepo) List(status string, offset, limit int) ([]model.Dispute, int64, error) {
	query := r.db.Model(&model.Dispute{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Dispute
	err := query.
		Preload("Deposit").
		Preload("Deposit.Customer").
		Preload("Deposit.Broker").
		Preload("Deposit.RealEstate").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&items).Error

	return items, total, err
}

func (r *disputeRepo) CountOpen() (int64, error) {
	var count int64
	err := r.db.Model(&model.Dispute{}).
		Where("status <> ?", model.DisputeStatusResolved).
		Count(&count).Error
	return count, err
}
