package models

import "fmt"

type Message struct {
	Subject string
	Body    string
}

func NewMessageFromRate(r RateProvider) *Message {
	subject := constructSubject(r)
	body := constructBody(r)
	return &Message{
		Subject: subject,
		Body:    body,
	}
}

type RateProvider interface {
	GetFrom() string
	GetTo() string
	GetRate() float64
}

func constructSubject(r RateProvider) string {
	return fmt.Sprintf("%s rate", r.GetFrom())
}

func constructBody(r RateProvider) string {
	return fmt.Sprintf("1 %s = %.2f %s", r.GetFrom(), r.GetRate(), r.GetTo())
}
