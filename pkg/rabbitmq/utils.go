package rabbitmq

import (
	"fmt"
	"net"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/logger"
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

	general := "general"
	_, err = ch.QueueDeclare(
		general, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return err
	}

	// bind the queue to the exchange
	err = ch.QueueBind(
		general,                           // queue name
		fmt.Sprintf("%s.#", ExchangeName), // routing key
		ExchangeName,                      // exchange
		false,                             // no-wait
		nil,                               // args
	)
	if err != nil {
		return err
	}

	for _, level := range []logger.Level{logger.Debug, logger.Info, logger.Error} {
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

func ConstructRabbitMQURL(rabbitmqHost, rabbitmqPort,
	rabbitmqUsername, rabbitmqPassword string) string {
	rabbitHostPort := net.JoinHostPort(rabbitmqHost, rabbitmqPort)

	return fmt.Sprintf("amqp://%s:%s@%s/",
		rabbitmqUsername, rabbitmqPassword, rabbitHostPort)
}
