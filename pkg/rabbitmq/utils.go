package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
)

func declareRabbitMQResources(ch *amqp091.Channel) error {
	err := ch.ExchangeDeclare(
		ExchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	// bind the queue to the exchange
	err = ch.QueueBind(
		QueueName,    // queue name
		QueueName,    // routing key
		ExchangeName, // exchange
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		return err
	}

	return nil
}
