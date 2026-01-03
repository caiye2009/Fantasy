package idgen

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// redisCounter provides atomic sequence generation using Redis
type redisCounter struct {
	rdb *redis.Client
	key string
}

// newRedisCounter creates a new Redis-based counter
func newRedisCounter(rdb *redis.Client) *redisCounter {
	return &redisCounter{
		rdb: rdb,
		key: "order:seq",
	}
}

// next returns the next sequence number atomically
// Uses Redis INCR command which is atomic and thread-safe
func (c *redisCounter) next(ctx context.Context) (uint64, error) {
	val, err := c.rdb.Incr(ctx, c.key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment Redis counter: %w", err)
	}
	return uint64(val), nil
}
