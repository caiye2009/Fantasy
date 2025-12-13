package application

import (
	"context"

	"back/internal/product/domain"
	"back/internal/product/infra"
	pricingDomain "back/internal/pricing/domain"
)

// MaterialServiceInterface Material 服务接口
type MaterialServiceInterface interface {
	Get(ctx context.Context, id uint) (*MaterialInfo, error)
}

// ProcessServiceInterface Process 服务接口
type ProcessServiceInterface interface {
	Get(ctx context.Context, id uint) (*ProcessInfo, error)
}

// MaterialPriceServiceInterface Material 价格服务接口
type MaterialPriceServiceInterface interface {
	GetMinPrice(ctx context.Context, materialID uint) (*pricingDomain.PriceData, error)
	GetMaxPrice(ctx context.Context, materialID uint) (*pricingDomain.PriceData, error)
}

// ProcessPriceServiceInterface Process 价格服务接口
type ProcessPriceServiceInterface interface {
	GetMinPrice(ctx context.Context, processID uint) (*pricingDomain.PriceData, error)
	GetMaxPrice(ctx context.Context, processID uint) (*pricingDomain.PriceData, error)
}

// MaterialInfo Material 信息
type MaterialInfo struct {
	ID   uint
	Name string
}

// ProcessInfo Process 信息
type ProcessInfo struct {
	ID   uint
	Name string
}

// CostCalculator 成本计算器
type CostCalculator struct {
	productRepo      *infra.ProductRepo
	materialService  MaterialServiceInterface
	processService   ProcessServiceInterface
	materialPriceSvc MaterialPriceServiceInterface
	processPriceSvc  ProcessPriceServiceInterface
}

// NewCostCalculator 创建成本计算器
func NewCostCalculator(
	productRepo *infra.ProductRepo,
	materialService MaterialServiceInterface,
	processService ProcessServiceInterface,
	materialPriceSvc MaterialPriceServiceInterface,
	processPriceSvc ProcessPriceServiceInterface,
) *CostCalculator {
	return &CostCalculator{
		productRepo:      productRepo,
		materialService:  materialService,
		processService:   processService,
		materialPriceSvc: materialPriceSvc,
		processPriceSvc:  processPriceSvc,
	}
}

// Calculate 计算产品成本
func (c *CostCalculator) Calculate(ctx context.Context, req *CalculateCostRequest) (*domain.CostResult, error) {
	// 1. 查询产品
	product, err := c.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}
	
	// 2. 验证产品配置
	if err := product.Validate(); err != nil {
		return nil, err
	}
	
	// 3. 计算原料成本
	materialCost := 0.0
	materialItems := []domain.MaterialCostItem{}
	
	for _, m := range product.Materials {
		// 获取价格
		var priceData *pricingDomain.PriceData
		if req.UseMinPrice {
			priceData, err = c.materialPriceSvc.GetMinPrice(ctx, m.MaterialID)
		} else {
			priceData, err = c.materialPriceSvc.GetMaxPrice(ctx, m.MaterialID)
		}
		if err != nil {
			return nil, err
		}
		
		// 获取材料信息
		mat, err := c.materialService.Get(ctx, m.MaterialID)
		if err != nil {
			return nil, err
		}
		
		// 计算成本
		cost := priceData.Price * m.Ratio
		materialCost += cost
		
		materialItems = append(materialItems, domain.MaterialCostItem{
			MaterialID:   m.MaterialID,
			MaterialName: mat.Name,
			Ratio:        m.Ratio,
			Price:        priceData.Price,
			Cost:         cost,
		})
	}
	
	// 4. 计算工艺成本
	processCost := 0.0
	processItems := []domain.ProcessCostItem{}
	
	for _, p := range product.Processes {
		// 获取价格
		var priceData *pricingDomain.PriceData
		if req.UseMinPrice {
			priceData, err = c.processPriceSvc.GetMinPrice(ctx, p.ProcessID)
		} else {
			priceData, err = c.processPriceSvc.GetMaxPrice(ctx, p.ProcessID)
		}
		if err != nil {
			return nil, err
		}
		
		// 获取工艺信息
		proc, err := c.processService.Get(ctx, p.ProcessID)
		if err != nil {
			return nil, err
		}
		
		processCost += priceData.Price
		
		processItems = append(processItems, domain.ProcessCostItem{
			ProcessID:   p.ProcessID,
			ProcessName: proc.Name,
			Price:       priceData.Price,
			Cost:        priceData.Price,
		})
	}
	
	// 5. 计算总成本
	unitCost := materialCost + processCost
	totalCost := unitCost * req.Quantity
	
	return &domain.CostResult{
		ProductID:    product.ID,
		ProductName:  product.Name,
		Quantity:     req.Quantity,
		MaterialCost: materialCost,
		ProcessCost:  processCost,
		UnitCost:     unitCost,
		TotalCost:    totalCost,
		Breakdown: &domain.CostBreakdown{
			Materials: materialItems,
			Processes: processItems,
		},
	}, nil
}