package queue

import "fmt"

func (b *Broker) Consume() error {
	const op = "queue.Consume"

	msgs, err := b.Ch.Consume(QUEUE_NAME, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	go func() {
		for d := range msgs {
			switch d.RoutingKey {
			//TODO: handle here
			}
		}
	}()

	return nil
}
