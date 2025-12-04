package infra

import (
	"time"
	
	"back/internal/pricing/domain"
)

// SupplierPricePO 供应商价格持久化对象
type SupplierPricePO struct {
	ID         uint      `gorm:"primaryKey"`
	TargetType string    `gorm:"size:20;not null;index:idx_target_price"`
	TargetID   uint      `gorm:"not null;index:idx_target_price"`
	SupplierID uint      `gorm:"not null;index"`
	Price      float64   `gorm:"type:decimal(10,2);not null;index:idx_target_price"`
	QuotedAt   time.Time `gorm:"not null;index:idx_target_time"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// TableName 表名
func (SupplierPricePO) TableName() string {
	return "supplier_prices"
}

// ToDomain 转换为领域模型
func (po *SupplierPricePO) ToDomain() *domain.SupplierPrice {
	return &domain.SupplierPrice{
		ID:         po.ID,
		TargetType: domain.TargetType(po.TargetType),
		TargetID:   po.TargetID,
		SupplierID: po.SupplierID,
		Price:      po.Price,
		QuotedAt:   po.QuotedAt,
		CreatedAt:  po.CreatedAt,
	}
}

// FromDomain 从领域模型转换
func FromDomain(sp *domain.SupplierPrice) *SupplierPricePO {
	return &SupplierPricePO{
		ID:         sp.ID,
		TargetType: string(sp.TargetType),
		TargetID:   sp.TargetID,
		SupplierID: sp.SupplierID,
		Price:      sp.Price,
		QuotedAt:   sp.QuotedAt,
		CreatedAt:  sp.CreatedAt,
	}
}