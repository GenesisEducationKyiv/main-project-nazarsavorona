package rabbitmq

import (
	"fmt"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/rabbitmq/amqp091-go"
)

func declareRabbitMQResources(ch *amqp091.Channel) error {
	err := ch.ExchangeDeclare(
		ExchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	// get every log level list
	for _, level := range []models.Level{models.Debug, models.Info, models.Error} {
		_, err = ch.QueueDeclare(
			level.String(), // name
			true,           // durable
			false,          // delete when unused
			false,          // exclusive
			false,          // no-wait
			nil,            // args
		)
		if err != nil {
			return err
		}

		// bind the queue to the exchange
		err = ch.QueueBind(
			level.String(), // queue name
			fmt.Sprintf("%s.%s", ExchangeName, level), // routing key
			ExchangeName, // exchange
			false,        // no-wait
			nil,          // args
		)
		if err != nil {
			return err
		}
	}

	return nil
}
