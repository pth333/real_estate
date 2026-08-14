package repo

import (
	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Register(item *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
	CreateUserByPhone(phone string) (*model.User, error)
	MarkPhoneVerified(phone string) error
	UpdatePhoneByEmail(email string, phone string) error
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(item *model.User) error {
	return r.db.Create(item).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUserByPhone(phone string) (*model.User, error) {
	user := &model.User{
		Phone:         phone,
		PhoneVerified: 1,
	}
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) MarkPhoneVerified(phone string) error {
	return r.db.Model(&model.User{}).
		Where("phone = ?", phone).
		Update("phone_verified", 1).
		Error
}

func (r *UserRepository) UpdatePhoneByEmail(email string, phone string) error {
	// 1. Kiểm tra xem số điện thoại này đã có ai dùng chưa
	var existingUser model.User
	err := r.db.Where("phone = ?", phone).First(&existingUser).Error
	if err == nil {
		// Nếu đã có user khác dùng số này:
		// Kiểm tra nếu user đó không có email (là tài khoản ảo tự tạo tự động bằng số điện thoại)
		if existingUser.Email == "" {
			// Xóa bản ghi ảo này đi để nhường chỗ cập nhật cho tài khoản chính
			r.db.Unscoped().Delete(&existingUser)
		} else if existingUser.Email != email {
			// Nếu là một tài khoản thật khác có email khác thì trả về lỗi trùng số điện thoại
			return r.db.Where("email = ?", email).Error // Trả về lỗi hoặc xử lý phù hợp
		}
	}

	// 2. Cập nhật số điện thoại và trạng thái đã xác thực cho user hiện tại theo email
	return r.db.Model(&model.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"phone":          phone,
			"phone_verified": 1,
		}).Error
}
