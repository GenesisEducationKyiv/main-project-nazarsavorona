package email

import (
	"net/smtp"
)

type (
	MailSender interface {
		SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
	}

	Sender struct {
		auth      smtp.Auth
		fromEmail string
		hostURI   string
		sender    MailSender
	}
)

func NewSender(fromEmail, hostURI string, auth smtp.Auth, sender MailSender) *Sender {
	return &Sender{
		auth:      auth,
		fromEmail: fromEmail,
		hostURI:   hostURI,
		sender:    sender,
	}
}

type MessageConstructStrategy interface {
	Construct(subject, body string) []byte
}

func (s *Sender) SendEmail(toEmail, subject, body string, strategy MessageConstructStrategy) error {
	message := strategy.Construct(subject, body)

	return s.sender.SendMail(s.hostURI, s.auth, s.fromEmail, []string{toEmail}, message)
}
