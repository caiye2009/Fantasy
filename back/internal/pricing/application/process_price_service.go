package application

import (
	"context"
	"time"

	"back/internal/pricing/domain"
	"back/internal/pricing/infra"
	processApp "back/internal/process/application"
	supplierApp "back/internal/supplier/application"
)

// ProcessPriceService Process 价格服务
type ProcessPriceService struct {
	repo            *infra.SupplierPriceRepo
	cache           domain.PriceCache
	processService  *processApp.ProcessService
	supplierService *supplierApp.SupplierService
}

// NewProcessPriceService 创建 Process 价格服务
func NewProcessPriceService(
	repo *infra.SupplierPriceRepo,
	cache domain.PriceCache,
	processService *processApp.ProcessService,
	supplierService *supplierApp.SupplierService,
) *ProcessPriceService {
	return &ProcessPriceService{
		repo:            repo,
		cache:           cache,
		processService:  processService,
		supplierService: supplierService,
	}
}

// Quote 供应商报价
func (s *ProcessPriceService) Quote(ctx context.Context, req *QuoteRequest) error {
	// 1. 验证 Process 是否存在
	exists, err := s.processService.Exists(ctx, req.TargetID)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrTargetNotFound
	}
	
	// 2. 创建价格记录
	price := &domain.SupplierPrice{
		TargetType: domain.TargetTypeProcess,
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
	minPrice, _ := s.repo.FindMinPrice(ctx, domain.TargetTypeProcess, req.TargetID)
	if minPrice != nil {
		priceData := &domain.PriceData{
			Price:        minPrice.Price,
			SupplierID:   minPrice.SupplierID,
			SupplierName: supplier.Name,
			QuotedAt:     minPrice.QuotedAt,
		}
		s.cache.UpdateMin(ctx, "process", req.TargetID, priceData)
	}
	
	maxPrice, _ := s.repo.FindMaxPrice(ctx, domain.TargetTypeProcess, req.TargetID)
	if maxPrice != nil {
		priceData := &domain.PriceData{
			Price:        maxPrice.Price,
			SupplierID:   maxPrice.SupplierID,
			SupplierName: supplier.Name,
			QuotedAt:     maxPrice.QuotedAt,
		}
		s.cache.UpdateMax(ctx, "process", req.TargetID, priceData)
	}
	
	return nil
}

// GetMinPrice 获取最低价
func (s *ProcessPriceService) GetMinPrice(ctx context.Context, processID uint) (*domain.PriceData, error) {
	// 1. 尝试从缓存获取
	cachedPrice, err := s.cache.GetMin(ctx, "process", processID)
	if err == nil {
		return cachedPrice, nil
	}
	
	// 2. 缓存未命中，从数据库查询
	price, err := s.repo.FindMinPrice(ctx, domain.TargetTypeProcess, processID)
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
	s.cache.SetMin(ctx, "process", processID, priceData)
	
	return priceData, nil
}

// GetMaxPrice 获取最高价
func (s *ProcessPriceService) GetMaxPrice(ctx context.Context, processID uint) (*domain.PriceData, error) {
	// 1. 尝试从缓存获取
	cachedPrice, err := s.cache.GetMax(ctx, "process", processID)
	if err == nil {
		return cachedPrice, nil
	}
	
	// 2. 缓存未命中，从数据库查询
	price, err := s.repo.FindMaxPrice(ctx, domain.TargetTypeProcess, processID)
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
	s.cache.SetMax(ctx, "process", processID, priceData)
	
	return priceData, nil
}

// GetCurrentPrice 获取当前价格（最新报价）
func (s *ProcessPriceService) GetCurrentPrice(ctx context.Context, processID uint) (*domain.PriceData, error) {
	// 获取最新的一条报价
	prices, err := s.repo.FindHistory(ctx, domain.TargetTypeProcess, processID, 1)
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
func (s *ProcessPriceService) GetHistory(ctx context.Context, processID uint, limit int) ([]*domain.PriceData, error) {
	prices, err := s.repo.FindHistory(ctx, domain.TargetTypeProcess, processID, limit)
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