package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitManagerRoutes - Khởi tạo các API route dành riêng cho Manager sử dụng Wire
func InitManagerRoutes(Router fiber.Router) {
	managerHandler, err := wire.InitializeManagerPostHandler()
	if err != nil {
		panic(err)
	}

	managerGroup := Router.Group("/manager")
	{
		// Áp dụng AuthMiddleware yêu cầu đăng nhập đối với mọi tác vụ của Manager
		authGroup := managerGroup.Group("/", middleware.AuthMiddleware)
		{
			// Lấy danh sách bài viết phân trang và lọc của manager
			authGroup.Get("/posts", managerHandler.GetManagerPostsList)
			authGroup.Post("/create-post", managerHandler.CreatePost)
			authGroup.Put("/update-post/:id", managerHandler.UpdatePost)
			// Xóa bài viết chính chủ của manager
			authGroup.Delete("/posts/:id", managerHandler.DeleteManagerPost)
			// Tạo dự án mới (tác nhân quản lý dự án)
			authGroup.Post("/create-project", managerHandler.CreateProject)
			// Danh sách dự án
			authGroup.Get("/projects", managerHandler.ListProjects)
			// Chi tiết dự án (điền form chỉnh sửa)
			authGroup.Get("/projects/:id", managerHandler.GetProjectDetail)
			// Cập nhật dự án
			authGroup.Put("/update-project/:id", managerHandler.UpdateProject)
		}
	}
}
