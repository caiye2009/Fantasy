package price

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"back/internal/vendor"
)

type ProcessPriceService struct {
	priceRepo     *ProcessPriceRepo
	vendorService interface {
		Get(id uint) (*vendor.Vendor, error)
	}
	rdb *redis.Client
	ctx context.Context
}

func NewProcessPriceService(
	priceRepo *ProcessPriceRepo,
	vendorService interface{ Get(id uint) (*vendor.Vendor, error) },
	rdb *redis.Client,
) *ProcessPriceService {
	return &ProcessPriceService{
		priceRepo:     priceRepo,
		vendorService: vendorService,
		rdb:           rdb,
		ctx:           context.Background(),
	}
}

// 厂商报价
func (s *ProcessPriceService) Quote(vendorID, processID uint, price float64) error {
	// 1. 写入数据库
	processPrice := &ProcessPrice{
		VendorID:  vendorID,
		ProcessID: processID,
		Price:     price,
		QuotedAt:  time.Now(),
	}
	if err := s.priceRepo.Create(processPrice); err != nil {
		return err
	}

	// 2. 获取厂商名称
	v, err := s.vendorService.Get(vendorID)
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
	if err := s.compareAndUpdateMin(processID, priceData); err != nil {
		return err
	}

	// 更新最高价
	if err := s.compareAndUpdateMax(processID, priceData); err != nil {
		return err
	}

	return nil
}

// 比对并更新最低价
func (s *ProcessPriceService) compareAndUpdateMin(processID uint, newPrice PriceData) error {
	key := fmt.Sprintf("price:process:%d:min", processID)

	// 读取当前最低价
	val, err := s.rdb.Get(s.ctx, key).Result()
	if err == redis.Nil {
		// Redis 中没有,直接写入
		return s.setPrice(key, newPrice)
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
		return s.setPrice(key, newPrice)
	}

	return nil
}

// 比对并更新最高价
func (s *ProcessPriceService) compareAndUpdateMax(processID uint, newPrice PriceData) error {
	key := fmt.Sprintf("price:process:%d:max", processID)

	// 读取当前最高价
	val, err := s.rdb.Get(s.ctx, key).Result()
	if err == redis.Nil {
		// Redis 中没有,直接写入
		return s.setPrice(key, newPrice)
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
		return s.setPrice(key, newPrice)
	}

	return nil
}

// 写入 Redis
func (s *ProcessPriceService) setPrice(key string, priceData PriceData) error {
	data, err := json.Marshal(priceData)
	if err != nil {
		return err
	}
	return s.rdb.Set(s.ctx, key, data, 0).Err() // 无 TTL
}

// 获取最低价
func (s *ProcessPriceService) GetMinPrice(processID uint) (*PriceData, error) {
	key := fmt.Sprintf("price:process:%d:min", processID)

	// 尝试从 Redis 读取
	val, err := s.rdb.Get(s.ctx, key).Result()
	if err == redis.Nil {
		// Redis 未命中,从数据库加载
		return s.loadMinPriceFromDB(processID)
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

// 获取最高价
func (s *ProcessPriceService) GetMaxPrice(processID uint) (*PriceData, error) {
	key := fmt.Sprintf("price:process:%d:max", processID)

	// 尝试从 Redis 读取
	val, err := s.rdb.Get(s.ctx, key).Result()
	if err == redis.Nil {
		// Redis 未命中,从数据库加载
		return s.loadMaxPriceFromDB(processID)
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

// 从数据库加载最低价并写入 Redis
func (s *ProcessPriceService) loadMinPriceFromDB(processID uint) (*PriceData, error) {
	price, err := s.priceRepo.GetMinPrice(processID)
	if err != nil {
		return nil, err
	}

	v, err := s.vendorService.Get(price.VendorID)
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
	if err := s.setPrice(key, *priceData); err != nil {
		return nil, err
	}

	return priceData, nil
}

// 从数据库加载最高价并写入 Redis
func (s *ProcessPriceService) loadMaxPriceFromDB(processID uint) (*PriceData, error) {
	price, err := s.priceRepo.GetMaxPrice(processID)
	if err != nil {
		return nil, err
	}

	v, err := s.vendorService.Get(price.VendorID)
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
	if err := s.setPrice(key, *priceData); err != nil {
		return nil, err
	}

	return priceData, nil
}

// 获取价格历史
func (s *ProcessPriceService) GetHistory(processID uint) ([]ProcessPrice, error) {
	return s.priceRepo.GetHistory(processID)
}