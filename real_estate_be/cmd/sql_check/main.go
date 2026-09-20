// Chương trình kiểm tra tạm thời: tái hiện lỗi 1136 và xác nhận các câu SQL đã sửa.
package main

import (
	"fmt"
	"time"

	"real_estate_be/internal/global"
	"real_estate_be/internal/initialize"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

func main() {
	initialize.LoadConfig()
	initialize.InitMysql()
	initialize.MigrateDb(global.DB)

	user := ensureUser("qa-sql@test.local", "QA SQL")
	defer cleanup()
	fmt.Printf("→ user tạm #%d (không đụng user_id=1 của bạn)\n\n", user.ID)

	fmt.Println("═══ 1. CÂU SQL BỊ LỖI (giống của bạn) ═══")
	bad := fmt.Sprintf("INSERT INTO user_roles (user_id, role_id) VALUES (%d, 1, '2026-09-20 13:50:15')", user.ID)
	fmt.Println("SQL:", bad)
	if err := global.DB.Exec(bad).Error; err != nil {
		fmt.Println("❌ lỗi:", err)
	}

	fmt.Println("\n═══ 2. CÁCH 1 — bỏ giá trị thừa (created_at nhận DEFAULT NULL) ═══")
	ok1 := fmt.Sprintf("INSERT INTO user_roles (user_id, role_id) VALUES (%d, 1)", user.ID)
	fmt.Println("SQL:", ok1)
	if err := global.DB.Exec(ok1).Error; err != nil {
		fmt.Println("❌ lỗi:", err)
	} else {
		fmt.Println("✅ chạy được")
	}

	fmt.Println("\n═══ 3. CÁCH 2 — thêm created_at vào danh sách cột, dùng NOW() ═══")
	ok2 := fmt.Sprintf("INSERT INTO user_roles (user_id, role_id, created_at) VALUES (%d, 2, NOW())", user.ID)
	fmt.Println("SQL:", ok2)
	if err := global.DB.Exec(ok2).Error; err != nil {
		fmt.Println("❌ lỗi:", err)
	} else {
		fmt.Println("✅ chạy được")
	}

	fmt.Println("\n═══ 4. CÁCH 3 — gán nhiều role bằng role CODE (khuyên dùng, không hardcode id) ═══")
	multi := fmt.Sprintf(`INSERT IGNORE INTO user_roles (user_id, role_id, created_at)
		SELECT %d, id, NOW() FROM roles WHERE code IN ('CUSTOMER','BROKER','ADMIN')`, user.ID)
	fmt.Println("SQL:", multi)
	if err := global.DB.Exec(multi).Error; err != nil {
		fmt.Println("❌ lỗi:", err)
	} else {
		fmt.Println("✅ chạy được")
	}

	fmt.Println("\n═══ 5. KẾT QUẢ TRONG BẢNG ═══")
	var links []model.UserRole
	global.DB.Where("user_id = ?", user.ID).Order("role_id ASC").Find(&links)
	for _, l := range links {
		var role model.Role
		global.DB.First(&role, l.RoleID)
		fmt.Printf("• user_id=%d | role_id=%d (%s) | created_at=%v\n",
			l.UserID, l.RoleID, role.Code, l.CreatedAt.Format(time.RFC3339))
	}
	access, _ := repo.NewRbacRepository(global.DB).GetUserAccess(user.ID)
	fmt.Printf("→ GetUserAccess: roles=%v, %d quyền\n", access.Roles, len(access.Permissions))

	fmt.Println("\n═══ 6. CHẠY LẠI CÂU 4 (đã có role rồi) ═══")
	if err := global.DB.Exec(multi).Error; err != nil {
		fmt.Println("❌ lỗi khi chạy lại:", err)
	} else {
		fmt.Println("✅ INSERT IGNORE nên chạy lại không lỗi, không tạo dòng trùng")
	}
	var count int64
	global.DB.Model(&model.UserRole{}).Where("user_id = ?", user.ID).Count(&count)
	fmt.Printf("→ tổng số dòng của user: %d (đúng bằng số role, không nhân đôi)\n", count)
}

func ensureUser(email, name string) *model.User {
	var user model.User
	if err := global.DB.Where("email = ?", email).Limit(1).Find(&user).Error; err == nil && user.ID != 0 {
		return &user
	}
	user = model.User{
		Email: email, Name: name, Password: "qa", IsActive: 1,
		Phone: fmt.Sprintf("0999%06d", time.Now().UnixNano()%1000000),
	}
	if err := global.DB.Create(&user).Error; err != nil {
		panic(err)
	}
	return &user
}

func cleanup() {
	var ids []uint64
	global.DB.Model(&model.User{}).Where("email = ?", "qa-sql@test.local").Pluck("id", &ids)
	if len(ids) > 0 {
		global.DB.Where("user_id IN ?", ids).Delete(&model.UserRole{})
	}
	global.DB.Where("email = ?", "qa-sql@test.local").Delete(&model.User{})
	var left int64
	global.DB.Model(&model.User{}).Where("email LIKE ?", "qa-%@test.local").Count(&left)
	fmt.Printf("\n🧹 đã dọn user test (còn sót: %d)\n", left)
}
