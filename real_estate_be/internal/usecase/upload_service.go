package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"real_estate_be/internal/dto"
	"strings"
	"real_estate_be/internal/global"
	"real_estate_be/internal/helpers"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

type UploadServiceInterface interface {
	CreatePresignURL(req dto.PresignRequest) (*dto.PresignResponse, error)
	ConfirmUpload(req dto.ConfirmUploadRequest) (*dto.ConfirmUploadResponse, error)
}

type uploadService struct {
	imageRepo     repo.ImageRepository
	s3Client      *s3.Client
	presignClient *s3.PresignClient
}

func NewUploadService(imageRepo repo.ImageRepository, s3Client *s3.Client) UploadServiceInterface {
	return &uploadService{
		imageRepo:     imageRepo,
		s3Client:      s3Client,
		presignClient: s3.NewPresignClient(s3Client),
	}
}

// generateObjectKey tạo key duy nhất cho file
func generateObjectKey(filename string) string {
	raw := make([]byte, 16)
	rand.Read(raw)
	id := hex.EncodeToString(raw)
	ext := filepath.Ext(filename)
	return fmt.Sprintf("uploads/%s%s", id, ext)
}

func (s *uploadService) CreatePresignURL(req dto.PresignRequest) (*dto.PresignResponse, error) {
	// Validate
	if err := helpers.ValidateFilename(req.Filename); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := helpers.ValidateContentType(req.ContentType); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Sinh object key duy nhất
	key := generateObjectKey(req.Filename)
	bucket := global.Config.R2.Bucket

	// Tạo presigned PUT URL — 15 phút
	ctx := context.Background()
	expiry := 15 * time.Minute

	reqObj := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(req.ContentType),
	}

	presignReq, err := s.presignClient.PresignPutObject(ctx, reqObj, s3.WithPresignExpires(expiry))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Không thể tạo presigned URL: "+err.Error())
	}

	expiresAt := time.Now().Add(expiry).Format(time.RFC3339)

	return &dto.PresignResponse{
		UploadURL: presignReq.URL,
		Key:       key,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *uploadService) ConfirmUpload(req dto.ConfirmUploadRequest) (*dto.ConfirmUploadResponse, error) {
	// Gọi HeadObject lên R2 để kiểm tra file tồn tại
	ctx := context.Background()
	headObj, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(global.Config.R2.Bucket),
		Key:    aws.String(req.Key),
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "File không tồn tại trên R2 hoặc đã hết hạn")
	}

	// Xác định content_type từ kết quả HeadObject
	contentType := ""
	if headObj.ContentType != nil {
		contentType = *headObj.ContentType
	}

	// Xác định filename từ key
	filename := filepath.Base(req.Key)

	// Tạo thumbnail URL cho video
	var thumbnailURL string
	if strings.HasPrefix(contentType, "video/") {
		thumbnailKey := fmt.Sprintf("thumbnails/%s.jpg", strings.TrimSuffix(filepath.Base(req.Key), filepath.Ext(req.Key)))
		thumbnailURL = global.Config.R2.PublicURL + "/" + thumbnailKey
	}

	// Lưu record vào DB
	image := &model.Image{
		Key:          req.Key,
		Filename:     filename,
		FileType:     contentType,
		FileSize:     aws.ToInt64(headObj.ContentLength),
		URL:          global.Config.R2.PublicURL + "/" + req.Key,
		ThumbnailURL: thumbnailURL,
	}

	if err := s.imageRepo.Create(image); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Lưu thông tin ảnh thất bại")
	}

	return &dto.ConfirmUploadResponse{
		ImageID:      image.ID,
		PublicURL:    image.URL,
		Key:          req.Key,
		ThumbnailURL: thumbnailURL,
	}, nil
}
