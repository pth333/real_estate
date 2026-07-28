package controller

import (
	"real_estate_be/internal/dto"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	service usecase.UploadServiceInterface
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
