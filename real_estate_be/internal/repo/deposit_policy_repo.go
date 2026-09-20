package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type IDepositPolicyRepository interface {
	// GetByPrice tra chính sách cọc theo giá BĐS (khoảng giá nửa khoảng [min, max))
	GetByPrice(priceVND float64) (*model.DepositPolicy, error)
	ListActive() ([]model.DepositPolicy, error)
}

type depositPolicyRepo struct {
	db *gorm.DB
}

func NewDepositPolicyRepository(db *gorm.DB) IDepositPolicyRepository {
	return &depositPolicyRepo{db: db}
}

// GetByPrice lấy chính sách khớp khoảng giá. Ưu tiên khoảng hẹp nhất (price_min lớn nhất)
// để sau này có thể thêm chính sách riêng cho phân khúc nhỏ nằm trong khoảng lớn.
func (r *depositPolicyRepo) GetByPrice(priceVND float64) (*model.DepositPolicy, error) {
	var policies []model.DepositPolicy
	err := r.db.
		Where("is_active = 1").
		Where("(price_min IS NULL OR price_min <= ?)", priceVND).
		Where("(price_max IS NULL OR price_max > ?)", priceVND).
		Order("COALESCE(price_min, 0) DESC").
		Limit(1).
		Find(&policies).Error
	if err != nil {
		return nil, err
	}
	if len(policies) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &policies[0], nil
}

func (r *depositPolicyRepo) ListActive() ([]model.DepositPolicy, error) {
	var policies []model.DepositPolicy
	err := r.db.
		Where("is_active = 1").
		Order("COALESCE(price_min, 0) ASC").
		Find(&policies).Error
	return policies, err
}
