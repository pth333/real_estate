package usecase

import (
	"real_estate_be/internal/delivery/https/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repository"
	"real_estate_be/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (h *AuthService) Register(req dto.CreateUserRequest) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	user := &model.User{
		Email:    req.Email,
		Password: string(hash),
		Name:     req.Name,
	}
	return h.userRepo.Register(user)

}

func (h *AuthService) Login(req dto.LoginRequest) (string, error) {
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	return jwt.GenerateToken(user.Email)
}
