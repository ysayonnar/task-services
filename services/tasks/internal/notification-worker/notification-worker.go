package notification_worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"tasks/internal/queue"
	"tasks/internal/storage"
	"time"
)

const defaultTickerDelta = 1

type NotificationWorker struct {
	Log     *slog.Logger
	Storage *storage.Storage
	Broker  *queue.Broker
}

func New(log *slog.Logger, s *storage.Storage, b *queue.Broker) *NotificationWorker {
	return &NotificationWorker{
		Log:     log,
		Storage: s,
		Broker:  b,
	}
}

func (w *NotificationWorker) Listen() {
	const op = "notification_worker.Listen"
	log := w.Log.With(slog.String("op", op))

	ticker := time.NewTicker(defaultTickerDelta * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

			tasks, err := w.Storage.GetTasksToNotify(ctx)
			if err != nil {
				log.Error("error while trying to get tasks to notify", "error", err.Error())
				continue
			}

			for _, task := range tasks {
				jsonTask, err := json.Marshal(task)
				if err != nil {
					log.Error("error while marshaling task", "error", err.Error())
					continue
				}

				err = w.Broker.Publish(ctx, queue.KEY_TASK_NOTIFICATE, jsonTask)
				if err != nil {
					log.Error("error while publishing message to rabbitmq", "error", err.Error())
				}
			}

			cancel()
		}
	}
}
