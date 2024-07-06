package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
