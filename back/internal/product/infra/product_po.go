package infra

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	
	"gorm.io/gorm"
	
	"back/internal/product/domain"
)

// MaterialConfigJSON 原料配置 JSON
type MaterialConfigJSON []domain.MaterialConfig

// Scan 实现 sql.Scanner 接口
func (m *MaterialConfigJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

// Value 实现 driver.Valuer 接口
func (m MaterialConfigJSON) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}

// ProcessConfigJSON 工艺配置 JSON
type ProcessConfigJSON []domain.ProcessConfig

// Scan 实现 sql.Scanner 接口
func (p *ProcessConfigJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// Value 实现 driver.Valuer 接口
func (p ProcessConfigJSON) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}
	return json.Marshal(p)
}

// ProductPO 产品持久化对象
type ProductPO struct {
	ID        uint              `gorm:"primaryKey"`
	Name      string            `gorm:"size:100;not null;index"`
	Status    string            `gorm:"size:20;default:'draft';index"`
	Materials MaterialConfigJSON `gorm:"type:jsonb"`
	Processes ProcessConfigJSON  `gorm:"type:jsonb"`
	CreatedAt time.Time         `gorm:"autoCreateTime"`
	UpdatedAt time.Time         `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt    `gorm:"index"`
}

// TableName 表名
func (ProductPO) TableName() string {
	return "products"
}

// ToDomain 转换为领域模型
func (po *ProductPO) ToDomain() *domain.Product {
	return &domain.Product{
		ID:        po.ID,
		Name:      po.Name,
		Status:    domain.ProductStatus(po.Status),
		Materials: []domain.MaterialConfig(po.Materials),
		Processes: []domain.ProcessConfig(po.Processes),
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

// FromDomain 从领域模型转换
func FromDomain(p *domain.Product) *ProductPO {
	return &ProductPO{
		ID:        p.ID,
		Name:      p.Name,
		Status:    string(p.Status),
		Materials: MaterialConfigJSON(p.Materials),
		Processes: ProcessConfigJSON(p.Processes),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}