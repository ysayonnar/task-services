package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type UserDeletedDto struct {
	UserId int64 `json:"user_id"`
}

type UserCreatedDto struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
}

type TaskNotificateDto struct {
	UserId      int64  `json:"user_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
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
			case KEY_USER_REGISTERED:
				b.handleUserRegistered(d)
			case KEY_USER_DELETED:
				b.handleUserDeleted(d)
			}
		}
	}()

	return nil
}

func (b *Broker) handleUserRegistered(d amqp.Delivery) {
	const op = "queue.handleUserRegistered"
	log := b.Log.With(slog.String("op", op))

	if len(d.Body) == 0 {
		d.Nack(false, false)
		return
	}

	var dto UserCreatedDto
	err := json.Unmarshal(d.Body, &dto)
	if err != nil {
		log.Error("error while parsing json", "error", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = b.RedisClient.Set(ctx, strconv.Itoa(int(dto.UserId)), dto.Email, 0).Err()
	if err != nil {
		log.Error("error while writing to redis", "error", err.Error())
	}
}

func (b *Broker) handleUserDeleted(d amqp.Delivery) {
	const op = "queue.handleUserDeleted"
	log := b.Log.With(slog.String("op", op))

	if len(d.Body) == 0 {
		d.Nack(false, false)
		return
	}

	var dto UserDeletedDto
	err := json.Unmarshal(d.Body, &dto)
	if err != nil {
		log.Error("error while parsing json", "error", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = b.RedisClient.Del(ctx, strconv.Itoa(int(dto.UserId))).Err()
	if err != nil {
		log.Error("error while writing to redis", "error", err.Error())
	}
}
