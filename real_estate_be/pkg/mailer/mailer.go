package mailer

import (
	"fmt"
	"log"
)

// Message — nội dung 1 email cần gửi
type Message struct {
	To      string
	Subject string
	Body    string
}

// Mailer — interface gửi email. Hiện dùng ConsoleMailer (log ra console),
// khi có SMTP/SES thật thì thêm adapter mới và đổi provider ở wire.
type Mailer interface {
	Send(msg Message) error
}

// ConsoleMailer — in email ra console để dev theo dõi luồng thông báo.
type ConsoleMailer struct{}

func NewConsoleMailer() Mailer {
	return &ConsoleMailer{}
}

func (m *ConsoleMailer) Send(msg Message) error {
	if msg.To == "" {
		return fmt.Errorf("thiếu email người nhận")
	}
	log.Printf("📧 [MAIL] to=%s | subject=%s\n%s", msg.To, msg.Subject, msg.Body)
	return nil
}
