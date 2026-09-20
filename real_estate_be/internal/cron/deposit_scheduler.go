package cron

import (
	"context"
	"log"
	"time"

	"real_estate_be/internal/usecase"
)

// DepositScheduler chạy các job định kỳ của luồng đặt cọc:
// tự huỷ đơn chưa thanh toán, tự từ chối đơn quá hạn xác nhận, nhắc lịch trước 24h,
// và đẩy các đơn quá hạn báo cáo sang tranh chấp.
type DepositScheduler struct {
	service  usecase.IDepositService
	interval time.Duration
}

func NewDepositScheduler(service usecase.IDepositService, intervalSeconds int) *DepositScheduler {
	if intervalSeconds <= 0 {
		intervalSeconds = 60
	}
	return &DepositScheduler{
		service:  service,
		interval: time.Duration(intervalSeconds) * time.Second,
	}
}

// Start chạy vòng lặp job tới khi context bị huỷ.
func (s *DepositScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("✅ [DepositCron] scheduler started, chu kỳ %s, cổng thanh toán: %s",
		s.interval, s.service.PaymentGatewayName())

	// Chạy ngay 1 lần lúc khởi động để không phải chờ hết chu kỳ đầu
	s.service.RunScheduledTasks()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 [DepositCron] scheduler stopped")
			return
		case <-ticker.C:
			s.service.RunScheduledTasks()
		}
	}
}
