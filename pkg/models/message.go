package models

import "fmt"

type Message struct {
	Subject string
	Body    string
}

func NewMessageFromRate(r *Rate) *Message {
	return &Message{
		Subject: fmt.Sprintf("%s rate", r.From),
		Body:    fmt.Sprintf("1 %s = %.2f %s", r.From, r.Rate, r.To),
	}
}
