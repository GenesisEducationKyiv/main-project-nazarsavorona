package rabbitmq

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log"
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

func (c *Consumer) ConsumeErrorLevelMsgs(context context.Context, topic string, handler func(message string)) {
	q, err := c.ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err)
	}

	err = c.ch.QueueBind(q.Name, topic, ExchangeName, false, nil)
	if err != nil {
		log.Fatalf("failed to bind a queue: %s", err)
	}

	msgs, err := c.ch.Consume(
		q.Name,       // queue
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
