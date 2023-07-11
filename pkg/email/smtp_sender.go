package email

import (
	"net/smtp"
)

type SMTPMailSender struct {
}

func (c *SMTPMailSender) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}
