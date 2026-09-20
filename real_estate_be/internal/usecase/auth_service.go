package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
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
	repo     repo.IUserRepository
	sms      sms.Provider
	rbacRepo repo.IRbacRepository
}

type AuthServiceInterface interface {
	Register(req dto.CreateUserRequest) error
	Login(req dto.LoginRequest) (string, string, *dto.UserResponse, error)
	RefreshToken(refreshToken string) (string, string, error)
	SendOTP(req dto.SendOTPRequest) error
	VerifyOTP(req dto.VerifyOTPRequest, currentUserEmail string) error
	// GetUserCurrentInfo trả thông tin user đang đăng nhập (role + permission)
	// để FE ẩn/hiện UI và chặn ở tầng route.
	GetUserCurrentInfo(userID uint64) (*dto.UserResponse, error)
}

func NewAuthService(repo repo.IUserRepository, smsProvider sms.Provider, rbacRepo repo.IRbacRepository) AuthServiceInterface {
	return &AuthService{
		repo:     repo,
		sms:      smsProvider,
		rbacRepo: rbacRepo,
	}
}

func (h *AuthService) Register(req dto.CreateUserRequest) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	user := &model.User{
		Email:    req.Email,
		Password: string(hash),
		Name:     req.Name,
	}
	if err := h.repo.Register(user); err != nil {
		return err
	}

	// Mặc định user mới là khách hàng.
	return h.rbacRepo.AddRoleByCode(user.ID, model.RoleCustomer)
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

	userRes, err := h.buildUserResponse(user)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, userRes, nil
}

// buildUserResponse gom role + permission của user để FE điều hướng giao diện.
func (h *AuthService) buildUserResponse(user *model.User) (*dto.UserResponse, error) {
	access, err := h.rbacRepo.GetUserAccess(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		Roles:       access.Roles,
		Permissions: access.Permissions,
	}, nil
}

// GetUserCurrentInfo đọc lại user + quyền từ DB theo user_id trong token.
// Dùng cho GET /auth/user-current-info: FE gọi lúc khởi động để luôn có role/permission
// mới nhất (sau khi admin đổi role không cần đăng nhập lại mới thấy đúng giao diện).
func (h *AuthService) GetUserCurrentInfo(userID uint64) (*dto.UserResponse, error) {
	user, err := h.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("không tìm thấy tài khoản")
	}
	if user.IsActive == 0 {
		return nil, errors.New("tài khoản đã bị khoá")
	}
	return h.buildUserResponse(user)
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

	if strings.HasPrefix(req.Phone, "0") {
		numberPhone := "+84" + req.Phone[1:]
		req.Phone = numberPhone
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
	existing, findErr := s.repo.FindByPhone(req.Phone)
	if findErr != nil {
		// Chưa có → tạo mới
		created, err := s.repo.CreateUserByPhone(req.Phone)
		if err != nil {
			return fmt.Errorf("không thể tạo tài khoản: %w", err)
		}
		// Tài khoản tạo qua SĐT mặc định là khách hàng
		if err := s.rbacRepo.AddRoleByCode(created.ID, model.RoleCustomer); err != nil {
			return fmt.Errorf("không thể gán role cho tài khoản: %w", err)
		}
	} else {
		// Đã có → đánh dấu verified
		if err := s.repo.MarkPhoneVerified(req.Phone); err != nil {
			return fmt.Errorf("không thể xác thực số điện thoại: %w", err)
		}
		if existing != nil && existing.ID != 0 {
			// Đảm bảo tài khoản cũ luôn có ít nhất role khách hàng
			if err := s.rbacRepo.AddRoleByCode(existing.ID, model.RoleCustomer); err != nil {
				return fmt.Errorf("không thể gán role cho tài khoản: %w", err)
			}
		}
	}

	return nil
}
