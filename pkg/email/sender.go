package email

import (
	"net/smtp"
	"strings"
)

type sendMailFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

type Sender struct {
	auth     smtp.Auth
	email    string
	hostURI  string
	sendMail sendMailFunc
}

func NewSender(email, hostURI string, auth smtp.Auth, sendMail sendMailFunc) *Sender {
	return &Sender{
		auth:     auth,
		email:    email,
		hostURI:  hostURI,
		sendMail: sendMail,
	}
}

func (s *Sender) SendEmail(toEmail string, subject string, body string) error {
	body = strings.ReplaceAll(body, "\n", "<br>")

	message := []byte("Subject: " + subject + "\n" +
		"To: " + toEmail + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"<html><body><h3>" + body + "</h3></body></html>")

	return s.sendMail(s.hostURI, s.auth, s.email, []string{toEmail}, message)
}
