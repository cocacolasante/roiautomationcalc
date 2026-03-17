package notification

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailNotifier struct {
	host      string
	port      int
	user      string
	pass      string
	fromEmail string
	fromName  string
}

func NewEmailNotifier(host string, port int, user, pass, fromEmail, fromName string) *EmailNotifier {
	return &EmailNotifier{host: host, port: port, user: user, pass: pass, fromEmail: fromEmail, fromName: fromName}
}

func (e *EmailNotifier) Send(_ context.Context, to, subject, htmlBody string) error {
	if e.host == "" {
		return fmt.Errorf("SMTP not configured")
	}
	from := fmt.Sprintf("%s <%s>", e.fromName, e.fromEmail)
	headers := []string{
		"From: " + from, "To: " + to, "Subject: " + subject,
		"MIME-Version: 1.0", "Content-Type: text/html; charset=UTF-8",
	}
	msg := strings.Join(headers, "\r\n") + "\r\n\r\n" + htmlBody
	auth := smtp.PlainAuth("", e.user, e.pass, e.host)
	addr := fmt.Sprintf("%s:%d", e.host, e.port)
	return smtp.SendMail(addr, auth, e.fromEmail, []string{to}, []byte(msg))
}
