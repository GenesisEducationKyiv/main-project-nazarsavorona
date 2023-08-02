package rabbitmq

import (
	"context"
	"fmt"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

const ExchangeName = "rate-app"

type Logger struct {
	ch *amqp.Channel
}

func NewLogger(ch *amqp.Channel) *Logger {
	err := declareRabbitMQResources(ch)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err)
	}

	return &Logger{
		ch: ch,
	}
}

func (l *Logger) Log(level models.Level, message string) {
	message = fmt.Sprintf("[%s]: %s", level, message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := l.ch.PublishWithContext(
		ctx,
		ExchangeName, // exchange
		fmt.Sprintf("%s.%s", ExchangeName, level), // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	); err != nil {
		log.Fatalf("failed to publish a message: %s", err)
	}
}
