package models

import "fmt"

type (
	Message struct {
		Subject string
		Body    string
	}

	RateProvider interface {
		From() string
		To() string
		Rate() float64
	}
)

func NewMessageFromRate(r RateProvider) *Message {
	subject := constructSubject(r)
	body := constructBody(r)
	return &Message{
		Subject: subject,
		Body:    body,
	}
}

func constructSubject(r RateProvider) string {
	return fmt.Sprintf("%s rate", r.From())
}

func constructBody(r RateProvider) string {
	return fmt.Sprintf("1 %s = %.2f %s", r.From(), r.Rate(), r.To())
}
