package price

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"back/internal/vendor"
)

type MaterialPriceService struct {
	priceRepo     *MaterialPriceRepo
	vendorService interface {
		Get(id uint) (*vendor.Vendor, error)
	}
	rdb *redis.Client
	ctx context.Context
}

func NewMaterialPriceService(
	priceRepo *MaterialPriceRepo,
	vendorService interface{ Get(id uint) (*vendor.Vendor, error) },
	rdb *redis.Client,
) *MaterialPriceService {
	return &MaterialPriceService{
		priceRepo:     priceRepo,
		vendorService: vendorService,
		rdb:           rdb,
		ctx:           context.Background(),
	}
}

// 厂商报价
func (s *MaterialPriceService) Quote(vendorID, materialID uint, price float64) error {
	// 1. 写入数据库
	materialPrice := &MaterialPrice{
		VendorID:   vendorID,
		MaterialID: materialID,
		Price:      price,
		QuotedAt:   time.Now(),
	}
	if err := s.priceRepo.Create(materialPrice); err != nil {
		return err
	}

	// 2. 获取厂商名称
	v, err := s.vendorService.Get(vendorID)
	if err != nil {
		return err
	}

	// 3. 更新 Redis (比对最低/最高)
	priceData := PriceData{
		ID:         materialPrice.ID,
		Price:      price,
		VendorID:   vendorID,
		VendorName: v.Name,
		QuotedAt:   materialPrice.QuotedAt,
	}

	// 更新最低价
	if err := s.compareAndUpdateMin(materialID, priceData); err != nil {
		return err
	}

	// 更新最高价
	if err := s.compareAndUpdateMax(materialID, priceData); err != nil {
		return err
	}

	return nil
}

// 比对并更新最低价
func (s *MaterialPriceService) compareAndUpdateMin(materialID uint, newPrice PriceData) error {
	key := fmt.Sprintf("price:material:%d:min", materialID)

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
func (s *MaterialPriceService) compareAndUpdateMax(materialID uint, newPrice PriceData) error {
	key := fmt.Sprintf("price:material:%d:max", materialID)

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
func (s *MaterialPriceService) setPrice(key string, priceData PriceData) error {
	data, err := json.Marshal(priceData)
	if err != nil {
		return err
	}
	return s.rdb.Set(s.ctx, key, data, 0).Err() // 无 TTL
}

// 获取最低价
func (s *MaterialPriceService) GetMinPrice(materialID uint) (*PriceData, error) {
	key := fmt.Sprintf("price:material:%d:min", materialID)

	// 尝试从 Redis 读取
	val, err := s.rdb.Get(s.ctx, key).Result()
	if err == redis.Nil {
		// Redis 未命中,从数据库加载
		return s.loadMinPriceFromDB(materialID)
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
func (s *MaterialPriceService) GetMaxPrice(materialID uint) (*PriceData, error) {
	key := fmt.Sprintf("price:material:%d:max", materialID)

	// 尝试从 Redis 读取
	val, err := s.rdb.Get(s.ctx, key).Result()
	if err == redis.Nil {
		// Redis 未命中,从数据库加载
		return s.loadMaxPriceFromDB(materialID)
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
func (s *MaterialPriceService) loadMinPriceFromDB(materialID uint) (*PriceData, error) {
	price, err := s.priceRepo.GetMinPrice(materialID)
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
	key := fmt.Sprintf("price:material:%d:min", materialID)
	if err := s.setPrice(key, *priceData); err != nil {
		return nil, err
	}

	return priceData, nil
}

// 从数据库加载最高价并写入 Redis
func (s *MaterialPriceService) loadMaxPriceFromDB(materialID uint) (*PriceData, error) {
	price, err := s.priceRepo.GetMaxPrice(materialID)
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
	key := fmt.Sprintf("price:material:%d:max", materialID)
	if err := s.setPrice(key, *priceData); err != nil {
		return nil, err
	}

	return priceData, nil
}

// 获取价格历史
func (s *MaterialPriceService) GetHistory(materialID uint) ([]MaterialPrice, error) {
	return s.priceRepo.GetHistory(materialID)
}