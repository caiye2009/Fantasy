package product

import (
	"strconv"
	"back/pkg/fields"
)

type MaterialConfig struct {
	MaterialID uint    `json:"material_id"`
	Ratio      float64 `json:"ratio"`
}

type ProcessConfig struct {
	ProcessID uint    `json:"process_id"`
	Quantity  float64 `json:"quantity"`
}

type Product struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	Name      string            `gorm:"size:100;not null" json:"name"`
	Status    string            `gorm:"size:20;default:'draft'" json:"status"`
	Materials []MaterialConfig  `gorm:"type:jsonb;serializer:json" json:"materials"`
	Processes []ProcessConfig   `gorm:"type:jsonb;serializer:json" json:"processes"`
	fields.DTFields
}

func (Product) TableName() string {
	return "products"
}

// ========== Indexable 接口实现 ==========

func (p *Product) GetIndexName() string {
	return "products"
}

func (p *Product) GetDocumentID() string {
	return strconv.Itoa(int(p.ID))
}

func (p *Product) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":         p.ID,
		"name":       p.Name,
		"status":     p.Status,
		"materials":  p.Materials,
		"processes":  p.Processes,
		"created_at": p.CreatedAt,
		"updated_at": p.UpdatedAt,
	}
}

// ========== 请求结构体 ==========

type CreateProductRequest struct {
	Name      string            `json:"name" binding:"required,min=2,max=100"`
	Materials []MaterialConfig  `json:"materials" binding:"required,dive"`
	Processes []ProcessConfig   `json:"processes"`
}

type UpdateProductRequest struct {
	Name      string            `json:"name" binding:"omitempty,min=2,max=100"`
	Materials []MaterialConfig  `json:"materials" binding:"omitempty,dive"`
	Processes []ProcessConfig   `json:"processes"`
}

// ========== 成本计算相关结构体 ==========

// MaterialCostItem 原料成本明细
type MaterialCostItem struct {
	MaterialID   uint    `json:"material_id"`
	MaterialName string  `json:"material_name"`
	Ratio        float64 `json:"ratio"`
	Price        float64 `json:"price"`
	Cost         float64 `json:"cost"`
}

// ProcessCostItem 工艺成本明细
type ProcessCostItem struct {
	ProcessID   uint    `json:"process_id"`
	ProcessName string  `json:"process_name"`
	Price       float64 `json:"price"`
	Cost        float64 `json:"cost"`
}

// CostBreakdown 成本明细
type CostBreakdown struct {
	Materials []MaterialCostItem `json:"materials"`
	Processes []ProcessCostItem  `json:"processes"`
}

// CostResult 成本计算结果
type CostResult struct {
	MaterialCost float64        `json:"material_cost"`
	ProcessCost  float64        `json:"process_cost"`
	UnitCost     float64        `json:"unit_cost"`
	TotalCost    float64        `json:"total_cost"`
	Breakdown    *CostBreakdown `json:"breakdown"`
}