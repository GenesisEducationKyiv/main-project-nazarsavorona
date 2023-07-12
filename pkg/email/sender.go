package email

import (
	"net/smtp"
)

type Sender struct {
	auth      smtp.Auth
	fromEmail string
	hostURI   string
}

func NewSender(fromEmail, hostURI string, auth smtp.Auth) *Sender {
	return &Sender{
		auth:      auth,
		fromEmail: fromEmail,
		hostURI:   hostURI,
	}
}

func (s *Sender) SendEmail(toEmail string, message []byte) error {
	return smtp.SendMail(s.hostURI, s.auth, s.fromEmail, []string{toEmail}, message)
}
