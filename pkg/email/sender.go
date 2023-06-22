package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Sender struct {
	auth    smtp.Auth
	email   string
	hostURI string
}

func NewSender(smtpHost, smtpPort, email, password string) *Sender {
	return &Sender{
		auth:    smtp.PlainAuth("", email, password, smtpHost),
		email:   email,
		hostURI: fmt.Sprintf("%s:%s", smtpHost, smtpPort),
	}
}

func (s *Sender) SendEmail(toEmail string, subject string, body string) error {
	body = strings.ReplaceAll(body, "\n", "<br>")

	message := []byte("Subject: " + subject + "\n" +
		"To: " + toEmail + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"<html><body><h3>" + body + "</h3></body></html>")

	return smtp.SendMail(s.hostURI, s.auth, s.email, []string{toEmail}, message)
}