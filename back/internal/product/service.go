package product

import (
	"context"
	"errors"
	"math"
	"strconv"
	
	"back/pkg/es"
	"back/pkg/repo"
	materialPrice "back/internal/material/price"
	processPrice "back/internal/process/price"
	"back/internal/material"
	"back/internal/process"
	"gorm.io/gorm"
)

type ProductService struct {
	db                   *gorm.DB
	materialPriceService interface {
		GetMinPrice(materialID uint) (*materialPrice.PriceData, error)
		GetMaxPrice(materialID uint) (*materialPrice.PriceData, error)
	}
	processPriceService interface {
		GetMinPrice(processID uint) (*processPrice.PriceData, error)
		GetMaxPrice(processID uint) (*processPrice.PriceData, error)
	}
	materialService interface {
		Get(ctx context.Context, id uint) (*material.Material, error)
	}
	processService interface {
		Get(ctx context.Context, id uint) (*process.Process, error)
	}
	esSync *es.ESSync
}

func NewProductService(
	db *gorm.DB,
	materialPriceService interface {
		GetMinPrice(materialID uint) (*materialPrice.PriceData, error)
		GetMaxPrice(materialID uint) (*materialPrice.PriceData, error)
	},
	processPriceService interface {
		GetMinPrice(processID uint) (*processPrice.PriceData, error)
		GetMaxPrice(processID uint) (*processPrice.PriceData, error)
	},
	materialService interface {
		Get(ctx context.Context, id uint) (*material.Material, error)
	},
	processService interface {
		Get(ctx context.Context, id uint) (*process.Process, error)
	},
	esSync *es.ESSync,
) *ProductService {
	return &ProductService{
		db:                   db,
		materialPriceService: materialPriceService,
		processPriceService:  processPriceService,
		materialService:      materialService,
		processService:       processService,
		esSync:               esSync,
	}
}

func (s *ProductService) Create(ctx context.Context, p *Product) error {
	if p.Status == "" {
		p.Status = "draft"
	}
	
	productRepo := repo.NewRepo[Product](s.db)
	if err := productRepo.Create(ctx, p); err != nil {
		return err
	}

	s.esSync.Index(p)
	return nil
}

func (s *ProductService) Get(ctx context.Context, id uint) (*Product, error) {
	productRepo := repo.NewRepo[Product](s.db)
	return productRepo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context, limit, offset int) ([]Product, error) {
	productRepo := repo.NewRepo[Product](s.db)
	return productRepo.List(ctx, limit, offset)
}

func (s *ProductService) Update(ctx context.Context, id uint, data map[string]interface{}) error {
	productRepo := repo.NewRepo[Product](s.db)
	
	product, err := productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := productRepo.UpdateFields(ctx, id, data); err != nil {
		return err
	}

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

	s.esSync.Update(product)
	return nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	productRepo := repo.NewRepo[Product](s.db)
	
	if err := productRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.esSync.Delete("products", strconv.Itoa(int(id)))
	return nil
}

func (s *ProductService) CalculateCost(ctx context.Context, productID uint, quantity float64, useMinPrice bool) (*CostResult, error) {
	productRepo := repo.NewRepo[Product](s.db)
	
	product, err := productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if len(product.Materials) == 0 {
		return nil, errors.New("产品未配置原料")
	}
	if len(product.Processes) == 0 {
		return nil, errors.New("产品未配置工艺")
	}

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

		mat, err := s.materialService.Get(ctx, m.MaterialID)
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

		proc, err := s.processService.Get(ctx, p.ProcessID)
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

func (s *ProductService) ValidateForSubmit(product *Product) error {
	if product.Name == "" {
		return errors.New("产品名称不能为空")
	}

	if len(product.Materials) == 0 {
		return errors.New("至少需要一个原料")
	}

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

	if len(product.Processes) == 0 {
		return errors.New("至少需要一个工艺")
	}

	return nil
}