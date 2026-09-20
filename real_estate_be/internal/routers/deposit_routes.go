package routers

import (
	"real_estate_be/internal/middleware"
	model "real_estate_be/internal/models"
	"real_estate_be/internal/wire"

	"github.com/gofiber/fiber/v2"
)

// InitDepositRoutes — API luồng đặt cọc escrow:
// khách đặt cọc & check-in, môi giới xác nhận & báo kết quả, admin xử lý tranh chấp.
//
// ⚠️ LƯU Ý QUAN TRỌNG (Fiber v2):
// `Group.Group(prefix, handlers...)` sẽ APPEND handlers vào slice Handlers của group CHA.
// Vì vậy KHÔNG được tạo 2 group có middleware từ cùng 1 group cha, và cũng không được
// dùng group cha để đăng ký route nghiệp vụ sau khi đã tạo group con có middleware —
// nếu không, group tạo sau sẽ thừa hưởng middleware của group tạo trước (dẫn tới 403 oan).
// Ở file này mọi middleware được gắn TRỰC TIẾP trên từng route để tránh hẳn vấn đề đó.
func InitDepositRoutes(Router fiber.Router) {
	depositHandler, err := wire.InitializeDepositHandler()
	if err != nil {
		panic(err)
	}
	adminHandler, err := wire.InitializeAdminDepositHandler()
	if err != nil {
		panic(err)
	}

	depositGroup := Router.Group("/deposits")

	// ── Cổng thanh toán gọi vào (KHÔNG cần token: cổng gọi server-to-server) ──
	depositGroup.Get("/payment/callback", depositHandler.PaymentCallback)
	depositGroup.Get("/payment/ipn", depositHandler.PaymentIPN)
	depositGroup.Post("/payment/ipn", depositHandler.PaymentIPN)

	// ── Khách hàng ──
	depositGroup.Get("/booking-options",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionDepositCreate),
		depositHandler.GetBookingOptions)
	depositGroup.Post("/",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionDepositCreate),
		depositHandler.CreateDeposit)
	depositGroup.Get("/my",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionDepositViewOwn),
		depositHandler.ListMyDeposits)
	depositGroup.Post("/:id/checkin",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionDepositCheckin),
		depositHandler.Checkin)
	depositGroup.Post("/:id/rating",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionDepositRate),
		depositHandler.RateBroker)

	// ── Khách hoặc môi giới của đơn ──
	depositGroup.Get("/:id",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionDepositView),
		depositHandler.GetDeposit)
	depositGroup.Post("/:id/report",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionDepositReport),
		depositHandler.SubmitReport)
	depositGroup.Post("/:id/dispute",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionDepositDispute),
		depositHandler.OpenDispute)

	// ── Môi giới ──
	brokerGroup := Router.Group("/broker/deposits")
	brokerGroup.Get("/",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionBrokerDepositList),
		depositHandler.ListBrokerDeposits)
	brokerGroup.Post("/:id/confirm",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionBrokerDepositConfirm),
		depositHandler.ConfirmDeposit)
	brokerGroup.Post("/:id/reject",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionBrokerDepositReject),
		depositHandler.RejectDeposit)
	brokerGroup.Post("/:id/otp",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionBrokerDepositOtp),
		depositHandler.GenerateOTP)

	// ── Bằng chứng tranh chấp (khách hoặc môi giới) ──
	disputeGroup := Router.Group("/disputes")
	disputeGroup.Post("/:id/evidence",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionDisputeEvidence),
		depositHandler.AddEvidence)

	// ── Admin panel ──
	adminGroup := Router.Group("/admin")
	adminGroup.Get("/deposits",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminDepositList),
		adminHandler.ListDeposits)
	// Chi tiết 1 đơn (route riêng vì /deposits/:id không mở cho admin)
	adminGroup.Get("/deposits/:id",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminDepositView),
		adminHandler.GetDeposit)
	// Duyệt / từ chối tài liệu mua nhà (trừ tồn kho dự án khi duyệt)
	adminGroup.Post("/deposits/:id/purchase-decision",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminDepositApprove),
		adminHandler.DecidePurchase)
	adminGroup.Get("/escrow",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminEscrowView),
		adminHandler.EscrowSummary)
	adminGroup.Get("/disputes",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminDisputeList),
		adminHandler.ListDisputes)
	adminGroup.Get("/disputes/:id",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminDisputeView),
		adminHandler.GetDispute)
	adminGroup.Post("/disputes/:id/resolve",
		middleware.AuthMiddleware, middleware.RequirePermission(model.PermissionAdminDisputeResolve),
		adminHandler.ResolveDispute)
}
