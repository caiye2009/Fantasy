package infra

import (
	"time"
	
	"gorm.io/gorm"
	
	"back/internal/plan/domain"
)

// PlanPO 计划持久化对象
type PlanPO struct {
	ID          uint           `gorm:"primaryKey"`
	PlanNo      string         `gorm:"size:50;uniqueIndex;not null"`
	OrderID     uint           `gorm:"not null;index"`
	ProductID   uint           `gorm:"not null;index"`
	Quantity    float64        `gorm:"type:decimal(10,2);not null"`
	Status      string         `gorm:"size:20;default:'planned';index"`
	ScheduledAt *time.Time     `gorm:"type:timestamp"`
	CompletedAt *time.Time     `gorm:"type:timestamp"`
	CreatedBy   uint           `gorm:"not null;index"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName 表名
func (PlanPO) TableName() string {
	return "plans"
}

// ToDomain 转换为领域模型
func (po *PlanPO) ToDomain() *domain.Plan {
	return &domain.Plan{
		ID:          po.ID,
		PlanNo:      po.PlanNo,
		OrderID:     po.OrderID,
		ProductID:   po.ProductID,
		Quantity:    po.Quantity,
		Status:      domain.PlanStatus(po.Status),
		ScheduledAt: po.ScheduledAt,
		CompletedAt: po.CompletedAt,
		CreatedBy:   po.CreatedBy,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

// FromDomain 从领域模型转换
func FromDomain(p *domain.Plan) *PlanPO {
	return &PlanPO{
		ID:          p.ID,
		PlanNo:      p.PlanNo,
		OrderID:     p.OrderID,
		ProductID:   p.ProductID,
		Quantity:    p.Quantity,
		Status:      string(p.Status),
		ScheduledAt: p.ScheduledAt,
		CompletedAt: p.CompletedAt,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}