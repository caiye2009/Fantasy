package idgen

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// OrderNumberGenerator generates unique order numbers using Crockford Base32 encoding
type OrderNumberGenerator struct {
	counter *redisCounter
}

// NewOrderNumberGenerator creates a new order number generator
func NewOrderNumberGenerator(rdb *redis.Client) *OrderNumberGenerator {
	return &OrderNumberGenerator{
		counter: newRedisCounter(rdb),
	}
}

// Generate creates a new unique order number
// Returns a 6-character Crockford Base32 encoded string (e.g., "B5MQ8G")
func (g *OrderNumberGenerator) Generate(ctx context.Context) (string, error) {
	seq, err := g.counter.next(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get next sequence: %w", err)
	}

	orderNo := EncodeCrockford(seq)
	return orderNo, nil
}
