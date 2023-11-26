package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/rb"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Sender struct {
	Client   rb.Client
	Logger   *logger.Logger
	Exchange string
	Key      string
}

func New(client rb.Client, logger *logger.Logger, exchange string, key string) *Sender {
	return &Sender{
		Client:   client,
		Logger:   logger,
		Exchange: exchange,
		Key:      key,
	}
}

func (s *Sender) Run(ctx context.Context, deliveries <-chan amqp.Delivery) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("sender context is done")
		case d := <-deliveries:
			var message rb.Message

			if err := json.Unmarshal(d.Body, &message); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf("unmarshal delivery message: %v %q", d.DeliveryTag, d.Body),
				)
			}

			s.Logger.Info(fmt.Sprintf("%v %v", d.DeliveryTag, &message))

			if err := d.Ack(false); err != nil {
				return errors.Wrap(err, "ack delivery message")
			}
		}
	}
}
