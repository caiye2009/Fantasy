package application

import (
	"context"

	"back/internal/product/domain"
	"back/internal/product/infra"
	pricingApp "back/internal/pricing/application"
)

// ProductPriceService 产品价格服务
type ProductPriceService struct {
	productRepo      *infra.ProductRepo
	materialPriceSvc *pricingApp.MaterialPriceService
	processPriceSvc  *pricingApp.ProcessPriceService
}

// NewProductPriceService 创建产品价格服务
func NewProductPriceService(
	productRepo *infra.ProductRepo,
	materialPriceSvc *pricingApp.MaterialPriceService,
	processPriceSvc *pricingApp.ProcessPriceService,
) *ProductPriceService {
	return &ProductPriceService{
		productRepo:      productRepo,
		materialPriceSvc: materialPriceSvc,
		processPriceSvc:  processPriceSvc,
	}
}

// ProductPriceResponse 产品价格响应
type ProductPriceResponse struct {
	CurrentPrice     float64 `json:"current_price"`
	HistoricalHigh   float64 `json:"historical_high"`
	HistoricalLow    float64 `json:"historical_low"`
}

// GetPrice 获取产品价格（当前价、历史最高价、历史最低价）
func (s *ProductPriceService) GetPrice(ctx context.Context, productID uint) (*ProductPriceResponse, error) {
	// 1. 查询产品
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 2. 验证产品配置
	if err := product.Validate(); err != nil {
		return nil, err
	}

	// 3. 计算当前价格
	currentPrice, err := s.calculateCurrentPrice(ctx, product)
	if err != nil {
		return nil, err
	}

	// 4. 计算历史最高价
	highPrice, err := s.calculateHighPrice(ctx, product)
	if err != nil {
		return nil, err
	}

	// 5. 计算历史最低价
	lowPrice, err := s.calculateLowPrice(ctx, product)
	if err != nil {
		return nil, err
	}

	return &ProductPriceResponse{
		CurrentPrice:   currentPrice,
		HistoricalHigh: highPrice,
		HistoricalLow:  lowPrice,
	}, nil
}

// calculateCurrentPrice 计算当前价格（基于最新报价）
func (s *ProductPriceService) calculateCurrentPrice(ctx context.Context, product *domain.Product) (float64, error) {
	// 计算材料成本
	materialCost := 0.0
	for _, m := range product.Materials {
		priceData, err := s.materialPriceSvc.GetCurrentPrice(ctx, m.MaterialID)
		if err != nil {
			return 0, err
		}
		materialCost += priceData.Price * m.Ratio
	}

	// 计算工艺成本
	processCost := 0.0
	for _, p := range product.Processes {
		priceData, err := s.processPriceSvc.GetCurrentPrice(ctx, p.ProcessID)
		if err != nil {
			return 0, err
		}
		processCost += priceData.Price
	}

	return materialCost + processCost, nil
}

// calculateHighPrice 计算历史最高价
func (s *ProductPriceService) calculateHighPrice(ctx context.Context, product *domain.Product) (float64, error) {
	// 计算材料最高成本
	materialCost := 0.0
	for _, m := range product.Materials {
		priceData, err := s.materialPriceSvc.GetMaxPrice(ctx, m.MaterialID)
		if err != nil {
			return 0, err
		}
		materialCost += priceData.Price * m.Ratio
	}

	// 计算工艺最高成本
	processCost := 0.0
	for _, p := range product.Processes {
		priceData, err := s.processPriceSvc.GetMaxPrice(ctx, p.ProcessID)
		if err != nil {
			return 0, err
		}
		processCost += priceData.Price
	}

	return materialCost + processCost, nil
}

// calculateLowPrice 计算历史最低价
func (s *ProductPriceService) calculateLowPrice(ctx context.Context, product *domain.Product) (float64, error) {
	// 计算材料最低成本
	materialCost := 0.0
	for _, m := range product.Materials {
		priceData, err := s.materialPriceSvc.GetMinPrice(ctx, m.MaterialID)
		if err != nil {
			return 0, err
		}
		materialCost += priceData.Price * m.Ratio
	}

	// 计算工艺最低成本
	processCost := 0.0
	for _, p := range product.Processes {
		priceData, err := s.processPriceSvc.GetMinPrice(ctx, p.ProcessID)
		if err != nil {
			return 0, err
		}
		processCost += priceData.Price
	}

	return materialCost + processCost, nil
}
