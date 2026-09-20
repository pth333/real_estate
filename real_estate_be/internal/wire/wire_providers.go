package wire

import (
	"real_estate_be/internal/global"
	"real_estate_be/pkg/mailer"
	"real_estate_be/pkg/payment"
	"real_estate_be/pkg/sms"

	"gorm.io/gorm"
)

func providerDB() *gorm.DB {
	return global.DB
}

func providerSMS() sms.Provider {
	return sms.NewConsoleProvider()
}

// providerPaymentGateway chọn cổng thanh toán theo cấu hình:
// thiếu merchant key → dùng cổng MOCK để chạy được toàn bộ luồng đặt cọc.
func providerPaymentGateway() payment.Gateway {
	return payment.NewGateway(payment.Config{
		ReturnURL: global.Config.Payment.ReturnURL,
		IPNURL:    global.Config.Payment.IPNURL,
		VNPay: payment.VNPayConfig{
			TmnCode:       global.Config.Payment.VNPay.TmnCode,
			HashSecret:    global.Config.Payment.VNPay.HashSecret,
			PaymentURL:    global.Config.Payment.VNPay.PaymentURL,
			Locale:        global.Config.Payment.VNPay.Locale,
			ExpireMinutes: global.Config.Payment.VNPay.ExpireMinutes,
		},
	})
}

func providerMailer() mailer.Mailer {
	return mailer.NewConsoleMailer()
}
