package repo

import (
	"time"

	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

// UserAccess — quyền truy cập đã gom sẵn của 1 user (dùng cho middleware + login)
type UserAccess struct {
	Roles       []string
	Permissions []string
}

type IRbacRepository interface {
	// ── Đọc quyền ──
	GetUserAccess(userID uint64) (*UserAccess, error)
	GetRolesByUserIDs(userIDs []uint64) (map[uint64][]string, error)

	// ── Role / permission ──
	ListRoles() ([]model.Role, error)
	ListPermissions() ([]model.Permission, error)
	GetRoleIDsByCodes(codes []string) ([]uint64, error)
	GetPermissionsByRoleIDs(roleIDs []uint64) ([]model.Permission, error)

	// ── Gán role cho user ──
	SetUserRoles(userID uint64, roleIDs []uint64) error
	HasRole(userID uint64, roleCode string) (bool, error)
	AddRoleByCode(userID uint64, roleCode string) error

	// ── Danh sách user cho admin ──
	ListUsers(search string, offset, limit int) ([]model.User, int64, error)
	GetUserByID(userID uint64) (*model.User, error)
}

type rbacRepo struct {
	db *gorm.DB
}

func NewRbacRepository(db *gorm.DB) IRbacRepository {
	return &rbacRepo{db: db}
}

// GetUserAccess lấy danh sách role code + permission code của user.
// Chỉ lấy role đang active.
func (r *rbacRepo) GetUserAccess(userID uint64) (*UserAccess, error) {
	access := &UserAccess{Roles: []string{}, Permissions: []string{}}

	var roleIDs []uint64
	if err := r.db.Model(&model.UserRole{}).
		Joins("JOIN roles ON roles.id = user_roles.role_id AND roles.is_active = 1").
		Where("user_roles.user_id = ?", userID).
		Pluck("user_roles.role_id", &roleIDs).Error; err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return access, nil
	}

	if err := r.db.Model(&model.Role{}).
		Where("id IN ?", roleIDs).
		Order("code ASC").
		Pluck("code", &access.Roles).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&model.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Distinct().
		Order("permissions.code ASC").
		Pluck("permissions.code", &access.Permissions).Error; err != nil {
		return nil, err
	}

	return access, nil
}

// GetRolesByUserIDs lấy role code của nhiều user trong 1 query (dùng cho danh sách admin).
func (r *rbacRepo) GetRolesByUserIDs(userIDs []uint64) (map[uint64][]string, error) {
	result := make(map[uint64][]string)
	if len(userIDs) == 0 {
		return result, nil
	}

	type row struct {
		UserID uint64
		Code   string
	}
	var rows []row
	err := r.db.Model(&model.UserRole{}).
		Select("user_roles.user_id AS user_id, roles.code AS code").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id IN ?", userIDs).
		Order("roles.code ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, item := range rows {
		result[item.UserID] = append(result[item.UserID], item.Code)
	}
	return result, nil
}

func (r *rbacRepo) ListRoles() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Order("id ASC").Find(&roles).Error
	return roles, err
}

func (r *rbacRepo) ListPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Order("module ASC, code ASC").Find(&permissions).Error
	return permissions, err
}

// GetRoleIDsByCodes đổi danh sách role code → id, bỏ qua code không tồn tại.
func (r *rbacRepo) GetRoleIDsByCodes(codes []string) ([]uint64, error) {
	if len(codes) == 0 {
		return []uint64{}, nil
	}
	var ids []uint64
	err := r.db.Model(&model.Role{}).
		Where("code IN ?", codes).
		Where("is_active = 1").
		Pluck("id", &ids).Error
	return ids, err
}

func (r *rbacRepo) GetPermissionsByRoleIDs(roleIDs []uint64) ([]model.Permission, error) {
	var permissions []model.Permission
	if len(roleIDs) == 0 {
		return permissions, nil
	}
	err := r.db.Model(&model.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Distinct().
		Order("permissions.module ASC, permissions.code ASC").
		Find(&permissions).Error
	return permissions, err
}

// SetUserRoles thay toàn bộ role của user bằng danh sách roleIDs (xoá hết rồi thêm lại).
func (r *rbacRepo) SetUserRoles(userID uint64, roleIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}

		now := time.Now()
		links := make([]model.UserRole, 0, len(roleIDs))
		for _, roleID := range roleIDs {
			links = append(links, model.UserRole{UserID: userID, RoleID: roleID, CreatedAt: now})
		}
		return tx.Create(&links).Error
	})
}

func (r *rbacRepo) HasRole(userID uint64, roleCode string) (bool, error) {
	var count int64
	err := r.db.Model(&model.UserRole{}).
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("roles.code = ?", roleCode).
		Count(&count).Error
	return count > 0, err
}

// AddRoleByCode gán thêm 1 role cho user, bỏ qua nếu đã có.
func (r *rbacRepo) AddRoleByCode(userID uint64, roleCode string) error {
	var role model.Role
	if err := r.db.Where("code = ?", roleCode).First(&role).Error; err != nil {
		return err
	}

	var count int64
	if err := r.db.Model(&model.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, role.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	return r.db.Create(&model.UserRole{UserID: userID, RoleID: role.ID, CreatedAt: time.Now()}).Error
}

func (r *rbacRepo) ListUsers(search string, offset, limit int) ([]model.User, int64, error) {
	query := r.db.Model(&model.User{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.User
	err := query.Order("id ASC").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *rbacRepo) GetUserByID(userID uint64) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
