package routers

import (
	"real_estate_be/internal/middleware"
	"real_estate_be/internal/wire"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

// InitUploadRoutes — upload ảnh/video lên storage. Tất cả đều CẦN ĐĂNG NHẬP.
// (Không truyền middleware vào Router.Group — xem ghi chú ở real_estate_routes.go)
func InitUploadRoutes(Router fiber.Router, s3Client *s3.Client) {
	uploadController, err := wire.InitializeUploadHandler(s3Client)
	if err != nil {
		panic(err)
	}

	uploadRouter := Router.Group("/upload")
	uploadRouter.Post("/image", middleware.AuthMiddleware, uploadController.UploadImage)
	uploadRouter.Post("/presign", middleware.AuthMiddleware, uploadController.Presign)
	uploadRouter.Post("/confirm", middleware.AuthMiddleware, uploadController.Confirm)
}
