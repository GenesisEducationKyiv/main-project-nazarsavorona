package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/internal/config"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultFilter  = "#" // all logs
	defaultTimeout = 15  // seconds
)

func main() {
	filter := flag.String("filter", defaultFilter, "filter to apply to the logs")
	timeout := flag.Int64("timeout", defaultTimeout, "timeout in seconds for the logs fetching")
	flag.Parse()

	conf, err := config.NewConfig()
	if err != nil {
		log.Panicln(err)
	}

	url := rabbitmq.ConstructRabbitMQURL(conf.RabbitMQHost, conf.RabbitMQPort,
		conf.RabbitMQUsername, conf.RabbitMQPassword)

	rabbitConn, err := amqp.Dial(url)
	if err != nil {
		log.Panicln(err)
	}

	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		log.Panicln(err)
	}

	defer rabbitChannel.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	consumer := rabbitmq.NewConsumer(rabbitChannel)
	consumer.ConsumeErrorLevelMsgs(ctx, *filter, func(message string) {
		log.Printf("Received message: %s", message)
	})
}
