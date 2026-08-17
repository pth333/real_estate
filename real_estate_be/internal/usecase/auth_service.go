package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	"real_estate_be/pkg/jwt"
	"real_estate_be/pkg/sms"

	"golang.org/x/crypto/bcrypt"
)

const (
	otpTTL      = 5 * time.Minute
	otpRedisKey = "otp:%s"
)

type AuthService struct {
	repo repo.IUserRepository
	sms  sms.Provider
}

type AuthServiceInterface interface {
	Register(req dto.CreateUserRequest) error
	Login(req dto.LoginRequest) (string, string, *dto.UserResponse, error)
	RefreshToken(refreshToken string) (string, string, error)
	SendOTP(req dto.SendOTPRequest) error
	VerifyOTP(req dto.VerifyOTPRequest, currentUserEmail string) error
}

func NewAuthService(repo repo.IUserRepository, smsProvider sms.Provider) AuthServiceInterface {
	return &AuthService{
		repo: repo,
		sms:  smsProvider,
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

func (h *AuthService) Login(req dto.LoginRequest) (string, string, *dto.UserResponse, error) {
	user, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		return "", "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", "", nil, err
	}

	accessToken, err := jwt.GenerateAccessToken(user.Email, user.ID)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.Email)
	if err != nil {
		return "", "", nil, err
	}

	userRes := &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return accessToken, refreshToken, userRes, nil
}

func (h *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	// Tìm user theo email để lấy user_id cho access token mới
	user, findErr := h.repo.FindByEmail(claims.Email)
	if findErr != nil {
		return "", "", errors.New("user not found")
	}

	newAccess, err := jwt.GenerateAccessToken(claims.Email, user.ID)
	if err != nil {
		return "", "", err
	}

	newRefresh, err := jwt.GenerateRefreshToken(claims.Email)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

func (s *AuthService) SendOTP(req dto.SendOTPRequest) error {
	ctx := context.Background()

	// Sinh mã OTP ngẫu nhiên
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Lưu vào Redis với TTL 5 phút
	key := fmt.Sprintf(otpRedisKey, req.Phone)
	if err := global.RedisClient.Set(ctx, key, otp, otpTTL).Err(); err != nil {
		return fmt.Errorf("failed to cache OTP: %w", err)
	}

	// Gửi OTP qua SMS
	if err := s.sms.Send(req.Phone, otp); err != nil {
		return fmt.Errorf("failed to send OTP: %w", err)
	}

	return nil
}

func (s *AuthService) VerifyOTP(req dto.VerifyOTPRequest, currentUserEmail string) error {
	ctx := context.Background()

	// Lấy OTP từ Redis
	key := fmt.Sprintf(otpRedisKey, req.Phone)
	cachedOTP, err := global.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("OTP không hợp lệ hoặc đã hết hạn")
	}

	if cachedOTP != req.OTP {
		return fmt.Errorf("OTP không đúng")
	}

	// Xoá OTP khỏi Redis (chỉ dùng 1 lần)
	global.RedisClient.Del(ctx, key)

	// Nếu có user đang đăng nhập hiện tại, lưu thẳng số điện thoại vào tài khoản của họ
	if currentUserEmail != "" {
		if err := s.repo.UpdatePhoneByEmail(currentUserEmail, req.Phone); err != nil {
			return fmt.Errorf("không thể cập nhật số điện thoại cho tài khoản: %w", err)
		}
		return nil
	}

	// Trường hợp Public (Không đăng nhập, ví dụ đăng nhập trực tiếp qua SĐT)
	// Kiểm tra user đã tồn tại chưa
	_, findErr := s.repo.FindByPhone(req.Phone)
	if findErr != nil {
		// Chưa có → tạo mới
		_, err = s.repo.CreateUserByPhone(req.Phone)
		if err != nil {
			return fmt.Errorf("không thể tạo tài khoản: %w", err)
		}
	} else {
		// Đã có → đánh dấu verified
		if err := s.repo.MarkPhoneVerified(req.Phone); err != nil {
			return fmt.Errorf("không thể xác thực số điện thoại: %w", err)
		}
	}

	return nil
}
