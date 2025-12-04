package price

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	
	"back/internal/vendor"
	"back/pkg/repo"
)

type ProcessPriceService struct {
	db            *gorm.DB
	vendorService interface {
		Get(ctx context.Context, id uint) (*vendor.Vendor, error)
	}
	rdb *redis.Client
}

func NewProcessPriceService(
	db *gorm.DB,
	vendorService interface{ Get(ctx context.Context, id uint) (*vendor.Vendor, error) },
	rdb *redis.Client,
) *ProcessPriceService {
	return &ProcessPriceService{
		db:            db,
		vendorService: vendorService,
		rdb:           rdb,
	}
}

// Quote 厂商报价
func (s *ProcessPriceService) Quote(ctx context.Context, vendorID, processID uint, price float64) error {
	// 1. 写入数据库
	processPrice := &ProcessPrice{
		VendorID:  vendorID,
		ProcessID: processID,
		Price:     price,
		QuotedAt:  time.Now(),
	}
	
	priceRepo := repo.NewRepo[ProcessPrice](s.db)
	if err := priceRepo.Create(ctx, processPrice); err != nil {
		return err
	}

	// 2. 获取厂商名称
	v, err := s.vendorService.Get(ctx, vendorID)
	if err != nil {
		return err
	}

	// 3. 更新 Redis (比对最低/最高)
	priceData := PriceData{
		ID:         processPrice.ID,
		Price:      price,
		VendorID:   vendorID,
		VendorName: v.Name,
		QuotedAt:   processPrice.QuotedAt,
	}

	// 更新最低价
	if err := s.compareAndUpdateMin(ctx, processID, priceData); err != nil {
		return err
	}

	// 更新最高价
	if err := s.compareAndUpdateMax(ctx, processID, priceData); err != nil {
		return err
	}

	return nil
}

// compareAndUpdateMin 比对并更新最低价
func (s *ProcessPriceService) compareAndUpdateMin(ctx context.Context, processID uint, newPrice PriceData) error {
	key := fmt.Sprintf("price:process:%d:min", processID)

	// 读取当前最低价
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Redis 中没有,直接写入
		return s.setPrice(ctx, key, newPrice)
	} else if err != nil {
		return err
	}

	// 解析当前最低价
	var currentMin PriceData
	if err := json.Unmarshal([]byte(val), &currentMin); err != nil {
		return err
	}

	// 比对
	if newPrice.Price < currentMin.Price {
		return s.setPrice(ctx, key, newPrice)
	}

	return nil
}

// compareAndUpdateMax 比对并更新最高价
func (s *ProcessPriceService) compareAndUpdateMax(ctx context.Context, processID uint, newPrice PriceData) error {
	key := fmt.Sprintf("price:process:%d:max", processID)

	// 读取当前最高价
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Redis 中没有,直接写入
		return s.setPrice(ctx, key, newPrice)
	} else if err != nil {
		return err
	}

	// 解析当前最高价
	var currentMax PriceData
	if err := json.Unmarshal([]byte(val), &currentMax); err != nil {
		return err
	}

	// 比对
	if newPrice.Price > currentMax.Price {
		return s.setPrice(ctx, key, newPrice)
	}

	return nil
}

// setPrice 写入 Redis
func (s *ProcessPriceService) setPrice(ctx context.Context, key string, priceData PriceData) error {
	data, err := json.Marshal(priceData)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, key, data, 0).Err() // 无 TTL
}

// GetMinPrice 获取最低价
func (s *ProcessPriceService) GetMinPrice(processID uint) (*PriceData, error) {
	ctx := context.Background()
	key := fmt.Sprintf("price:process:%d:min", processID)

	// 尝试从 Redis 读取
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Redis 未命中,从数据库加载
		return s.loadMinPriceFromDB(ctx, processID)
	} else if err != nil {
		return nil, err
	}

	// 解析
	var priceData PriceData
	if err := json.Unmarshal([]byte(val), &priceData); err != nil {
		return nil, err
	}

	return &priceData, nil
}

// GetMaxPrice 获取最高价
func (s *ProcessPriceService) GetMaxPrice(processID uint) (*PriceData, error) {
	ctx := context.Background()
	key := fmt.Sprintf("price:process:%d:max", processID)

	// 尝试从 Redis 读取
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Redis 未命中,从数据库加载
		return s.loadMaxPriceFromDB(ctx, processID)
	} else if err != nil {
		return nil, err
	}

	// 解析
	var priceData PriceData
	if err := json.Unmarshal([]byte(val), &priceData); err != nil {
		return nil, err
	}

	return &priceData, nil
}

// loadMinPriceFromDB 从数据库加载最低价并写入 Redis
func (s *ProcessPriceService) loadMinPriceFromDB(ctx context.Context, processID uint) (*PriceData, error) {
	// 查询最低价
	var price ProcessPrice
	err := s.db.WithContext(ctx).
		Where("process_id = ?", processID).
		Order("price ASC").
		First(&price).Error
	if err != nil {
		return nil, err
	}

	// 获取厂商名称
	v, err := s.vendorService.Get(ctx, price.VendorID)
	if err != nil {
		return nil, err
	}

	priceData := &PriceData{
		ID:         price.ID,
		Price:      price.Price,
		VendorID:   price.VendorID,
		VendorName: v.Name,
		QuotedAt:   price.QuotedAt,
	}

	// 写入 Redis
	key := fmt.Sprintf("price:process:%d:min", processID)
	if err := s.setPrice(ctx, key, *priceData); err != nil {
		return nil, err
	}

	return priceData, nil
}

// loadMaxPriceFromDB 从数据库加载最高价并写入 Redis
func (s *ProcessPriceService) loadMaxPriceFromDB(ctx context.Context, processID uint) (*PriceData, error) {
	// 查询最高价
	var price ProcessPrice
	err := s.db.WithContext(ctx).
		Where("process_id = ?", processID).
		Order("price DESC").
		First(&price).Error
	if err != nil {
		return nil, err
	}

	// 获取厂商名称
	v, err := s.vendorService.Get(ctx, price.VendorID)
	if err != nil {
		return nil, err
	}

	priceData := &PriceData{
		ID:         price.ID,
		Price:      price.Price,
		VendorID:   price.VendorID,
		VendorName: v.Name,
		QuotedAt:   price.QuotedAt,
	}

	// 写入 Redis
	key := fmt.Sprintf("price:process:%d:max", processID)
	if err := s.setPrice(ctx, key, *priceData); err != nil {
		return nil, err
	}

	return priceData, nil
}

// GetHistory 获取价格历史
func (s *ProcessPriceService) GetHistory(ctx context.Context, processID uint) ([]ProcessPrice, error) {
	var prices []ProcessPrice
	err := s.db.WithContext(ctx).
		Where("process_id = ?", processID).
		Order("quoted_at DESC").
		Find(&prices).Error
	
	return prices, err
}