package infra

import (
	"time"
	
	"gorm.io/gorm"
	
	"back/internal/material/domain"
)

// MaterialPO 材料持久化对象
type MaterialPO struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"size:100;not null;index"`
	Spec        string         `gorm:"size:200"`
	Unit        string         `gorm:"size:20"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName 表名
func (MaterialPO) TableName() string {
	return "materials"
}

// ToDomain 转换为领域模型
func (po *MaterialPO) ToDomain() *domain.Material {
	return &domain.Material{
		ID:          po.ID,
		Name:        po.Name,
		Spec:        po.Spec,
		Unit:        po.Unit,
		Description: po.Description,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

// FromDomain 从领域模型转换
func FromDomain(m *domain.Material) *MaterialPO {
	return &MaterialPO{
		ID:          m.ID,
		Name:        m.Name,
		Spec:        m.Spec,
		Unit:        m.Unit,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}