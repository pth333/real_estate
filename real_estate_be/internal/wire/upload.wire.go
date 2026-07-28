package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitializeUploadHandler(s3Client *s3.Client) (*controller.UploadHandler, error) {
	db := providerDB()
	imageRepo := repo.NewImageRepository(db)
	uploadService := usecase.NewUploadService(imageRepo, s3Client)
	uploadHandler := controller.NewUploadHandler(uploadService)
	return uploadHandler, nil
}
