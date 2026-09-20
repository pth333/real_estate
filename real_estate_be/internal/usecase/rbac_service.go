package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/repo"
)

// IRbacService — quản lý role/permission và gán role cho user.
type IRbacService interface {
	// GetAccess trả role + permission của user (dùng cho login)
	GetAccess(userID uint64) (roles []string, permissions []string, err error)

	ListRoles() ([]dto.RoleResponse, error)
	ListPermissions() ([]dto.PermissionResponse, error)

	ListUsers(search string, page, size int) ([]dto.AdminUserResponse, int64, error)
	SetUserRoles(userID uint64, req dto.SetUserRolesRequest) (*dto.AdminUserResponse, error)
}

type rbacService struct {
	rbacRepo repo.IRbacRepository
}

func NewRbacService(rbacRepo repo.IRbacRepository) IRbacService {
	return &rbacService{rbacRepo: rbacRepo}
}

func (s *rbacService) GetAccess(userID uint64) ([]string, []string, error) {
	access, err := s.rbacRepo.GetUserAccess(userID)
	if err != nil {
		return nil, nil, err
	}
	return access.Roles, access.Permissions, nil
}

func (s *rbacService) ListRoles() ([]dto.RoleResponse, error) {
	roles, err := s.rbacRepo.ListRoles()
	if err != nil {
		return nil, err
	}

	result := make([]dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		permissions, err := s.rbacRepo.GetPermissionsByRoleIDs([]uint64{role.ID})
		if err != nil {
			return nil, err
		}

		codes := make([]string, 0, len(permissions))
		for _, permission := range permissions {
			codes = append(codes, permission.Code)
		}

		result = append(result, dto.RoleResponse{
			ID:          role.ID,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			IsActive:    role.IsActive == 1,
			Permissions: codes,
		})
	}
	return result, nil
}

func (s *rbacService) ListPermissions() ([]dto.PermissionResponse, error) {
	permissions, err := s.rbacRepo.ListPermissions()
	if err != nil {
		return nil, err
	}

	result := make([]dto.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		result = append(result, dto.PermissionResponse{
			ID:          permission.ID,
			Code:        permission.Code,
			Name:        permission.Name,
			Module:      permission.Module,
			Description: permission.Description,
		})
	}
	return result, nil
}

func (s *rbacService) ListUsers(search string, page, size int) ([]dto.AdminUserResponse, int64, error) {
	offset, limit := paginate(page, size)

	users, total, err := s.rbacRepo.ListUsers(search, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(users) == 0 {
		return []dto.AdminUserResponse{}, total, nil
	}

	ids := make([]uint64, 0, len(users))
	for _, user := range users {
		ids = append(ids, user.ID)
	}
	roleMap, err := s.rbacRepo.GetRolesByUserIDs(ids)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.AdminUserResponse, 0, len(users))
	for _, user := range users {
		roles := roleMap[user.ID]
		if roles == nil {
			roles = []string{}
		}
		result = append(result, dto.AdminUserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			IsActive:  user.IsActive == 1,
			Roles:     roles,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, total, nil
}

// SetUserRoles thay toàn bộ role của user. Yêu cầu ít nhất 1 role còn hiệu lực
// để tránh tài khoản mất hết quyền truy cập.
func (s *rbacService) SetUserRoles(userID uint64, req dto.SetUserRolesRequest) (*dto.AdminUserResponse, error) {
	user, err := s.rbacRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("không tìm thấy người dùng")
	}

	// Chuẩn hoá mã role: viết hoa, bỏ khoảng trắng, bỏ trùng
	codes := make([]string, 0, len(req.Roles))
	seen := make(map[string]bool)
	for _, code := range req.Roles {
		normalized := strings.ToUpper(strings.TrimSpace(code))
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true
		codes = append(codes, normalized)
	}
	if len(codes) == 0 {
		return nil, errors.New("phải chọn ít nhất 1 role")
	}

	roleIDs, err := s.rbacRepo.GetRoleIDsByCodes(codes)
	if err != nil {
		return nil, err
	}
	if len(roleIDs) != len(codes) {
		return nil, fmt.Errorf("có role không tồn tại hoặc đã bị khoá: %s", strings.Join(codes, ", "))
	}

	if err := s.rbacRepo.SetUserRoles(userID, roleIDs); err != nil {
		return nil, err
	}

	access, err := s.rbacRepo.GetUserAccess(userID)
	if err != nil {
		return nil, err
	}

	return &dto.AdminUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		IsActive:  user.IsActive == 1,
		Roles:     access.Roles,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
