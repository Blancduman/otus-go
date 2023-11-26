package rb

import (
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Client interface {
	Connect() error
	Close() error
	ExchangeDeclare(exchange string, exchangeType string) error
	QueueDeclare(name string) error
	QueueBind(name string, key string, exchange string) error
	Publish(exchange string, routingKey string, body []byte) error
	NotifyPublish(ch chan amqp.Confirmation) (chan amqp.Confirmation, error)
	Consume(name string, tag string) (<-chan amqp.Delivery, error)
}

type Message struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	OwnerID   int64     `json:"ownerId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type client struct {
	Client
	dsn     string
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

var ErrNoChannel = errors.New("rb channel is not initialized")

func New(dsn string) Client {
	return &client{
		dsn: dsn,
	}
}

func (c *client) Connect() error {
	conn, err := amqp.Dial(c.dsn)
	if err != nil {
		return errors.Wrap(err, "connect to rabbitmq")
	}

	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "rabbitmq create channel")
	}

	c.conn = conn
	c.channel = ch

	return nil
}

func (c *client) Close() error {
	if err := c.conn.Close(); err != nil {
		return errors.Wrap(err, "close connection rabbitmq")
	}

	return nil
}

func (c *client) ExchangeDeclare(exchange string, exchangeType string) error {
	if c.channel == nil {
		return ErrNoChannel
	}

	err := c.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "exchange declare")
	}

	return nil
}

func (c *client) QueueDeclare(name string) error {
	if c.channel == nil {
		return ErrNoChannel
	}

	queue, err := c.channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "rb queue declare")
	}

	c.queue = &queue

	return nil
}

func (c *client) QueueBind(name string, key string, exchange string) error {
	if c.channel == nil {
		return ErrNoChannel
	}

	err := c.channel.QueueBind(name, key, exchange, false, nil)
	if err != nil {
		return errors.Wrap(err, "rb queue bind")
	}

	return nil
}

func (c *client) Publish(exchange string, routingKey string, body []byte) error {
	if c.channel == nil {
		return ErrNoChannel
	}

	err := c.channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		Headers:      amqp.Table{},
		ContentType:  "text/plain",
		DeliveryMode: amqp.Transient,
		Body:         body,
	})
	if err != nil {
		return errors.Wrap(err, "rb publishing")
	}

	return nil
}

func (c *client) NotifyPublish(ch chan amqp.Confirmation) (chan amqp.Confirmation, error) {
	if c.channel == nil {
		return nil, ErrNoChannel
	}

	return c.channel.NotifyPublish(ch), nil
}

func (c *client) Consume(name string, tag string) (<-chan amqp.Delivery, error) {
	if c.channel == nil {
		return nil, ErrNoChannel
	}

	deliveries, err := c.channel.Consume(name, tag, false, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "rb consume")
	}

	return deliveries, nil
}
