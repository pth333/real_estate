package controller

import (
	"strconv"

	"real_estate_be/internal/dto"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// AdminDepositHandler — Admin panel: quản lý tranh chấp, tổng quan escrow, toàn bộ đơn đặt cọc.
type AdminDepositHandler struct {
	service usecase.IDepositService
}

func NewAdminDepositHandler(service usecase.IDepositService) *AdminDepositHandler {
	return &AdminDepositHandler{service: service}
}

// ListDeposits — tất cả đơn đặt cọc, lọc theo trạng thái.
func (h *AdminDepositHandler) ListDeposits(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	items, total, err := h.service.ListAllDeposits(c.Query("status", ""), page, size)
	if err != nil {
		return response.InternalServerError(c, "Lấy danh sách đặt cọc thất bại", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "", items, fiber.Map{"total": total, "page": page, "size": size})
}

// DecidePurchase — admin duyệt/từ chối tài liệu mua nhà.
// Duyệt thì trừ 1 căn tồn kho dự án và tất toán hoàn 100% cho khách.
func (h *AdminDepositHandler) DecidePurchase(c *fiber.Ctx) error {
	adminID, _ := currentActor(c)
	depositID, err := parseIDParam(c)
	if err != nil || depositID == 0 {
		return response.BadRequest(c, "ID đơn đặt cọc không hợp lệ", nil)
	}

	var req dto.PurchaseDecisionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	if _, err := h.service.DecidePurchase(depositID, adminID, req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	if req.Approved {
		return response.OK(c, "Đã duyệt tài liệu mua nhà, hoàn 100% tiền cọc và trừ 1 căn của dự án")
	}
	return response.OK(c, "Đã từ chối tài liệu, đơn chuyển sang tranh chấp để xử lý tiếp")
}

// EscrowSummary — số dư escrow, tiền đang giữ, tiền đã release.
func (h *AdminDepositHandler) EscrowSummary(c *fiber.Ctx) error {
	summary, err := h.service.GetEscrowSummary()
	if err != nil {
		return response.InternalServerError(c, "Lấy số liệu escrow thất bại", err.Error())
	}
	return response.OK(c, summary)
}

// ListDisputes — danh sách tranh chấp cần xử lý.
func (h *AdminDepositHandler) ListDisputes(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	items, total, err := h.service.ListDisputes(c.Query("status", ""), page, size)
	if err != nil {
		return response.InternalServerError(c, "Lấy danh sách tranh chấp thất bại", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "", items, fiber.Map{"total": total, "page": page, "size": size})
}

// GetDispute — chi tiết 1 tranh chấp kèm bằng chứng 2 bên.
func (h *AdminDepositHandler) GetDispute(c *fiber.Ctx) error {
	disputeID, err := parseIDParam(c)
	if err != nil || disputeID == 0 {
		return response.BadRequest(c, "ID tranh chấp không hợp lệ", nil)
	}

	result, err := h.service.GetDisputeDetail(disputeID)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// ResolveDispute — admin ra quyết định và release tiền khỏi escrow.
func (h *AdminDepositHandler) ResolveDispute(c *fiber.Ctx) error {
	adminID, _ := currentActor(c)
	disputeID, err := parseIDParam(c)
	if err != nil || disputeID == 0 {
		return response.BadRequest(c, "ID tranh chấp không hợp lệ", nil)
	}

	var req dto.ResolveDisputeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	result, err := h.service.ResolveDispute(disputeID, adminID, req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}
