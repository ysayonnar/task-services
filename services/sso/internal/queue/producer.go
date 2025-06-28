package queue

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sso/internal/config"
)

const (
	KEY_USER_REGISTERED = "user.registered"
	KEY_USER_DELETED    = "user.deleted"
)

type Broker struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func New(cfg config.Config) (*Broker, error) {
	const op = `queue.New`

	url := fmt.Sprintf("amqp://%s:%s@%s/", cfg.RabbitMQ.Username, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	err = ch.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return &Broker{Conn: conn, Ch: ch}, nil
}

func (b *Broker) GracefulStop() {
	b.Ch.Close()
	b.Conn.Close()
}
