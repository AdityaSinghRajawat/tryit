// All methods are nil-tolerant: when Redis isn't configured the wrapper /
// client is nil and calls become safe no-ops (GetKey returns a miss). This
// keeps the server runnable without Redis at the cost of no parse cache.
package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil struct {
	RedisClient *redis.Client
}

func NewRedisUtilManager(client *redis.Client) *RedisUtil {
	return &RedisUtil{RedisClient: client}
}

func (r *RedisUtil) SetKey(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	if r == nil || r.RedisClient == nil {
		return nil
	}
	if err := r.RedisClient.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis set %q: %w", key, err)
	}
	return nil
}

func (r *RedisUtil) GetKey(ctx context.Context, key string) (string, error) {
	if r == nil || r.RedisClient == nil {
		return "", nil
	}
	result, err := r.RedisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("redis get %q: %w", key, err)
	}
	return result, nil
}

func (r *RedisUtil) DeleteKey(ctx context.Context, key string) error {
	if r == nil || r.RedisClient == nil {
		return nil
	}
	if err := r.RedisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis del %q: %w", key, err)
	}
	return nil
}
