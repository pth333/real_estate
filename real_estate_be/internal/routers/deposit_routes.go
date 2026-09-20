package routers

import (
	"real_estate_be/internal/middleware"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitDepositRoutes — API luồng đặt cọc escrow:
// khách đặt cọc & check-in, môi giới xác nhận & báo kết quả, admin xử lý tranh chấp.
func InitDepositRoutes(Router fiber.Router) {
	depositHandler, err := wire.InitializeDepositHandler()
	if err != nil {
		panic(err)
	}
	adminHandler, err := wire.InitializeAdminDepositHandler()
	if err != nil {
		panic(err)
	}

	// ── Khách hàng + dùng chung 2 bên ──
	depositGroup := Router.Group("/deposits")
	{
		// Cổng thanh toán gọi vào (public, không cần token)
		depositGroup.Get("/payment/callback", depositHandler.PaymentCallback)
		depositGroup.Get("/payment/ipn", depositHandler.PaymentIPN)
		depositGroup.Post("/payment/ipn", depositHandler.PaymentIPN)

		customerGroup := depositGroup.Group("/", middleware.AuthMiddleware, middleware.RequireRole(model.RoleCustomer))
		{
			// Mức cọc đề xuất theo giá BĐS (khách chỉ xem, không sửa)
			customerGroup.Get("/booking-options", depositHandler.GetBookingOptions)
			customerGroup.Post("/", depositHandler.CreateDeposit)
			customerGroup.Get("/my", depositHandler.ListMyDeposits)
			customerGroup.Post("/:id/checkin", depositHandler.Checkin)
			customerGroup.Post("/:id/rating", depositHandler.RateBroker)
		}

		sharedGroup := depositGroup.Group("/", middleware.AuthMiddleware,
			middleware.RequireRole(model.RoleCustomer, model.RoleBroker))
		{
			sharedGroup.Get("/:id", depositHandler.GetDeposit)
			sharedGroup.Post("/:id/report", depositHandler.SubmitReport)
			sharedGroup.Post("/:id/dispute", depositHandler.OpenDispute)
		}
	}

	// ── Môi giới ──
	brokerGroup := Router.Group("/broker/deposits", middleware.AuthMiddleware, middleware.RequireRole(model.RoleBroker))
	{
		brokerGroup.Get("/", depositHandler.ListBrokerDeposits)
		brokerGroup.Post("/:id/confirm", depositHandler.ConfirmDeposit)
		brokerGroup.Post("/:id/reject", depositHandler.RejectDeposit)
		brokerGroup.Post("/:id/otp", depositHandler.GenerateOTP)
	}

	// ── Bằng chứng tranh chấp (khách hoặc môi giới) ──
	disputeGroup := Router.Group("/disputes", middleware.AuthMiddleware,
		middleware.RequireRole(model.RoleCustomer, model.RoleBroker))
	{
		disputeGroup.Post("/:id/evidence", depositHandler.AddEvidence)
	}

	// ── Admin panel ──
	adminGroup := Router.Group("/admin", middleware.AuthMiddleware, middleware.RequireRole(model.RoleAdmin))
	{
		adminGroup.Get("/deposits", adminHandler.ListDeposits)
		// Duyệt / từ chối tài liệu mua nhà (trừ tồn kho dự án khi duyệt)
		adminGroup.Post("/deposits/:id/purchase-decision", adminHandler.DecidePurchase)
		adminGroup.Get("/escrow", adminHandler.EscrowSummary)
		adminGroup.Get("/disputes", adminHandler.ListDisputes)
		adminGroup.Get("/disputes/:id", adminHandler.GetDispute)
		adminGroup.Post("/disputes/:id/resolve", adminHandler.ResolveDispute)
	}
}
