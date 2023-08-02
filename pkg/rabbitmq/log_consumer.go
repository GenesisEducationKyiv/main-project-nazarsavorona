package rabbitmq

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch *amqp091.Channel
}

func NewConsumer(ch *amqp091.Channel) *Consumer {
	return &Consumer{
		ch: ch,
	}
}

const ConsumerName = "rate-app-consumer"

func (c *Consumer) ConsumeMessages(context context.Context, queue string, handler func(message string)) {
	msgs, err := c.ch.Consume(
		queue,        // queue
		ConsumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		log.Fatalf("failed to register a consumer: %s", err)
	}

	for {
		select {
		case <-context.Done():
			return
		case msg := <-msgs:
			handler(string(msg.Body))
		}
	}
}
