package domain

import "context"

// SupplierPriceRepository 供应商价格仓储接口
type SupplierPriceRepository interface {
	// 基础操作
	Save(ctx context.Context, price *SupplierPrice) error
	
	// 查询操作
	FindMinPrice(ctx context.Context, targetType TargetType, targetID uint) (*SupplierPrice, error)
	FindMaxPrice(ctx context.Context, targetType TargetType, targetID uint) (*SupplierPrice, error)
	FindHistory(ctx context.Context, targetType TargetType, targetID uint, limit int) ([]*SupplierPrice, error)
	FindBySupplier(ctx context.Context, supplierID uint, limit int) ([]*SupplierPrice, error)
}