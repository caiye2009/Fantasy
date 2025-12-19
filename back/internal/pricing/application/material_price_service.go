package application

import (
	"context"
	"time"

	"back/internal/pricing/domain"
	"back/internal/pricing/infra"
	materialApp "back/internal/material/application"
	supplierApp "back/internal/supplier/application"
)

// MaterialPriceService Material 价格服务
type MaterialPriceService struct {
	repo            *infra.SupplierPriceRepo
	cache           domain.PriceCache
	materialService *materialApp.MaterialService
	supplierService *supplierApp.SupplierService
}

// NewMaterialPriceService 创建 Material 价格服务
func NewMaterialPriceService(
	repo *infra.SupplierPriceRepo,
	cache domain.PriceCache,
	materialService *materialApp.MaterialService,
	supplierService *supplierApp.SupplierService,
) *MaterialPriceService {
	return &MaterialPriceService{
		repo:            repo,
		cache:           cache,
		materialService: materialService,
		supplierService: supplierService,
	}
}

// Quote 供应商报价
func (s *MaterialPriceService) Quote(ctx context.Context, req *QuoteRequest) error {
	// 1. 验证 Material 是否存在
	exists, err := s.materialService.Exists(ctx, req.TargetID)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrTargetNotFound
	}
	
	// 2. 创建价格记录
	price := &domain.SupplierPrice{
		TargetType: domain.TargetTypeMaterial,
		TargetID:   req.TargetID,
		SupplierID: req.SupplierID,
		Price:      req.Price,
		QuotedAt:   time.Now(),
	}
	
	// 3. 领域验证
	if err := price.Validate(); err != nil {
		return err
	}
	
	// 4. 保存到数据库
	if err := s.repo.Save(ctx, price); err != nil {
		return err
	}
	
	// 5. 获取供应商信息
	supplier, err := s.supplierService.GetSupplierInfo(ctx, req.SupplierID)
	if err != nil {
		return err
	}
	
	// 6. 更新缓存（最低价和最高价）
	minPrice, _ := s.repo.FindMinPrice(ctx, domain.TargetTypeMaterial, req.TargetID)
	if minPrice != nil {
		priceData := &domain.PriceData{
			Price:        minPrice.Price,
			SupplierID:   minPrice.SupplierID,
			SupplierName: supplier.Name,
			QuotedAt:     minPrice.QuotedAt,
		}
		s.cache.UpdateMin(ctx, "material", req.TargetID, priceData)
	}
	
	maxPrice, _ := s.repo.FindMaxPrice(ctx, domain.TargetTypeMaterial, req.TargetID)
	if maxPrice != nil {
		priceData := &domain.PriceData{
			Price:        maxPrice.Price,
			SupplierID:   maxPrice.SupplierID,
			SupplierName: supplier.Name,
			QuotedAt:     maxPrice.QuotedAt,
		}
		s.cache.UpdateMax(ctx, "material", req.TargetID, priceData)
	}
	
	return nil
}

// GetMinPrice 获取最低价
func (s *MaterialPriceService) GetMinPrice(ctx context.Context, materialID uint) (*domain.PriceData, error) {
	// 1. 尝试从缓存获取
	cachedPrice, err := s.cache.GetMin(ctx, "material", materialID)
	if err == nil {
		return cachedPrice, nil
	}
	
	// 2. 缓存未命中，从数据库查询
	price, err := s.repo.FindMinPrice(ctx, domain.TargetTypeMaterial, materialID)
	if err != nil {
		return nil, err
	}
	
	// 3. 获取供应商信息
	supplier, err := s.supplierService.GetSupplierInfo(ctx, price.SupplierID)
	if err != nil {
		return nil, err
	}
	
	priceData := &domain.PriceData{
		Price:        price.Price,
		SupplierID:   price.SupplierID,
		SupplierName: supplier.Name,
		QuotedAt:     price.QuotedAt,
	}
	
	// 4. 写入缓存
	s.cache.SetMin(ctx, "material", materialID, priceData)
	
	return priceData, nil
}

// GetMaxPrice 获取最高价
func (s *MaterialPriceService) GetMaxPrice(ctx context.Context, materialID uint) (*domain.PriceData, error) {
	// 1. 尝试从缓存获取
	cachedPrice, err := s.cache.GetMax(ctx, "material", materialID)
	if err == nil {
		return cachedPrice, nil
	}
	
	// 2. 缓存未命中，从数据库查询
	price, err := s.repo.FindMaxPrice(ctx, domain.TargetTypeMaterial, materialID)
	if err != nil {
		return nil, err
	}
	
	// 3. 获取供应商信息
	supplier, err := s.supplierService.GetSupplierInfo(ctx, price.SupplierID)
	if err != nil {
		return nil, err
	}
	
	priceData := &domain.PriceData{
		Price:        price.Price,
		SupplierID:   price.SupplierID,
		SupplierName: supplier.Name,
		QuotedAt:     price.QuotedAt,
	}
	
	// 4. 写入缓存
	s.cache.SetMax(ctx, "material", materialID, priceData)
	
	return priceData, nil
}

// GetCurrentPrice 获取当前价格（最新报价）
func (s *MaterialPriceService) GetCurrentPrice(ctx context.Context, materialID uint) (*domain.PriceData, error) {
	// 获取最新的一条报价
	prices, err := s.repo.FindHistory(ctx, domain.TargetTypeMaterial, materialID, 1)
	if err != nil {
		return nil, err
	}

	if len(prices) == 0 {
		return nil, domain.ErrPriceNotFound
	}

	// 获取供应商信息
	supplier, err := s.supplierService.GetSupplierInfo(ctx, prices[0].SupplierID)
	if err != nil {
		return nil, err
	}

	return &domain.PriceData{
		Price:        prices[0].Price,
		SupplierID:   prices[0].SupplierID,
		SupplierName: supplier.Name,
		QuotedAt:     prices[0].QuotedAt,
	}, nil
}

// GetHistory 获取报价历史
func (s *MaterialPriceService) GetHistory(ctx context.Context, materialID uint, limit int) ([]*domain.PriceData, error) {
	prices, err := s.repo.FindHistory(ctx, domain.TargetTypeMaterial, materialID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.PriceData, len(prices))
	for i, p := range prices {
		supplier, err := s.supplierService.GetSupplierInfo(ctx, p.SupplierID)
		if err != nil {
			return nil, err
		}

		result[i] = &domain.PriceData{
			Price:        p.Price,
			SupplierID:   p.SupplierID,
			SupplierName: supplier.Name,
			QuotedAt:     p.QuotedAt,
		}
	}

	return result, nil
}