package infra

import (
	"time"
	
	"gorm.io/gorm"
	
	"back/internal/process/domain"
)

// ProcessPO 工序持久化对象
type ProcessPO struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"size:100;not null;index"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName 表名
func (ProcessPO) TableName() string {
	return "processes"
}

// ToDomain 转换为领域模型
func (po *ProcessPO) ToDomain() *domain.Process {
	return &domain.Process{
		ID:          po.ID,
		Name:        po.Name,
		Description: po.Description,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

// FromDomain 从领域模型转换
func FromDomain(p *domain.Process) *ProcessPO {
	return &ProcessPO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}