package controller

import (
	"real_estate_be/internal/dto"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"
	"real_estate_be/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type TrackingHandler struct {
	trackingService usecase.ITrackingService
}

func NewTrackingHandler(trackingService usecase.ITrackingService) *TrackingHandler {
	return &TrackingHandler{trackingService: trackingService}
}

// helper trích xuất UserID từ Authorization Header nếu có (Public tracking)
func (h *TrackingHandler) getUserIDFromHeader(c *fiber.Ctx) uint64 {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0
	}

	tokenStr := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenStr = authHeader[7:]
	}

	claims, err := jwt.ParseAccessToken(tokenStr)
	if err != nil {
		return 0
	}

	return claims.UserID
}

// RecordSearch lưu lịch sử tìm kiếm
func (h *TrackingHandler) RecordSearch(c *fiber.Ctx) error {
	var req dto.TrackingSearchRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Tự động gán UserID từ token nếu có đăng nhập
	if userID := h.getUserIDFromHeader(c); userID > 0 {
		req.UserID = string(rune(userID)) // Hoặc convert phù hợp với kiểu dữ liệu của SearchHistory (string char(10))
	}

	if err := h.trackingService.RecordSearch(req); err != nil {
		return response.InternalServerError(c, "Failed to record search", err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RecordView lưu lịch sử xem chi tiết BĐS
func (h *TrackingHandler) RecordView(c *fiber.Ctx) error {
	var req dto.TrackingViewRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Tự động gán UserID giải mã từ token bảo mật
	req.UserID = h.getUserIDFromHeader(c)

	if err := h.trackingService.RecordView(req); err != nil {
		return response.InternalServerError(c, "Failed to record view", err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// MergeSession sáp nhập session của Guest sang User khi đăng nhập
func (h *TrackingHandler) MergeSession(c *fiber.Ctx) error {
	var req dto.MergeSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if req.SessionID == "" {
		return response.BadRequest(c, "session_id is required", nil)
	}

	// API này phải qua AuthMiddleware nên chắc chắn giải mã được token lấy UserID
	userID := h.getUserIDFromHeader(c)
	if userID == 0 {
		return response.Unauthorized(c, "Unauthorized to merge session", nil)
	}

	if err := h.trackingService.MergeSession(req.SessionID, userID); err != nil {
		return response.InternalServerError(c, "Failed to merge session history", err.Error())
	}

	return response.OK(c, fiber.Map{
		"message": "Session merged successfully",
	})
}
