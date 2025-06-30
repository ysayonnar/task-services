package queue

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"tasks/internal/config"
	"tasks/internal/storage"
)

const (
	EXCHANGE_NAME       = "events"
	QUEUE_NAME          = "tasks_queue"
	KEY_USER_REGISTERED = "user.registered"
	KEY_USER_DELETED    = "user.deleted"
	KEY_TASK_NOTIFICATE = "task.notificate"
)

type Broker struct {
	Conn    *amqp.Connection
	Ch      *amqp.Channel
	Storage *storage.Storage
	Log     *slog.Logger
}

func New(cfg config.Config, logger *slog.Logger, storage *storage.Storage) (*Broker, error) {
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

	err = ch.ExchangeDeclare(EXCHANGE_NAME, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	q, err := ch.QueueDeclare(
		QUEUE_NAME,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	err = ch.QueueBind(q.Name, KEY_USER_DELETED, EXCHANGE_NAME, false, nil)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return &Broker{Conn: conn, Ch: ch, Log: logger, Storage: storage}, nil
}

func (b *Broker) Publish(ctx context.Context, key string, body []byte) error {
	const op = "queue.Publish"

	err := b.Ch.PublishWithContext(ctx, EXCHANGE_NAME, key, false, false, amqp.Publishing{ContentType: "application/json", Body: body})
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	return nil
}

func (b *Broker) GracefulStop() {
	b.Ch.Close()
	b.Conn.Close()
}
