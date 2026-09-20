package sms

import (
	"fmt"
	"io"
	"net/http"
	"real_estate_be/internal/global"
	"strings"
)

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
	payload := fmt.Sprintf(`{
		"messages": [{
			"destinations": [{"to": "%s"}],
			"sender": "447491163443",
			"content": {"text": "Ma OTP cua ban la: %s"}
		}]
	}`, phone, otp)

	req, err := http.NewRequest("POST",
		"https://ndgln2.api.infobip.com/sms/3/messages",
		strings.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "App "+global.Config.Infobip.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send OTP: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Infobip response: %s\n", body)

	return nil
}
