package model

import "time"

// ══════════════════════════════════════════════════════════
// RBAC: user —(n-n)— role —(n-n)— permission
// Một user có thể giữ NHIỀU role, mỗi role gồm nhiều permission.
// Kiểm quyền ở route dựa trên PERMISSION (không dựa tên role) nên thêm role
// mới chỉ cần gán permission, không phải sửa code.
// ══════════════════════════════════════════════════════════

type Role struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"column:code;size:30;uniqueIndex" json:"code"`
	Name        string    `gorm:"column:name;size:100" json:"name"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	IsActive    int       `gorm:"column:is_active;default:1" json:"is_active"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Role) TableName() string { return "roles" }

type Permission struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"column:code;size:60;uniqueIndex" json:"code"`
	Name        string `gorm:"column:name;size:150" json:"name"`
	// Module nhóm quyền theo nghiệp vụ: deposit / broker / admin
	Module      string    `gorm:"column:module;size:30;index" json:"module"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Permission) TableName() string { return "permissions" }

// RolePermission — bảng nối role ↔ permission
type RolePermission struct {
	RoleID       uint64 `gorm:"primaryKey;column:role_id" json:"role_id"`
	PermissionID uint64 `gorm:"primaryKey;column:permission_id" json:"permission_id"`
}

func (RolePermission) TableName() string { return "role_permissions" }

// UserRole — bảng nối user ↔ role (thay cho cột users.role cũ)
type UserRole struct {
	UserID    uint64    `gorm:"primaryKey;column:user_id" json:"user_id"`
	RoleID    uint64    `gorm:"primaryKey;column:role_id" json:"role_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (UserRole) TableName() string { return "user_roles" }

// ── Mã role (giá trị của roles.code) ──
const (
	RoleCustomer = "CUSTOMER"
	RoleBroker   = "BROKER"
	RoleAdmin    = "ADMIN"
)

// ── Mã permission (giá trị của permissions.code) ──
const (
	// Khách hàng
	PermissionDepositCreate  = "deposit.create"
	PermissionDepositViewOwn = "deposit.view.own"
	PermissionDepositCheckin = "deposit.checkin"
	PermissionDepositRate    = "deposit.rate"

	// Dùng chung khách + môi giới của đơn
	PermissionDepositView     = "deposit.view"
	PermissionDepositReport   = "deposit.report"
	PermissionDepositDispute  = "deposit.dispute"
	PermissionDisputeEvidence = "dispute.evidence"

	// Môi giới
	PermissionBrokerDepositList    = "broker.deposit.list"
	PermissionBrokerDepositConfirm = "broker.deposit.confirm"
	PermissionBrokerDepositReject  = "broker.deposit.reject"
	PermissionBrokerDepositOtp     = "broker.deposit.otp"

	// Admin
	PermissionAdminDepositList     = "admin.deposit.list"
	PermissionAdminDepositView     = "admin.deposit.view"
	PermissionAdminDepositApprove  = "admin.deposit.approve"
	PermissionAdminEscrowView      = "admin.escrow.view"
	PermissionAdminDisputeList     = "admin.dispute.list"
	PermissionAdminDisputeView     = "admin.dispute.view"
	PermissionAdminDisputeResolve  = "admin.dispute.resolve"
	PermissionAdminUserList        = "admin.user.list"
	PermissionAdminUserAssignRole  = "admin.user.assign_role"
	PermissionAdminRoleList        = "admin.role.list"
)
