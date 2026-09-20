package initialize

import (
	"log"

	model "real_estate_be/internal/models"

	"gorm.io/gorm"
)

// permissionSeed — danh sách quyền theo hành động, nhóm theo module.
type permissionSeed struct {
	Code   string
	Name   string
	Module string
}

var permissionSeeds = []permissionSeed{
	// ── Khách hàng ──
	{model.PermissionDepositCreate, "Tạo đơn đặt cọc xem nhà", "deposit"},
	{model.PermissionDepositViewOwn, "Xem danh sách đơn của mình", "deposit"},
	{model.PermissionDepositCheckin, "Check-in OTP tại buổi xem", "deposit"},
	{model.PermissionDepositRate, "Đánh giá môi giới", "deposit"},

	// ── Dùng chung khách + môi giới của đơn ──
	{model.PermissionDepositView, "Xem chi tiết đơn mình liên quan", "deposit"},
	{model.PermissionDepositReport, "Báo cáo kết quả buổi xem", "deposit"},
	{model.PermissionDepositDispute, "Mở tranh chấp", "deposit"},
	{model.PermissionDisputeEvidence, "Gửi bằng chứng tranh chấp", "dispute"},

	// ── Môi giới ──
	{model.PermissionBrokerDepositList, "Xem danh sách đơn được giao", "broker"},
	{model.PermissionBrokerDepositConfirm, "Xác nhận lịch xem nhà", "broker"},
	{model.PermissionBrokerDepositReject, "Từ chối lịch xem nhà", "broker"},
	{model.PermissionBrokerDepositOtp, "Sinh mã OTP check-in", "broker"},

	// ── Admin ──
	{model.PermissionAdminDepositList, "Xem tất cả đơn đặt cọc", "admin"},
	{model.PermissionAdminDepositView, "Xem chi tiết mọi đơn đặt cọc", "admin"},
	{model.PermissionAdminDepositApprove, "Duyệt tài liệu mua nhà", "admin"},
	{model.PermissionAdminEscrowView, "Xem tổng quan escrow", "admin"},
	{model.PermissionAdminDisputeList, "Xem danh sách tranh chấp", "admin"},
	{model.PermissionAdminDisputeView, "Xem chi tiết tranh chấp", "admin"},
	{model.PermissionAdminDisputeResolve, "Ra quyết định xử lý tranh chấp", "admin"},
	{model.PermissionAdminUserList, "Xem danh sách người dùng", "admin"},
	{model.PermissionAdminUserAssignRole, "Gán role cho người dùng", "admin"},
	{model.PermissionAdminRoleList, "Xem danh sách role", "admin"},
}

// roleSeed — role kèm danh sách permission code được gán.
type roleSeed struct {
	Code        string
	Name        string
	Description string
	Permissions []string
}

var roleSeeds = []roleSeed{
	{
		Code:        model.RoleCustomer,
		Name:        "Khách hàng",
		Description: "Tìm kiếm BĐS, đặt cọc giữ lịch xem nhà, đánh giá môi giới",
		Permissions: []string{
			model.PermissionDepositCreate,
			model.PermissionDepositViewOwn,
			model.PermissionDepositView,
			model.PermissionDepositCheckin,
			model.PermissionDepositRate,
			model.PermissionDepositReport,
			model.PermissionDepositDispute,
			model.PermissionDisputeEvidence,
		},
	},
	{
		Code:        model.RoleBroker,
		Name:        "Môi giới",
		Description: "Đăng BĐS, xác nhận lịch, dẫn khách xem nhà, báo cáo kết quả",
		Permissions: []string{
			model.PermissionDepositView,
			model.PermissionDepositReport,
			model.PermissionDepositDispute,
			model.PermissionDisputeEvidence,
			model.PermissionBrokerDepositList,
			model.PermissionBrokerDepositConfirm,
			model.PermissionBrokerDepositReject,
			model.PermissionBrokerDepositOtp,
		},
	},
	{
		Code:        model.RoleAdmin,
		Name:        "Quản trị viên",
		Description: "Giữ escrow, duyệt tài liệu mua nhà, xử lý tranh chấp, gán role",
		Permissions: []string{
			model.PermissionAdminDepositList,
			model.PermissionAdminDepositView,
			model.PermissionAdminDepositApprove,
			model.PermissionAdminEscrowView,
			model.PermissionAdminDisputeList,
			model.PermissionAdminDisputeView,
			model.PermissionAdminDisputeResolve,
			model.PermissionAdminUserList,
			model.PermissionAdminUserAssignRole,
			model.PermissionAdminRoleList,
		},
	},
}

// seedRbac tạo role + permission + gán quyền cho role. Chạy lại nhiều lần vẫn đúng
// (chỉ thêm cái còn thiếu, không xoá quyền đã gán thủ công trong DB).
func seedRbac(db *gorm.DB) {
	// ── Permissions ──
	for _, seed := range permissionSeeds {
		var existing model.Permission
		err := db.Where("code = ?", seed.Code).Limit(1).Find(&existing).Error
		if err != nil {
			log.Fatalf("❌ [RBAC] đọc permission %s thất bại: %v", seed.Code, err)
		}
		if existing.ID != 0 {
			continue
		}
		if err := db.Create(&model.Permission{
			Code: seed.Code, Name: seed.Name, Module: seed.Module,
		}).Error; err != nil {
			log.Fatalf("❌ [RBAC] tạo permission %s thất bại: %v", seed.Code, err)
		}
	}

	// ── Roles + gán quyền ──
	createdLinks := 0
	for _, seed := range roleSeeds {
		var role model.Role
		err := db.Where("code = ?", seed.Code).Limit(1).Find(&role).Error
		if err != nil {
			log.Fatalf("❌ [RBAC] đọc role %s thất bại: %v", seed.Code, err)
		}
		if role.ID == 0 {
			role = model.Role{Code: seed.Code, Name: seed.Name, Description: seed.Description, IsActive: 1}
			if err := db.Create(&role).Error; err != nil {
				log.Fatalf("❌ [RBAC] tạo role %s thất bại: %v", seed.Code, err)
			}
		}

		for _, permissionCode := range seed.Permissions {
			var permission model.Permission
			if err := db.Where("code = ?", permissionCode).First(&permission).Error; err != nil {
				log.Fatalf("❌ [RBAC] không tìm thấy permission %s: %v", permissionCode, err)
			}

			var count int64
			db.Model(&model.RolePermission{}).
				Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).
				Count(&count)
			if count > 0 {
				continue
			}
			if err := db.Create(&model.RolePermission{RoleID: role.ID, PermissionID: permission.ID}).Error; err != nil {
				log.Fatalf("❌ [RBAC] gán %s cho %s thất bại: %v", permissionCode, seed.Code, err)
			}
			createdLinks++
		}
	}

	log.Printf("✅ [RBAC] %d permission, %d role, %d liên kết quyền mới được tạo",
		len(permissionSeeds), len(roleSeeds), createdLinks)
}

// migrateLegacyUserRole chuyển dữ liệu từ cột users.role cũ sang bảng nối user_roles,
// sau đó XOÁ cột cũ vì từ nay role là quan hệ nhiều-nhiều.
func migrateLegacyUserRole(db *gorm.DB) {
	var columnCount int64
	db.Raw(`SELECT COUNT(*) FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name = 'users' AND column_name = 'role'`).Scan(&columnCount)
	if columnCount == 0 {
		return
	}

	// Backfill: mỗi user cũ được gán đúng role đang có
	result := db.Exec(`INSERT IGNORE INTO user_roles (user_id, role_id, created_at)
		SELECT u.id, r.id, NOW()
		FROM users u
		JOIN roles r ON r.code = u.role
		WHERE u.role IS NOT NULL AND u.role <> ''`)
	if result.Error != nil {
		log.Fatalf("❌ [RBAC] migrate users.role → user_roles thất bại: %v", result.Error)
	}

	if err := db.Migrator().DropColumn("users", "role"); err != nil {
		log.Printf("⚠️ [RBAC] không xoá được cột users.role: %v", err)
		return
	}

	log.Printf("✅ [RBAC] đã chuyển %d user từ users.role sang user_roles và xoá cột cũ", result.RowsAffected)
}
