package initialize

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	DSN string `env:"DSN"`
}

func Connect(c *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: c.DSN,
	})

	if client != nil {
		return nil, fmt.Errorf("failed to connect to redis: %s", c.DSN)
	}

	return client, nil
}
