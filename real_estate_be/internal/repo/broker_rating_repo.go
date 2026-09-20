package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type IBrokerRatingRepository interface {
	Create(rating *model.BrokerRating) error
	GetByDeposit(depositID uint64) (*model.BrokerRating, error)
	// RecomputedRating trả về điểm trung bình mới và tổng số đánh giá của môi giới
	RecomputedRating(brokerID uint64) (float64, int, error)
}

type brokerRatingRepo struct {
	db *gorm.DB
}

func NewBrokerRatingRepository(db *gorm.DB) IBrokerRatingRepository {
	return &brokerRatingRepo{db: db}
}

func (r *brokerRatingRepo) Create(rating *model.BrokerRating) error {
	return r.db.Create(rating).Error
}

// GetByDeposit lấy đánh giá của 1 đơn. Dùng Find để không log lỗi khi chưa có đánh giá.
func (r *brokerRatingRepo) GetByDeposit(depositID uint64) (*model.BrokerRating, error) {
	var ratings []model.BrokerRating
	err := r.db.Where("deposit_id = ?", depositID).Limit(1).Find(&ratings).Error
	if err != nil {
		return nil, err
	}
	if len(ratings) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &ratings[0], nil
}

func (r *brokerRatingRepo) RecomputedRating(brokerID uint64) (float64, int, error) {
	type row struct {
		Avg   float64
		Total int
	}
	var result row
	err := r.db.Model(&model.BrokerRating{}).
		Select("COALESCE(AVG(rating), 5) AS avg, COUNT(*) AS total").
		Where("broker_id = ?", brokerID).
		Scan(&result).Error
	if err != nil {
		return 5, 0, err
	}
	return result.Avg, result.Total, nil
}
