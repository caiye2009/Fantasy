package domain

import (
	"time"
)

// 目标类型常量
const (
	TargetTypeMaterial = "material"
	TargetTypeProcess  = "process"
)

// SupplierPrice 供应商价格聚合根
type SupplierPrice struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TargetType string    `gorm:"size:20;not null;index:idx_target_price" json:"target_type"`
	TargetID   uint      `gorm:"not null;index:idx_target_price" json:"target_id"`
	SupplierID uint      `gorm:"not null;index" json:"supplier_id"`
	Price      float64   `gorm:"type:decimal(10,2);not null;index:idx_target_price" json:"price"`
	QuotedAt   time.Time `gorm:"not null;index:idx_target_time" json:"quoted_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 表名
func (SupplierPrice) TableName() string {
	return "supplier_prices"
}

// Validate 验证价格数据
func (sp *SupplierPrice) Validate() error {
	if sp.TargetType == "" {
		return ErrTargetTypeRequired
	}

	if sp.TargetType != TargetTypeMaterial && sp.TargetType != TargetTypeProcess {
		return ErrInvalidTargetType
	}

	if sp.TargetID == 0 {
		return ErrTargetIDRequired
	}

	if sp.SupplierID == 0 {
		return ErrSupplierIDRequired
	}

	if sp.Price <= 0 {
		return ErrInvalidPrice
	}

	return nil
}

// IsLowerThan 是否低于指定价格
func (sp *SupplierPrice) IsLowerThan(price float64) bool {
	return sp.Price < price
}

// IsHigherThan 是否高于指定价格
func (sp *SupplierPrice) IsHigherThan(price float64) bool {
	return sp.Price > price
}