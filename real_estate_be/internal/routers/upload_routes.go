package routers

import (
	"real_estate_be/internal/wire"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

func InitUploadRoutes(Router fiber.Router, s3Client *s3.Client) {
	uploadController, err := wire.InitializeUploadHandler(s3Client)
	if err != nil {
		panic(err)
	}

	uploadRouter := Router.Group("/upload")
	{
		uploadRouter.Post("/presign", uploadController.Presign)
		uploadRouter.Post("/confirm", uploadController.Confirm)
	}
}
