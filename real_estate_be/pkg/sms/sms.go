package sms

import "fmt"

// Provider là interface cho các dịch vụ gửi SMS
type Provider interface {
	Send(phone, otp string) error
}

// ConsoleProvider in OTP ra console (dùng cho môi trường dev)
type ConsoleProvider struct{}

func NewConsoleProvider() Provider {
	return &ConsoleProvider{}
}

func (p *ConsoleProvider) Send(phone, otp string) error {
	fmt.Printf("[SMS] Gửi tới %s: mã OTP là %s\n", phone, otp)
	return nil
}
