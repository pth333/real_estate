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
