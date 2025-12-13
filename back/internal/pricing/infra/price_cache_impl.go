package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"back/internal/pricing/domain"
)

var ErrCacheMiss = errors.New("cache miss")

// PriceCacheImpl Redis 缓存实现
type PriceCacheImpl struct {
	rdb *redis.Client
}

// NewPriceCacheImpl 创建缓存实现
func NewPriceCacheImpl(rdb *redis.Client) domain.PriceCache {
	return &PriceCacheImpl{rdb: rdb}
}

// GetMin 获取最低价
func (c *PriceCacheImpl) GetMin(ctx context.Context, targetType string, targetID uint) (*domain.PriceData, error) {
	key := fmt.Sprintf("price:%s:%d:min", targetType, targetID)
	
	val, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	
	var data domain.PriceData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	
	return &data, nil
}

// SetMin 设置最低价
func (c *PriceCacheImpl) SetMin(ctx context.Context, targetType string, targetID uint, data *domain.PriceData) error {
	key := fmt.Sprintf("price:%s:%d:min", targetType, targetID)
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	return c.rdb.Set(ctx, key, jsonData, 0).Err()
}

// UpdateMin 更新最低价
func (c *PriceCacheImpl) UpdateMin(ctx context.Context, targetType string, targetID uint, newPrice *domain.PriceData) error {
	current, err := c.GetMin(ctx, targetType, targetID)
	if errors.Is(err, ErrCacheMiss) {
		// 没有缓存，直接写入
		return c.SetMin(ctx, targetType, targetID, newPrice)
	}
	if err != nil {
		return err
	}
	
	// 比对
	if newPrice.Price < current.Price {
		return c.SetMin(ctx, targetType, targetID, newPrice)
	}
	
	return nil
}

// GetMax 获取最高价
func (c *PriceCacheImpl) GetMax(ctx context.Context, targetType string, targetID uint) (*domain.PriceData, error) {
	key := fmt.Sprintf("price:%s:%d:max", targetType, targetID)
	
	val, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	
	var data domain.PriceData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	
	return &data, nil
}

// SetMax 设置最高价
func (c *PriceCacheImpl) SetMax(ctx context.Context, targetType string, targetID uint, data *domain.PriceData) error {
	key := fmt.Sprintf("price:%s:%d:max", targetType, targetID)
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	return c.rdb.Set(ctx, key, jsonData, 0).Err()
}

// UpdateMax 更新最高价
func (c *PriceCacheImpl) UpdateMax(ctx context.Context, targetType string, targetID uint, newPrice *domain.PriceData) error {
	current, err := c.GetMax(ctx, targetType, targetID)
	if errors.Is(err, ErrCacheMiss) {
		// 没有缓存，直接写入
		return c.SetMax(ctx, targetType, targetID, newPrice)
	}
	if err != nil {
		return err
	}
	
	// 比对
	if newPrice.Price > current.Price {
		return c.SetMax(ctx, targetType, targetID, newPrice)
	}
	
	return nil
}