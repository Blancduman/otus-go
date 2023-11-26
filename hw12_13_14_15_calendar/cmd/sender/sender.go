package main

import (
	"context"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/sender"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func senderCmd(ctx context.Context, conf Config) *cobra.Command {
	return &cobra.Command{
		Use:   "sender",
		Short: "start scheduler",
		RunE: func(cmd *cobra.Command, args []string) error {
			logg := logger.New(conf.Logger.Level)

			rbClient, err := RabbitMqBuild(conf)
			if err != nil {
				return errors.Wrap(err, "could not build rabbitmq client")
			}
			defer func() {
				if err := rbClient.Close(); err != nil {
					logg.Error(errors.Wrap(err, "fail to close rb client").Error())
				}
			}()

			deliveries, err := rbClient.Consume(conf.RabbitMQ.Queue, conf.RabbitMQ.ConsumerTag)
			if err != nil {
				return errors.Wrap(err, "rb queue consume")
			}

			send := sender.New(rbClient, logg, conf.RabbitMQ.Exchange, conf.RabbitMQ.Key)

			if err := send.Run(ctx, deliveries); err != nil {
				return errors.Wrap(err, "fail to start sender")
			}

			return nil
		},
	}
}
