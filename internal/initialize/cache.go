package initialize

import (
	"fmt"
	"net/url"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	DSN url.URL `env:"DSN"`
}

func Connect(c *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.DSN.String(),
		Password: "",
		DB:       0,
	})

	if client == nil {
		return nil, fmt.Errorf("failed to connect to redis: %s", c.DSN)
	}

	return client, nil
}
