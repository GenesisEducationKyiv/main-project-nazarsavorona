package email

import (
	"net/smtp"
)

type MockSender struct {
	err error
}

func NewMockSender(err error) *MockSender {
	return &MockSender{err: err}
}

func (m *MockSender) SendMail(_ string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
	return m.err
}
