package usecase

import (
	"fmt"
	"log"
	"time"

	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
)

// ══════════════════════════════════════════════════════════
// Cron jobs (plan mục 6)
// ══════════════════════════════════════════════════════════

// RunScheduledTasks chạy toàn bộ job định kỳ. Được gọi từ scheduler mỗi chu kỳ.
func (s *depositService) RunScheduledTasks() {
	now := time.Now()
	s.cancelUnpaidDeposits(now)
	s.autoRejectExpiredDeposits(now)
	s.openReportWindows(now)
	s.sendViewingReminders(now)
	s.escalateOverdueReports(now)
}

// openReportWindows — fallback plan mục 2.3: sau viewing_start + ân hạn mà chưa check-in
// thì mở cửa sổ 24h để 2 bên tự báo cáo kết quả.
func (s *depositService) openReportWindows(now time.Time) {
	items, err := s.depositRepo.ListConfirmedWithoutDeadline()
	if err != nil {
		log.Printf("⚠️ [DepositCron] không lấy được đơn chưa mở cửa sổ báo cáo: %v", err)
		return
	}

	grace := time.Duration(s.cfg.CheckinGraceHours) * time.Hour
	window := time.Duration(s.cfg.ReportWindowHours) * time.Hour

	for i := range items {
		deposit := &items[i]
		startAt, _, err := viewingRange(deposit)
		if err != nil {
			continue
		}
		// Chưa tới thời điểm mở cửa sổ báo cáo
		if now.Before(startAt.Add(grace)) {
			continue
		}

		deadline := now.Add(window)
		if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
			"report_deadline": deadline,
		}); err != nil {
			log.Printf("⚠️ [DepositCron] mở cửa sổ báo cáo đơn %d thất bại: %v", deposit.ID, err)
			continue
		}

		s.notifyBoth(deposit, "report_window_opened",
			fmt.Sprintf("Xác nhận kết quả buổi xem nhà #%d", deposit.ID),
			fmt.Sprintf("Hệ thống chưa ghi nhận check-in cho buổi xem nhà #%d.\nVui lòng báo cáo kết quả trong %d giờ tới. Quá hạn, phí môi giới sẽ được chuyển cho admin xác minh.",
				deposit.ID, s.cfg.ReportWindowHours),
			fmt.Sprintf("Hệ thống chưa ghi nhận check-in cho buổi xem nhà #%d.\nVui lòng báo cáo kết quả trong %d giờ tới. Quá hạn, phí môi giới sẽ được chuyển cho admin xác minh.",
				deposit.ID, s.cfg.ReportWindowHours),
		)
	}
}

// cancelUnpaidDeposits — huỷ đơn quá hạn thanh toán để trả lại khung giờ xem nhà.
func (s *depositService) cancelUnpaidDeposits(now time.Time) {
	deadline := now.Add(-time.Duration(s.cfg.PaymentTimeoutMinutes) * time.Minute)
	items, err := s.depositRepo.ListUnpaidBefore(deadline)
	if err != nil {
		log.Printf("⚠️ [DepositCron] không lấy được đơn quá hạn thanh toán: %v", err)
		return
	}

	for i := range items {
		if err := s.depositRepo.UpdateFields(items[i].ID, map[string]interface{}{
			"status":        model.DepositStatusCancelled,
			"reject_reason": "Quá hạn thanh toán, hệ thống tự huỷ",
		}); err != nil {
			log.Printf("⚠️ [DepositCron] huỷ đơn %d thất bại: %v", items[i].ID, err)
		}
	}
	if len(items) > 0 {
		log.Printf("🧹 [DepositCron] đã huỷ %d đơn quá hạn thanh toán", len(items))
	}
}

// autoRejectExpiredDeposits — môi giới không phản hồi trong 24h → từ chối + hoàn 100%.
func (s *depositService) autoRejectExpiredDeposits(now time.Time) {
	deadline := now.Add(-time.Duration(s.cfg.BrokerConfirmHours) * time.Hour)
	items, err := s.depositRepo.ListPendingBefore(deadline)
	if err != nil {
		log.Printf("⚠️ [DepositCron] không lấy được đơn quá hạn xác nhận: %v", err)
		return
	}

	for i := range items {
		deposit := &items[i]
		if err := s.settleWithNotify(deposit, model.DepositStatusBrokerRejected,
			"Môi giới không phản hồi trong thời hạn xác nhận"); err != nil {
			log.Printf("⚠️ [DepositCron] tự động từ chối đơn %d thất bại: %v", deposit.ID, err)
		}
	}
	if len(items) > 0 {
		log.Printf("⏰ [DepositCron] đã tự động từ chối %d đơn quá hạn môi giới xác nhận", len(items))
	}
}

// sendViewingReminders — gửi mail nhắc trước 24h cho cả khách và môi giới.
func (s *depositService) sendViewingReminders(now time.Time) {
	tomorrow := now.AddDate(0, 0, 1)
	items, err := s.depositRepo.ListReminderDue(tomorrow)
	if err != nil {
		log.Printf("⚠️ [DepositCron] không lấy được đơn cần nhắc lịch: %v", err)
		return
	}

	for i := range items {
		deposit := &items[i]
		contactName, contactPhone := contactOf(deposit)
		s.notifyBoth(deposit, "viewing_reminder",
			fmt.Sprintf("Nhắc lịch xem nhà ngày mai - đơn #%d", deposit.ID),
			fmt.Sprintf("Bạn có lịch xem nhà vào ngày mai.\nĐịa chỉ: %s\nThời gian: %s - %s\nMôi giới: %s - %s",
				estateAddress(deposit), deposit.ViewingStart, deposit.ViewingEnd, deposit.Broker.Name, deposit.Broker.Phone),
			fmt.Sprintf("Bạn có lịch dẫn khách xem nhà vào ngày mai.\nĐịa chỉ: %s\nThời gian: %s - %s\nKhách: %s - %s",
				estateAddress(deposit), deposit.ViewingStart, deposit.ViewingEnd, contactName, contactPhone),
		)

		if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
			"reminder_sent_at": now,
		}); err != nil {
			log.Printf("⚠️ [DepositCron] đánh dấu đã nhắc đơn %d thất bại: %v", deposit.ID, err)
		}
	}
}

// escalateOverdueReports — quá hạn báo cáo mà chưa có kết quả → xử lý theo log check-in.
func (s *depositService) escalateOverdueReports(now time.Time) {
	// Chưa check-in được và đã hết cửa sổ 24h báo cáo.
	// Chỉ còn phí môi giới nên không thể treo tiền chờ admin xử tay: áp suy đoán theo
	// LOG CHECK-IN (vị trí/OTP) — bên có căn cứ thắng khi bên kia im lặng.
	// Admin chỉ nhận ca 2 bên đều có log mà mâu thuẫn.
	notCheckedIn, err := s.depositRepo.ListReportDeadlinePassed(now)
	if err != nil {
		log.Printf("⚠️ [DepositCron] không lấy được đơn quá hạn báo cáo: %v", err)
	}
	for i := range notCheckedIn {
		deposit := &notCheckedIn[i]
		if err := s.settleOverdueByCheckinLog(deposit); err != nil {
			log.Printf("⚠️ [DepositCron] xử lý đơn quá hạn %d thất bại: %v", deposit.ID, err)
		}
	}

	// Đã check-in xong và hết thời gian giữ phí (4 ngày sau buổi xem).
	// Phần mua/không mua đã bỏ: phí thuộc môi giới, TRỪ KHI khách đặt cọc mua BĐS trong thời gian giữ.
	checkedIn, err := s.depositRepo.ListCheckedInBefore(now)
	if err != nil {
		log.Printf("⚠️ [DepositCron] không lấy được đơn hết thời gian giữ phí: %v", err)
	}
	for i := range checkedIn {
		deposit := &checkedIn[i]

		// Đơn cũ còn báo cáo "đã mua nhà": giữ nguyên luồng duyệt tài liệu mua như trước
		if deposit.BrokerReport == model.ReportBought || deposit.CustomerReport == model.ReportBought {
			s.settleByReport(deposit, model.ReportBought, "Đơn cũ: một bên khai khách đã mua nhà")
			continue
		}

		if err := s.releaseFeeAfterHold(deposit); err != nil {
			log.Printf("⚠️ [DepositCron] tất toán đơn %d sau thời gian giữ phí thất bại: %v", deposit.ID, err)
		}
	}
}

// releaseFeeAfterHold — hết thời gian giữ phí (4 ngày sau buổi xem):
// - Khách ĐÃ đặt cọc mua bất động sản → hoàn 100% phí cho khách.
// - Khách KHÔNG đặt cọc mua → phí thuộc về môi giới.
func (s *depositService) releaseFeeAfterHold(deposit *model.Deposit) error {
	if deposit.PurchaseDepositAt != nil {
		return s.settleWithNotify(deposit, model.DepositStatusRefunded,
			"Khách đã đặt cọc mua bất động sản trong thời gian giữ phí → hoàn 100% phí")
	}
	return s.settleWithNotify(deposit, model.DepositStatusCompleted,
		"Hết thời gian giữ phí mà khách không đặt cọc mua → phí thuộc về môi giới")
}

// settleOverdueByCheckinLog — hết hạn báo cáo mà chưa check-in: quyết định theo LOG check-in.
//
// Nguyên tắc: bên có căn cứ (vị trí/OTP) thắng khi bên kia không có gì; chỉ khi cả 2 đều
// có log mà vẫn mâu thuẫn mới cần admin. Không ai có log ⇒ hoàn 100% cho khách, vì môi giới
// là bên nắm công cụ sinh OTP để chứng minh mình đã tới mà đã không dùng.
func (s *depositService) settleOverdueByCheckinLog(deposit *model.Deposit) error {
	brokerPresent := deposit.BrokerCheckin != nil && *deposit.BrokerCheckin
	customerPresent := deposit.CustomerCheckin != nil && *deposit.CustomerCheckin

	switch {
	case brokerPresent && !customerPresent:
		note := "Quá hạn: chỉ môi giới có log check-in → môi giới nhận phí"
		if deposit.CustomerReport == model.ReportAttended {
			note = "Quá hạn: khách khai có mặt nhưng không có log vị trí/OTP → môi giới nhận phí"
		}
		return s.settleWithNotify(deposit, model.DepositStatusNoShowCustomer, note)

	case !brokerPresent && customerPresent:
		note := "Quá hạn: chỉ khách có log check-in → hoàn 100% phí cho khách"
		if deposit.BrokerReport == model.ReportAttended {
			note = "Quá hạn: môi giới khai có mặt nhưng không check-in tại chỗ → hoàn 100% phí cho khách"
		}
		return s.settleWithNotify(deposit, model.DepositStatusNoShowBroker, note)

	case brokerPresent && customerPresent:
		// Cả 2 đều có log nhưng vẫn chưa tất toán được (log/vị trí lệch nhau) → admin xác minh
		return s.openSystemDispute(deposit,
			"Hai bên đều có log check-in nhưng chưa thống nhất được kết quả buổi xem")

	default:
		return s.settleWithNotify(deposit, model.DepositStatusRefunded,
			"Không bên nào xác nhận buổi xem → hoàn 100% phí cho khách")
	}
}

// settleByReport tất toán theo giá trị báo cáo mua/không mua đã nhận.
// Riêng báo cáo BOUGHT vẫn phải qua admin duyệt tài liệu trước khi trừ tồn kho.
func (s *depositService) settleByReport(deposit *model.Deposit, report, note string) {
	if report == model.ReportBought {
		if err := s.awaitPurchaseApproval(deposit); err != nil {
			log.Printf("⚠️ [DepositCron] chuyển đơn %d sang chờ duyệt tài liệu thất bại: %v", deposit.ID, err)
		}
		return
	}
	if err := s.settleWithNotify(deposit, model.DepositStatusVisitedNotBuy, note); err != nil {
		log.Printf("⚠️ [DepositCron] tất toán đơn %d theo báo cáo thất bại: %v", deposit.ID, err)
	}
}

// ══════════════════════════════════════════════════════════
// Admin: tổng quan quỹ escrow (plan mục 8)
// ══════════════════════════════════════════════════════════

// GetEscrowSummary tính số dư escrow từ sổ giao dịch (transactions).
func (s *depositService) GetEscrowSummary() (*dto.EscrowSummaryResponse, error) {
	sums, err := s.transactionRepo.SumByType()
	if err != nil {
		return nil, err
	}

	refunded := sums[model.TransactionRefundFull] + sums[model.TransactionRefundPartial]
	transferred := sums[model.TransactionTransferToBroker]
	penalty := sums[model.TransactionPenaltyBroker]
	// Phí còn đang giữ = đã thu - đã hoàn - đã chuyển cho môi giới
	holding := sums[model.TransactionDeposit] - refunded - transferred

	counts, err := s.depositRepo.CountByStatus()
	if err != nil {
		return nil, err
	}
	// Điền đủ trạng thái = 0 để FE dựng bảng không bị thiếu dòng
	for _, status := range allDepositStatuses() {
		if _, ok := counts[status]; !ok {
			counts[status] = 0
		}
	}

	openDisputes, err := s.disputeRepo.CountOpen()
	if err != nil {
		return nil, err
	}

	return &dto.EscrowSummaryResponse{
		HoldingAmount:     round2(holding),
		RefundedAmount:    round2(refunded),
		TransferredAmount: round2(transferred),
		PenaltyAmount:     round2(penalty),
		StatusCounts:      counts,
		OpenDisputes:      openDisputes,
	}, nil
}

func allDepositStatuses() []string {
	return []string{
		model.DepositStatusAwaitingPayment,
		model.DepositStatusPending,
		model.DepositStatusBrokerRejected,
		model.DepositStatusBrokerConfirmed,
		model.DepositStatusCheckedIn,
		model.DepositStatusPendingPurchase,
		model.DepositStatusVisitedBought,
		model.DepositStatusVisitedNotBuy,
		model.DepositStatusNoShowCustomer,
		model.DepositStatusNoShowBroker,
		model.DepositStatusDispute,
		model.DepositStatusRefunded,
		model.DepositStatusCompleted,
		model.DepositStatusCancelled,
	}
}
