package usecase

import (
	"errors"
	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	"real_estate_be/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repo.IUserRepository
}

type AuthServiceInterface interface {
	Register(req dto.CreateUserRequest) error
	Login(req dto.LoginRequest) (string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
}

func NewAuthService(repo repo.IUserRepository) AuthServiceInterface {
	return &AuthService{
		repo: repo,
	}
}

func (h *AuthService) Register(req dto.CreateUserRequest) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	user := &model.User{
		Email:    req.Email,
		Password: string(hash),
		Name:     req.Name,
	}
	return h.repo.Register(user)
}

func (h *AuthService) Login(req dto.LoginRequest) (string, string, error) {
	user, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", "", err
	}

	accessToken, err := jwt.GenerateAccessToken(user.Email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.Email)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (h *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	// Tạo cặp token mới
	newAccess, err := jwt.GenerateAccessToken(claims.Email)
	if err != nil {
		return "", "", err
	}

	newRefresh, err := jwt.GenerateRefreshToken(claims.Email)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}
