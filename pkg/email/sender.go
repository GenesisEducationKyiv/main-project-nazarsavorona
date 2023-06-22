package email

import (
	"net/smtp"
	"strings"
)

type Sender struct {
	auth    smtp.Auth
	email   string
	hostURI string
}

func NewSender(hostURI, email, password string) *Sender {
	return &Sender{
		auth:    NewLoginAuth(email, password),
		email:   email,
		hostURI: hostURI,
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
