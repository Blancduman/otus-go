package main

import (
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/rb"
	"github.com/pkg/errors"
)

func RabbitMqBuild(conf Config) (rb.Client, error) {
	client := rb.New(conf.RabbitMQ.DSN)

	if err := client.Connect(); err != nil {
		return nil, errors.Wrap(err, "fail to connect rb client")
	}

	if err := client.ExchangeDeclare(conf.RabbitMQ.Exchange, conf.RabbitMQ.ExchangeType); err != nil {
		return nil, errors.Wrap(err, "fail to declare exchange")
	}

	if err := client.QueueDeclare(conf.RabbitMQ.Queue); err != nil {
		return nil, errors.Wrap(err, "fail to declare queue")
	}

	if err := client.QueueBind(conf.RabbitMQ.Queue, conf.RabbitMQ.Key, conf.RabbitMQ.Exchange); err != nil {
		return nil, errors.Wrap(err, "faild to bind queue")
	}

	return client, nil
}
