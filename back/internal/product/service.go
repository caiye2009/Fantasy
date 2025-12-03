package product

import (
	"errors"
	"math"
	"strconv"
	
	"back/pkg/es"
	materialPrice "back/internal/material/price"
	processPrice "back/internal/process/price"
	"back/internal/material"
	"back/internal/process"
)

type ProductService struct {
	productRepo          *ProductRepo
	materialPriceService interface {
		GetMinPrice(materialID uint) (*materialPrice.PriceData, error)
		GetMaxPrice(materialID uint) (*materialPrice.PriceData, error)
	}
	processPriceService interface {
		GetMinPrice(processID uint) (*processPrice.PriceData, error)
		GetMaxPrice(processID uint) (*processPrice.PriceData, error)
	}
	materialService interface {
		Get(id uint) (*material.Material, error)
	}
	processService interface {
		Get(id uint) (*process.Process, error)
	}
	esSync *es.ESSync
}

func NewProductService(
	productRepo *ProductRepo,
	materialPriceService interface {
		GetMinPrice(materialID uint) (*materialPrice.PriceData, error)
		GetMaxPrice(materialID uint) (*materialPrice.PriceData, error)
	},
	processPriceService interface {
		GetMinPrice(processID uint) (*processPrice.PriceData, error)
		GetMaxPrice(processID uint) (*processPrice.PriceData, error)
	},
	materialService interface {
		Get(id uint) (*material.Material, error)
	},
	processService interface {
		Get(id uint) (*process.Process, error)
	},
	esSync *es.ESSync,
) *ProductService {
	return &ProductService{
		productRepo:          productRepo,
		materialPriceService: materialPriceService,
		processPriceService:  processPriceService,
		materialService:      materialService,
		processService:       processService,
		esSync:               esSync,
	}
}

func (s *ProductService) Create(p *Product) error {
	// 默认状态为 draft
	if p.Status == "" {
		p.Status = "draft"
	}
	
	// 暂不验证,允许创建不完整的产品
	if err := s.productRepo.Create(p); err != nil {
		return err
	}

	// 异步同步到 ES
	s.esSync.Index(p)

	return nil
}

func (s *ProductService) Get(id uint) (*Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *ProductService) List() ([]Product, error) {
	return s.productRepo.List()
}

func (s *ProductService) Update(id uint, data map[string]interface{}) error {
	// 先获取产品
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 更新数据库
	if err := s.productRepo.Update(id, data); err != nil {
		return err
	}

	// 更新内存中的 product 对象
	if name, ok := data["name"].(string); ok {
		product.Name = name
	}
	if status, ok := data["status"].(string); ok {
		product.Status = status
	}
	if materials, ok := data["materials"].([]MaterialConfig); ok {
		product.Materials = materials
	}
	if processes, ok := data["processes"].([]ProcessConfig); ok {
		product.Processes = processes
	}

	// 异步同步到 ES
	s.esSync.Update(product)

	return nil
}

func (s *ProductService) Delete(id uint) error {
	if err := s.productRepo.Delete(id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("products", strconv.Itoa(int(id)))

	return nil
}

// CalculateCost 计算产品成本
func (s *ProductService) CalculateCost(productID uint, quantity float64, useMinPrice bool) (*CostResult, error) {
	// 1. 获取产品
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	// 2. 检查配置
	if len(product.Materials) == 0 {
		return nil, errors.New("产品未配置原料")
	}
	if len(product.Processes) == 0 {
		return nil, errors.New("产品未配置工艺")
	}

	// 3. 计算原料成本
	materialCost := 0.0
	materialItems := []MaterialCostItem{}

	for _, m := range product.Materials {
		var priceData *materialPrice.PriceData
		if useMinPrice {
			priceData, err = s.materialPriceService.GetMinPrice(m.MaterialID)
		} else {
			priceData, err = s.materialPriceService.GetMaxPrice(m.MaterialID)
		}
		if err != nil {
			return nil, err
		}

		// 获取原料名称
		mat, err := s.materialService.Get(m.MaterialID)
		if err != nil {
			return nil, err
		}

		cost := priceData.Price * m.Ratio
		materialCost += cost

		materialItems = append(materialItems, MaterialCostItem{
			MaterialID:   m.MaterialID,
			MaterialName: mat.Name,
			Ratio:        m.Ratio,
			Price:        priceData.Price,
			Cost:         cost,
		})
	}

	// 4. 计算工艺成本
	processCost := 0.0
	processItems := []ProcessCostItem{}

	for _, p := range product.Processes {
		var priceData *processPrice.PriceData
		if useMinPrice {
			priceData, err = s.processPriceService.GetMinPrice(p.ProcessID)
		} else {
			priceData, err = s.processPriceService.GetMaxPrice(p.ProcessID)
		}
		if err != nil {
			return nil, err
		}

		// 获取工艺名称
		proc, err := s.processService.Get(p.ProcessID)
		if err != nil {
			return nil, err
		}

		processCost += priceData.Price

		processItems = append(processItems, ProcessCostItem{
			ProcessID:   p.ProcessID,
			ProcessName: proc.Name,
			Price:       priceData.Price,
			Cost:        priceData.Price,
		})
	}

	// 5. 计算总成本
	unitCost := materialCost + processCost
	totalCost := unitCost * quantity

	return &CostResult{
		MaterialCost: materialCost,
		ProcessCost:  processCost,
		UnitCost:     unitCost,
		TotalCost:    totalCost,
		Breakdown: &CostBreakdown{
			Materials: materialItems,
			Processes: processItems,
		},
	}, nil
}

// ValidateForSubmit 验证产品是否可以提交审核 (预留)
func (s *ProductService) ValidateForSubmit(product *Product) error {
	// 验证名称
	if product.Name == "" {
		return errors.New("产品名称不能为空")
	}

	// 验证原料
	if len(product.Materials) == 0 {
		return errors.New("至少需要一个原料")
	}

	// 验证占比总和
	sum := 0.0
	for _, m := range product.Materials {
		if m.Ratio <= 0 {
			return errors.New("原料占比必须大于0")
		}
		sum += m.Ratio
	}
	if math.Abs(sum-1.0) > 0.0001 {
		return errors.New("原料占比总和必须为1")
	}

	// 验证工艺
	if len(product.Processes) == 0 {
		return errors.New("至少需要一个工艺")
	}

	return nil
}