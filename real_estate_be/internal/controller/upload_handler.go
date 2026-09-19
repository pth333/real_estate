package controller

import (
	"io"
	"real_estate_be/internal/dto"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	service usecase.UploadServiceInterface
}

func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, "file ảnh là bắt buộc", err.Error())
	}

	source, err := file.Open()
	if err != nil {
		return response.BadRequest(c, "không thể đọc file ảnh", err.Error())
	}
	defer source.Close()

	payload, err := io.ReadAll(source)
	if err != nil {
		return response.BadRequest(c, "không thể đọc nội dung ảnh", err.Error())
	}

	result, err := h.service.UploadImage(
		file.Filename,
		file.Header.Get("Content-Type"),
		payload,
		c.FormValue("kind"),
	)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Upload image successful", result)
}

func NewUploadHandler(service usecase.UploadServiceInterface) *UploadHandler {
	return &UploadHandler{service: service}
}

func (h *UploadHandler) Presign(c *fiber.Ctx) error {
	var req dto.PresignRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if req.Filename == "" || req.ContentType == "" {
		return response.BadRequest(c, "filename và content_type là bắt buộc", nil)
	}

	result, err := h.service.CreatePresignURL(req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, result)
}

func (h *UploadHandler) Confirm(c *fiber.Ctx) error {
	var req dto.ConfirmUploadRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if req.Key == "" {
		return response.BadRequest(c, "key là bắt buộc", nil)
	}

	result, err := h.service.ConfirmUpload(req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Upload confirmed", result)
}
