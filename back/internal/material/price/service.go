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

type MaterialPriceService struct {
	db            *gorm.DB
	vendorService interface {
		Get(ctx context.Context, id uint) (*vendor.Vendor, error)
	}
	rdb *redis.Client
}

func NewMaterialPriceService(
	db *gorm.DB,
	vendorService interface{ Get(ctx context.Context, id uint) (*vendor.Vendor, error) },
	rdb *redis.Client,
) *MaterialPriceService {
	return &MaterialPriceService{
		db:            db,
		vendorService: vendorService,
		rdb:           rdb,
	}
}

// Quote 厂商报价
func (s *MaterialPriceService) Quote(ctx context.Context, vendorID, materialID uint, price float64) error {
	// 1. 写入数据库
	materialPrice := &MaterialPrice{
		VendorID:   vendorID,
		MaterialID: materialID,
		Price:      price,
		QuotedAt:   time.Now(),
	}
	
	priceRepo := repo.NewRepo[MaterialPrice](s.db)
	if err := priceRepo.Create(ctx, materialPrice); err != nil {
		return err
	}

	// 2. 获取厂商名称
	v, err := s.vendorService.Get(ctx, vendorID)
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
	if err := s.compareAndUpdateMin(ctx, materialID, priceData); err != nil {
		return err
	}

	// 更新最高价
	if err := s.compareAndUpdateMax(ctx, materialID, priceData); err != nil {
		return err
	}

	return nil
}

// compareAndUpdateMin 比对并更新最低价
func (s *MaterialPriceService) compareAndUpdateMin(ctx context.Context, materialID uint, newPrice PriceData) error {
	key := fmt.Sprintf("price:material:%d:min", materialID)

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
func (s *MaterialPriceService) compareAndUpdateMax(ctx context.Context, materialID uint, newPrice PriceData) error {
	key := fmt.Sprintf("price:material:%d:max", materialID)

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
func (s *MaterialPriceService) setPrice(ctx context.Context, key string, priceData PriceData) error {
	data, err := json.Marshal(priceData)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, key, data, 0).Err() // 无 TTL
}

// GetMinPrice 获取最低价
func (s *MaterialPriceService) GetMinPrice(materialID uint) (*PriceData, error) {
	ctx := context.Background()
	key := fmt.Sprintf("price:material:%d:min", materialID)

	// 尝试从 Redis 读取
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Redis 未命中,从数据库加载
		return s.loadMinPriceFromDB(ctx, materialID)
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
func (s *MaterialPriceService) GetMaxPrice(materialID uint) (*PriceData, error) {
	ctx := context.Background()
	key := fmt.Sprintf("price:material:%d:max", materialID)

	// 尝试从 Redis 读取
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Redis 未命中,从数据库加载
		return s.loadMaxPriceFromDB(ctx, materialID)
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
func (s *MaterialPriceService) loadMinPriceFromDB(ctx context.Context, materialID uint) (*PriceData, error) {
	// 查询最低价
	var price MaterialPrice
	err := s.db.WithContext(ctx).
		Where("material_id = ?", materialID).
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
	key := fmt.Sprintf("price:material:%d:min", materialID)
	if err := s.setPrice(ctx, key, *priceData); err != nil {
		return nil, err
	}

	return priceData, nil
}

// loadMaxPriceFromDB 从数据库加载最高价并写入 Redis
func (s *MaterialPriceService) loadMaxPriceFromDB(ctx context.Context, materialID uint) (*PriceData, error) {
	// 查询最高价
	var price MaterialPrice
	err := s.db.WithContext(ctx).
		Where("material_id = ?", materialID).
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
	key := fmt.Sprintf("price:material:%d:max", materialID)
	if err := s.setPrice(ctx, key, *priceData); err != nil {
		return nil, err
	}

	return priceData, nil
}

// GetHistory 获取价格历史
func (s *MaterialPriceService) GetHistory(ctx context.Context, materialID uint) ([]MaterialPrice, error) {
	var prices []MaterialPrice
	err := s.db.WithContext(ctx).
		Where("material_id = ?", materialID).
		Order("quoted_at DESC").
		Find(&prices).Error
	
	return prices, err
}