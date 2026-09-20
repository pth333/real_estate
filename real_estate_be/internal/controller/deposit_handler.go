package controller

import (
	"strconv"

	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// DepositHandler — API cho khách hàng, môi giới và cổng thanh toán.
// Handlers chỉ nhận/trả dữ liệu, toàn bộ nghiệp vụ nằm ở DepositService.
type DepositHandler struct {
	service usecase.IDepositService
}

func NewDepositHandler(service usecase.IDepositService) *DepositHandler {
	return &DepositHandler{service: service}
}

// currentActor đọc user_id + role do RequireRole nạp vào Locals.
func currentActor(c *fiber.Ctx) (uint64, string) {
	userID, _ := c.Locals("user_id").(uint64)
	role, _ := c.Locals("role").(string)
	return userID, role
}

func parseIDParam(c *fiber.Ctx) (uint64, error) {
	return strconv.ParseUint(c.Params("id"), 10, 64)
}

// ── Khách hàng ──────────────────────────────────────────

// GetBookingOptions — mức cọc + phí môi giới hệ thống đề xuất cho 1 BĐS.
// FE dùng để hiển thị trong form đặt cọc (khách không sửa được số tiền).
func (h *DepositHandler) GetBookingOptions(c *fiber.Ctx) error {
	realEstateID, err := strconv.ParseUint(c.Query("real_estate_id"), 10, 64)
	if err != nil || realEstateID == 0 {
		return response.BadRequest(c, "ID bất động sản không hợp lệ", nil)
	}

	result, err := h.service.GetBookingOptions(realEstateID)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// CreateDeposit — khách đặt cọc xem nhà, trả về URL thanh toán.
func (h *DepositHandler) CreateDeposit(c *fiber.Ctx) error {
	customerID, _ := currentActor(c)

	var req dto.CreateDepositRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	result, err := h.service.CreateDeposit(customerID, req, c.IP())
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Khởi tạo đặt cọc thành công, vui lòng thanh toán", result)
}

// ListMyDeposits — danh sách đơn đặt cọc của khách đang đăng nhập.
func (h *DepositHandler) ListMyDeposits(c *fiber.Ctx) error {
	customerID, _ := currentActor(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	items, total, err := h.service.ListCustomerDeposits(customerID, c.Query("status", ""), page, size)
	if err != nil {
		return response.InternalServerError(c, "Lấy danh sách đặt cọc thất bại", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "", items, fiber.Map{"total": total, "page": page, "size": size})
}

// Checkin — khách nhập OTP do môi giới hiển thị để xác nhận 2 bên có mặt.
func (h *DepositHandler) Checkin(c *fiber.Ctx) error {
	customerID, _ := currentActor(c)
	depositID, err := parseIDParam(c)
	if err != nil || depositID == 0 {
		return response.BadRequest(c, "ID đơn đặt cọc không hợp lệ", nil)
	}

	var req dto.CheckinRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	result, err := h.service.CustomerCheckin(depositID, customerID, req.OTP)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// RateBroker — khách đánh giá môi giới sau buổi xem.
func (h *DepositHandler) RateBroker(c *fiber.Ctx) error {
	customerID, _ := currentActor(c)
	depositID, err := parseIDParam(c)
	if err != nil || depositID == 0 {
		return response.BadRequest(c, "ID đơn đặt cọc không hợp lệ", nil)
	}

	var req dto.RateBrokerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	if err := h.service.RateBroker(depositID, customerID, req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, "Cảm ơn bạn đã đánh giá môi giới")
}

// ── Môi giới ────────────────────────────────────────────

// ListBrokerDeposits — danh sách đơn đặt cọc môi giới đang phụ trách.
func (h *DepositHandler) ListBrokerDeposits(c *fiber.Ctx) error {
	brokerID, _ := currentActor(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	items, total, err := h.service.ListBrokerDeposits(brokerID, c.Query("status", ""), page, size)
	if err != nil {
		return response.InternalServerError(c, "Lấy danh sách đặt cọc thất bại", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "", items, fiber.Map{"total": total, "page": page, "size": size})
}

// ConfirmDeposit — môi giới xác nhận lịch xem nhà.
func (h *DepositHandler) ConfirmDeposit(c *fiber.Ctx) error {
	brokerID, _ := currentActor(c)
	depositID, err := parseIDParam(c)
	if err != nil || depositID == 0 {
		return response.BadRequest(c, "ID đơn đặt cọc không hợp lệ", nil)
	}

	result, err := h.service.ConfirmDeposit(depositID, brokerID)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// RejectDeposit — môi giới từ chối lịch kèm lý do, hệ thống hoàn 100% tiền cọc.
func (h *DepositHandler) RejectDeposit(c *fiber.Ctx) error {
	brokerID, _ := currentActor(c)
	depositID, err := parseIDParam(c)
	if err != nil || depositID == 0 {
		return response.BadRequest(c, "ID đơn đặt cọc không hợp lệ", nil)
	}

	var req dto.RejectDepositRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	result, err := h.service.RejectDeposit(depositID, brokerID, req.Reason)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// GenerateOTP — môi giới sinh mã OTP check-in tại chỗ.
func (h *DepositHandler) GenerateOTP(c *fiber.Ctx) error {
	brokerID, _ := currentActor(c)
	depositID, err := parseIDParam(c)
	if err != nil || depositID == 0 {
		return response.BadRequest(c, "ID đơn đặt cọc không hợp lệ", nil)
	}

	result, err := h.service.GenerateCheckinOTP(depositID, brokerID)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// ── Dùng chung 2 bên ────────────────────────────────────

// GetDeposit — chi tiết 1 đơn đặt cọc.
func (h *DepositHandler) GetDeposit(c *fiber.Ctx) error {
	userID, role := currentActor(c)
	depositID, err := parseIDParam(c)
	if err != nil || depositID == 0 {
		return response.BadRequest(c, "ID đơn đặt cọc không hợp lệ", nil)
	}

	result, err := h.service.GetDepositDetail(depositID, userID, role)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// SubmitReport — 1 bên báo cáo kết quả buổi xem nhà.
func (h *DepositHandler) SubmitReport(c *fiber.Ctx) error {
	userID, role := currentActor(c)
	depositID, err := parseIDParam(c)
	if err != nil || depositID == 0 {
		return response.BadRequest(c, "ID đơn đặt cọc không hợp lệ", nil)
	}

	var req dto.ReportResultRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	result, err := h.service.SubmitReport(depositID, userID, role, req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// OpenDispute — khách hoặc môi giới mở tranh chấp.
func (h *DepositHandler) OpenDispute(c *fiber.Ctx) error {
	userID, role := currentActor(c)
	depositID, err := parseIDParam(c)
	if err != nil || depositID == 0 {
		return response.BadRequest(c, "ID đơn đặt cọc không hợp lệ", nil)
	}
	if role == model.RoleAdmin {
		return response.BadRequest(c, "admin không mở tranh chấp thay các bên", nil)
	}

	var req dto.CreateDisputeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	result, err := h.service.OpenDispute(depositID, userID, role, req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.Created(c, "Đã mở tranh chấp, tiền cọc được tạm giữ", result)
}

// AddEvidence — bổ sung bằng chứng cho tranh chấp đang mở.
func (h *DepositHandler) AddEvidence(c *fiber.Ctx) error {
	userID, _ := currentActor(c)
	disputeID, err := parseIDParam(c)
	if err != nil || disputeID == 0 {
		return response.BadRequest(c, "ID tranh chấp không hợp lệ", nil)
	}

	var req dto.AddEvidenceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Dữ liệu gửi lên không hợp lệ", err.Error())
	}

	result, err := h.service.AddDisputeEvidence(disputeID, userID, req)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// ── Cổng thanh toán ─────────────────────────────────────

// PaymentCallback — return URL: cổng redirect khách về đây sau khi thanh toán.
func (h *DepositHandler) PaymentCallback(c *fiber.Ctx) error {
	result, err := h.service.HandlePaymentCallback(queryParams(c))
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, result)
}

// PaymentIPN — IPN: cổng gọi server-to-server, trả đúng format VNPay yêu cầu.
func (h *DepositHandler) PaymentIPN(c *fiber.Ctx) error {
	result, err := h.service.HandlePaymentCallback(queryParams(c))
	if err != nil {
		return c.JSON(fiber.Map{"RspCode": "99", "Message": err.Error()})
	}
	if !result.Success {
		return c.JSON(fiber.Map{"RspCode": "02", "Message": "Giao dịch không thành công"})
	}
	return c.JSON(fiber.Map{"RspCode": "00", "Message": "Confirm Success"})
}

// queryParams gom toàn bộ query string thành map để service xác thực chữ ký.
func queryParams(c *fiber.Ctx) map[string]string {
	params := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		params[string(key)] = string(value)
	})
	return params
}
