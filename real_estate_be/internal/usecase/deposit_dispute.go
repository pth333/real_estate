package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"real_estate_be/internal/dto"
	model "real_estate_be/internal/models"
)

// ══════════════════════════════════════════════════════════
// 2.5 Xử lý tranh chấp (Dispute)
// ══════════════════════════════════════════════════════════

// OpenDispute — khách hoặc môi giới mở tranh chấp. Tiền bị freeze tại escrow.
func (s *depositService) OpenDispute(depositID, userID uint64, role string, req dto.CreateDisputeRequest) (*dto.DisputeResponse, error) {
	deposit, err := s.getOwnedDeposit(depositID, userID, role)
	if err != nil {
		return nil, err
	}
	if deposit.Status == model.DepositStatusDispute {
		return nil, errors.New("đơn đặt cọc này đang trong quá trình xử lý tranh chấp")
	}
	if deposit.IsFinished() || deposit.Status == model.DepositStatusAwaitingPayment {
		return nil, errors.New("đơn đặt cọc đã kết thúc, không thể mở tranh chấp")
	}
	if strings.TrimSpace(req.Reason) == "" {
		return nil, errors.New("vui lòng nhập lý do tranh chấp")
	}

	raisedBy := model.DisputeRaisedByCustomer
	if role == model.RoleBroker {
		raisedBy = model.DisputeRaisedByBroker
	}

	deadline := time.Now().Add(time.Duration(s.cfg.DisputeEvidenceHours) * time.Hour)
	dispute := &model.Dispute{
		DepositID:        deposit.ID,
		RaisedBy:         raisedBy,
		Reason:           req.Reason,
		EvidenceURLs:     encodeEvidenceURLs(req.EvidenceURLs),
		Status:           model.DisputeStatusOpen,
		EvidenceDeadline: &deadline,
	}
	if err := s.disputeRepo.Create(dispute); err != nil {
		return nil, err
	}

	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"status": model.DepositStatusDispute,
	}); err != nil {
		return nil, err
	}
	deposit.Status = model.DepositStatusDispute

	s.notifyBoth(deposit, "dispute_opened",
		fmt.Sprintf("Đơn đặt cọc #%d đã bị tạm giữ để xử lý tranh chấp", deposit.ID),
		fmt.Sprintf("Tranh chấp đã được mở cho đơn #%d.\nLý do: %s\nTiền cọc đang bị tạm giữ. Vui lòng upload bằng chứng trong %d giờ tới.",
			deposit.ID, req.Reason, s.cfg.DisputeEvidenceHours),
		fmt.Sprintf("Tranh chấp đã được mở cho đơn #%d.\nLý do: %s\nTiền cọc đang bị tạm giữ. Vui lòng upload bằng chứng trong %d giờ tới.",
			deposit.ID, req.Reason, s.cfg.DisputeEvidenceHours),
	)

	return &dto.DisputeResponse{
		ID:               dispute.ID,
		DepositID:        dispute.DepositID,
		RaisedBy:         dispute.RaisedBy,
		Reason:           dispute.Reason,
		EvidenceURLs:     decodeEvidenceURLs(dispute.EvidenceURLs),
		Status:           dispute.Status,
		EvidenceDeadline: formatOptionalTime(dispute.EvidenceDeadline),
		CreatedAt:        dispute.CreatedAt.Format(time.RFC3339),
	}, nil
}

// AddDisputeEvidence — 1 trong 2 bên bổ sung bằng chứng trước hạn 48h.
func (s *depositService) AddDisputeEvidence(disputeID, userID uint64, req dto.AddEvidenceRequest) (*dto.DisputeResponse, error) {
	dispute, err := s.disputeRepo.GetByID(disputeID)
	if err != nil {
		return nil, errors.New("không tìm thấy tranh chấp")
	}
	if dispute.Status == model.DisputeStatusResolved {
		return nil, errors.New("tranh chấp đã được admin xử lý xong")
	}
	if dispute.Deposit == nil || (dispute.Deposit.CustomerID != userID && dispute.Deposit.BrokerID != userID) {
		return nil, errors.New("bạn không thuộc đơn đặt cọc của tranh chấp này")
	}
	if dispute.EvidenceDeadline != nil && time.Now().After(*dispute.EvidenceDeadline) {
		return nil, errors.New("đã quá hạn upload bằng chứng")
	}
	if len(req.EvidenceURLs) == 0 {
		return nil, errors.New("chưa có bằng chứng nào được gửi lên")
	}

	// Nối thêm bằng chứng mới vào danh sách hiện có
	merged := append(decodeEvidenceURLs(dispute.EvidenceURLs), req.EvidenceURLs...)
	dispute.EvidenceURLs = encodeEvidenceURLs(merged)

	// Bên nào upload thì chuyển tranh chấp sang REVIEWING để admin biết đã có phản hồi
	if dispute.Status == model.DisputeStatusOpen {
		dispute.Status = model.DisputeStatusReviewing
	}
	if err := s.disputeRepo.Save(dispute); err != nil {
		return nil, err
	}

	resp := toDisputeResponse(dispute, nil)
	return &resp, nil
}

// ListDisputes — danh sách tranh chấp cho admin panel.
func (s *depositService) ListDisputes(status string, page, size int) ([]dto.DisputeResponse, int64, error) {
	offset, limit := paginate(page, size)
	items, total, err := s.disputeRepo.List(status, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.DisputeResponse, 0, len(items))
	for i := range items {
		var depositResp *dto.DepositResponse
		if items[i].Deposit != nil {
			mapped := s.toDepositResponse(items[i].Deposit)
			depositResp = &mapped
		}
		result = append(result, toDisputeResponse(&items[i], depositResp))
	}
	return result, total, nil
}

func (s *depositService) GetDisputeDetail(disputeID uint64) (*dto.DisputeResponse, error) {
	dispute, err := s.disputeRepo.GetByID(disputeID)
	if err != nil {
		return nil, errors.New("không tìm thấy tranh chấp")
	}
	var depositResp *dto.DepositResponse
	if dispute.Deposit != nil {
		mapped := s.toDepositResponse(dispute.Deposit)
		depositResp = &mapped
	}
	resp := toDisputeResponse(dispute, depositResp)
	return &resp, nil
}

// ResolveDispute — admin ra quyết định và release tiền.
// 🚨 Tiền KHÔNG bao giờ được release tự động khi đang DISPUTE.
func (s *depositService) ResolveDispute(disputeID, adminID uint64, req dto.ResolveDisputeRequest) (*dto.DisputeResponse, error) {
	dispute, err := s.disputeRepo.GetByID(disputeID)
	if err != nil {
		return nil, errors.New("không tìm thấy tranh chấp")
	}
	if dispute.Status == model.DisputeStatusResolved {
		return nil, errors.New("tranh chấp này đã được xử lý")
	}

	deposit, err := s.depositRepo.GetByID(dispute.DepositID)
	if err != nil {
		return nil, errors.New("không tìm thấy đơn đặt cọc của tranh chấp")
	}

	plan, err := s.buildDisputePlan(req, deposit)
	if err != nil {
		return nil, err
	}
	if req.Note != "" {
		plan.Note = req.Note
	}

	if err := s.settle(deposit, plan); err != nil {
		return nil, err
	}

	now := time.Now()
	dispute.Status = model.DisputeStatusResolved
	dispute.Resolution = req.Resolution
	dispute.ResolvedBy = &adminID
	dispute.ResolvedAt = &now
	if err := s.disputeRepo.Save(dispute); err != nil {
		return nil, err
	}

	// Thông báo kết quả xử lý cho 2 bên
	customerBody := fmt.Sprintf("Admin đã xử lý tranh chấp đơn #%d.\nQuyết định: %s\nSố tiền hoàn cho bạn: %s.",
		deposit.ID, resolutionLabel(req.Resolution), formatMoney(plan.Refund))
	brokerBody := fmt.Sprintf("Admin đã xử lý tranh chấp đơn #%d.\nQuyết định: %s\nSố tiền bạn nhận được: %s.",
		deposit.ID, resolutionLabel(req.Resolution), formatMoney(plan.Transfer))
	if req.Note != "" {
		customerBody += "\nGhi chú từ admin: " + req.Note
		brokerBody += "\nGhi chú từ admin: " + req.Note
	}
	s.notifyBoth(deposit, "dispute_resolved",
		fmt.Sprintf("Kết quả xử lý tranh chấp đơn #%d", deposit.ID),
		customerBody, brokerBody)

	resp := toDisputeResponse(dispute, nil)
	return &resp, nil
}

// buildDisputePlan dựng kế hoạch chia tiền theo quyết định của admin.
func (s *depositService) buildDisputePlan(req dto.ResolveDisputeRequest, deposit *model.Deposit) (settlementPlan, error) {
	switch req.Resolution {
	case model.DisputeResolutionRefundCustomer:
		return settlementPlan{
			Status: model.DepositStatusRefunded,
			Refund: deposit.Amount,
			Note:   "Admin quyết định hoàn toàn bộ tiền cọc cho khách",
		}, nil

	case model.DisputeResolutionTransferBroker:
		return settlementPlan{
			Status:   model.DepositStatusCompleted,
			Transfer: deposit.Amount,
			Note:     "Admin quyết định chuyển toàn bộ tiền cọc cho môi giới",
		}, nil

	case model.DisputeResolutionSplit:
		if req.SplitCustomerPercent < 0 || req.SplitCustomerPercent > 100 {
			return settlementPlan{}, errors.New("tỷ lệ chia cho khách phải trong khoảng 0-100")
		}
		refund := round2(deposit.Amount * float64(req.SplitCustomerPercent) / 100)
		return settlementPlan{
			Status:   model.DepositStatusCompleted,
			Refund:   refund,
			Transfer: round2(deposit.Amount - refund),
			Note:     fmt.Sprintf("Admin quyết định chia tiền: khách %d%%", req.SplitCustomerPercent),
		}, nil
	}

	return settlementPlan{}, errors.New("quyết định xử lý không hợp lệ (REFUND_CUSTOMER / TRANSFER_BROKER / SPLIT)")
}

// openSystemDispute — hệ thống tự mở tranh chấp khi 2 bên báo cáo mâu thuẫn
// hoặc quá hạn báo cáo. Tiền được giữ nguyên tại escrow.
func (s *depositService) openSystemDispute(deposit *model.Deposit, reason string) error {
	deadline := time.Now().Add(time.Duration(s.cfg.DisputeEvidenceHours) * time.Hour)
	dispute := &model.Dispute{
		DepositID:        deposit.ID,
		RaisedBy:         model.DisputeRaisedBySystem,
		Reason:           reason,
		EvidenceURLs:     "[]",
		Status:           model.DisputeStatusOpen,
		EvidenceDeadline: &deadline,
	}
	if err := s.disputeRepo.Create(dispute); err != nil {
		return err
	}

	if err := s.depositRepo.UpdateFields(deposit.ID, map[string]interface{}{
		"status": model.DepositStatusDispute,
	}); err != nil {
		return err
	}
	deposit.Status = model.DepositStatusDispute

	s.notifyBoth(deposit, "dispute_opened",
		fmt.Sprintf("Đơn đặt cọc #%d đang được tạm giữ để xác minh", deposit.ID),
		fmt.Sprintf("Hệ thống ghi nhận mâu thuẫn ở đơn #%d.\nLý do: %s\nTiền cọc bị tạm giữ cho tới khi admin xử lý. Vui lòng upload bằng chứng trong %d giờ.",
			deposit.ID, reason, s.cfg.DisputeEvidenceHours),
		fmt.Sprintf("Hệ thống ghi nhận mâu thuẫn ở đơn #%d.\nLý do: %s\nTiền cọc bị tạm giữ cho tới khi admin xử lý. Vui lòng upload bằng chứng trong %d giờ.",
			deposit.ID, reason, s.cfg.DisputeEvidenceHours),
	)
	return nil
}

func toDisputeResponse(dispute *model.Dispute, depositResp *dto.DepositResponse) dto.DisputeResponse {
	return dto.DisputeResponse{
		ID:               dispute.ID,
		DepositID:        dispute.DepositID,
		RaisedBy:         dispute.RaisedBy,
		Reason:           dispute.Reason,
		EvidenceURLs:     decodeEvidenceURLs(dispute.EvidenceURLs),
		Status:           dispute.Status,
		Resolution:       dispute.Resolution,
		ResolvedBy:       dispute.ResolvedBy,
		ResolvedAt:       formatOptionalTime(dispute.ResolvedAt),
		EvidenceDeadline: formatOptionalTime(dispute.EvidenceDeadline),
		CreatedAt:        dispute.CreatedAt.Format(time.RFC3339),
		Deposit:          depositResp,
	}
}

func formatOptionalTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(time.RFC3339)
}

func resolutionLabel(resolution string) string {
	switch resolution {
	case model.DisputeResolutionRefundCustomer:
		return "Hoàn tiền cho khách"
	case model.DisputeResolutionTransferBroker:
		return "Chuyển tiền cho môi giới"
	case model.DisputeResolutionSplit:
		return "Chia tiền cho cả 2 bên"
	}
	return resolution
}
