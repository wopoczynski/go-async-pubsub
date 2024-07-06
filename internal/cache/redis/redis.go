package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/wopoczynski/playground/internal/cache"
)

var _ cache.Cache = (*RedisCache)(nil)

var ErrNoDataFound = errors.New("no data found in cache")

type RedisCache struct {
	client Client
}

type Client interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}

func (c *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	_, err := c.client.Set(ctx, key, value, ttl).Result()
	if err != nil {
		return fmt.Errorf("save in redis: %w", err)
	}

	return nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	v, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(redis.Nil, err) {
			err = ErrNoDataFound
		}
		return "", fmt.Errorf("read from redis: %w", err)
	}

	return v, nil
}
