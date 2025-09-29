package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedis(addr string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr, DB: db,
	})
}

func Ping(ctx context.Context, c *redis.Client) error {
	return c.Ping(ctx).Err()
}

func TTLSetJSON(ctx context.Context, c *redis.Client, key string, val []byte, ttl time.Duration) error {
	return c.Set(ctx, key, val, ttl).Err()
}
