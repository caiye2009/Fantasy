package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/pricing/domain"
	"back/pkg/repo"
)

// SupplierPriceRepo 供应商价格仓储实现
type SupplierPriceRepo struct {
	*repo.Repo[domain.SupplierPrice]
	db *gorm.DB
}

// NewSupplierPriceRepo 创建仓储实现
func NewSupplierPriceRepo(db *gorm.DB) *SupplierPriceRepo {
	return &SupplierPriceRepo{
		Repo: repo.NewRepo[domain.SupplierPrice](db),
		db:   db,
	}
}

// Save 保存价格
func (r *SupplierPriceRepo) Save(ctx context.Context, price *domain.SupplierPrice) error {
	return r.Create(ctx, price)
}

// FindMinPrice 查找最低价
func (r *SupplierPriceRepo) FindMinPrice(ctx context.Context, targetType string, targetID uint) (*domain.SupplierPrice, error) {
	var result domain.SupplierPrice
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("price ASC").
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPriceNotFound
		}
		return nil, err
	}

	return &result, nil
}

// FindMaxPrice 查找最高价
func (r *SupplierPriceRepo) FindMaxPrice(ctx context.Context, targetType string, targetID uint) (*domain.SupplierPrice, error) {
	var result domain.SupplierPrice
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("price DESC").
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPriceNotFound
		}
		return nil, err
	}

	return &result, nil
}

// FindHistory 查找报价历史
func (r *SupplierPriceRepo) FindHistory(ctx context.Context, targetType string, targetID uint, limit int) ([]*domain.SupplierPrice, error) {
	var results []*domain.SupplierPrice
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("quoted_at DESC").
		Limit(limit).
		Find(&results).Error

	return results, err
}

// FindBySupplier 根据供应商查找
func (r *SupplierPriceRepo) FindBySupplier(ctx context.Context, supplierID uint, limit int) ([]*domain.SupplierPrice, error) {
	var results []*domain.SupplierPrice
	err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("quoted_at DESC").
		Limit(limit).
		Find(&results).Error

	return results, err
}
