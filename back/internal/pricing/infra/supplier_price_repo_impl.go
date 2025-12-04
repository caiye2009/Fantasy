package infra

import (
	"context"
	
	"gorm.io/gorm"
	
	"back/internal/pricing/domain"
)

// SupplierPriceRepoImpl 供应商价格仓储实现
type SupplierPriceRepoImpl struct {
	db *gorm.DB
}

// NewSupplierPriceRepoImpl 创建仓储实现
func NewSupplierPriceRepoImpl(db *gorm.DB) domain.SupplierPriceRepository {
	return &SupplierPriceRepoImpl{db: db}
}

// Save 保存价格
func (r *SupplierPriceRepoImpl) Save(ctx context.Context, price *domain.SupplierPrice) error {
	po := FromDomain(price)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return err
	}
	price.ID = po.ID
	price.CreatedAt = po.CreatedAt
	return nil
}

// FindMinPrice 查找最低价
func (r *SupplierPriceRepoImpl) FindMinPrice(ctx context.Context, targetType domain.TargetType, targetID uint) (*domain.SupplierPrice, error) {
	var po SupplierPricePO
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("price ASC").
		First(&po).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPriceNotFound
		}
		return nil, err
	}
	
	return po.ToDomain(), nil
}

// FindMaxPrice 查找最高价
func (r *SupplierPriceRepoImpl) FindMaxPrice(ctx context.Context, targetType domain.TargetType, targetID uint) (*domain.SupplierPrice, error) {
	var po SupplierPricePO
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("price DESC").
		First(&po).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPriceNotFound
		}
		return nil, err
	}
	
	return po.ToDomain(), nil
}

// FindHistory 查找报价历史
func (r *SupplierPriceRepoImpl) FindHistory(ctx context.Context, targetType domain.TargetType, targetID uint, limit int) ([]*domain.SupplierPrice, error) {
	var pos []SupplierPricePO
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("quoted_at DESC").
		Limit(limit).
		Find(&pos).Error
	
	if err != nil {
		return nil, err
	}
	
	prices := make([]*domain.SupplierPrice, len(pos))
	for i, po := range pos {
		prices[i] = po.ToDomain()
	}
	
	return prices, nil
}

// FindBySupplier 根据供应商查找
func (r *SupplierPriceRepoImpl) FindBySupplier(ctx context.Context, supplierID uint, limit int) ([]*domain.SupplierPrice, error) {
	var pos []SupplierPricePO
	err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("quoted_at DESC").
		Limit(limit).
		Find(&pos).Error
	
	if err != nil {
		return nil, err
	}
	
	prices := make([]*domain.SupplierPrice, len(pos))
	for i, po := range pos {
		prices[i] = po.ToDomain()
	}
	
	return prices, nil
}