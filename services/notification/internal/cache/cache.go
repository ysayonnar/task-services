package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return client
}
