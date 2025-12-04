package domain

import (
	"time"
)

// TargetType 目标类型
type TargetType string

const (
	TargetTypeMaterial TargetType = "material"
	TargetTypeProcess  TargetType = "process"
)

// SupplierPrice 供应商价格聚合根
type SupplierPrice struct {
	ID         uint
	TargetType TargetType
	TargetID   uint
	SupplierID uint
	Price      float64
	QuotedAt   time.Time
	CreatedAt  time.Time
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