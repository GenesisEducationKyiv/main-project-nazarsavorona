package email

import (
	"net/smtp"
	"strings"
)

type MailSender interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

type Sender struct {
	auth      smtp.Auth
	fromEmail string
	hostURI   string
	sender    MailSender
}

func NewSender(fromEmail, hostURI string, auth smtp.Auth, sender MailSender) *Sender {
	return &Sender{
		auth:      auth,
		fromEmail: fromEmail,
		hostURI:   hostURI,
		sender:    sender,
	}
}

func (s *Sender) SendEmail(toEmail, subject, body string) error {
	body = strings.ReplaceAll(body, "\n", "<br>")
	message := constructEmailMessage(subject, body)

	return s.sender.SendMail(s.hostURI, s.auth, s.fromEmail, []string{toEmail}, message)
}

func constructEmailMessage(subject, body string) []byte {
	htmlBody := "<html><body><h3>" + body + "</h3></body></html>"
	message := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		htmlBody)

	return message
}
