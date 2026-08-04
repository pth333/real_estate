package controller

import (
	"real_estate_be/internal/dto"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type TrackingHandler struct {
	trackingService usecase.ITrackingService
}

func NewTrackingHandler(trackingService usecase.ITrackingService) *TrackingHandler {
	return &TrackingHandler{trackingService: trackingService}
}

// RecordSearch lưu lịch sử tìm kiếm
func (h *TrackingHandler) RecordSearch(c *fiber.Ctx) error {
	var req dto.TrackingSearchRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := h.trackingService.RecordSearch(req); err != nil {
		return response.InternalServerError(c, "Failed to record search", err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
