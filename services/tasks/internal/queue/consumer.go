package queue

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

type UserDeletedDto struct {
	UserId int64 `json:"user_id"`
}

func (b *Broker) Consume() error {
	const op = "queue.Consume"

	msgs, err := b.Ch.Consume(QUEUE_NAME, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	go func() {
		for d := range msgs {
			switch d.RoutingKey {
			case KEY_USER_DELETED:
				b.handleUserDeleted(d)
			}
		}
	}()

	return nil
}

func (b *Broker) handleUserDeleted(d amqp.Delivery) {
	const op = "queue.handleUserDeleted"
	log := b.Log.With(slog.String("op", op))

	var dto UserDeletedDto
	err := json.Unmarshal(d.Body, &dto)
	if err != nil {
		log.Error("error while parsing json", "error", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = b.Storage.DeleteAllTasksByUserId(ctx, dto.UserId)
	if err != nil {
		log.Error("error while parsing json", "error", err.Error())
	}
}
