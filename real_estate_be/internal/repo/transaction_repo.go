package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type ITransactionRepository interface {
	Create(tx *model.Transaction) error
	ListByDeposit(depositID uint64) ([]model.Transaction, error)
	// SumByType trả về tổng tiền theo từng loại giao dịch (dùng tính số dư escrow)
	SumByType() (map[string]float64, error)
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepo) ListByDeposit(depositID uint64) ([]model.Transaction, error) {
	var items []model.Transaction
	err := r.db.Where("deposit_id = ?", depositID).Order("created_at ASC").Find(&items).Error
	return items, err
}

func (r *transactionRepo) SumByType() (map[string]float64, error) {
	type row struct {
		Type  string
		Total float64
	}
	var rows []row
	err := r.db.Model(&model.Transaction{}).
		Select("type, COALESCE(SUM(amount), 0) AS total").
		Group("type").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	sums := make(map[string]float64, len(rows))
	for _, item := range rows {
		sums[item.Type] = item.Total
	}
	return sums, nil
}
