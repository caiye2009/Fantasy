package infra

import (
	"context"
	"errors"
	
	"gorm.io/gorm"
	
	"back/internal/product/domain"
)

// ProductRepoImpl 产品仓储实现
type ProductRepoImpl struct {
	db *gorm.DB
}

// NewProductRepoImpl 创建仓储实现
func NewProductRepoImpl(db *gorm.DB) domain.ProductRepository {
	return &ProductRepoImpl{db: db}
}

// Save 保存产品
func (r *ProductRepoImpl) Save(ctx context.Context, product *domain.Product) error {
	po := FromDomain(product)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return err
	}
	product.ID = po.ID
	product.CreatedAt = po.CreatedAt
	product.UpdatedAt = po.UpdatedAt
	return nil
}

// FindByID 根据 ID 查询
func (r *ProductRepoImpl) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	var po ProductPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindAll 查询所有
func (r *ProductRepoImpl) FindAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	var pos []ProductPO
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	products := make([]*domain.Product, len(pos))
	for i, po := range pos {
		products[i] = po.ToDomain()
	}
	return products, nil
}

// Update 更新产品
func (r *ProductRepoImpl) Update(ctx context.Context, product *domain.Product) error {
	po := FromDomain(product)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除产品
func (r *ProductRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&ProductPO{}, id).Error
}

// ExistsByID 检查是否存在
func (r *ProductRepoImpl) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ProductPO{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByName 根据名称查询
func (r *ProductRepoImpl) FindByName(ctx context.Context, name string) (*domain.Product, error) {
	var po ProductPO
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindByStatus 根据状态查询
func (r *ProductRepoImpl) FindByStatus(ctx context.Context, status domain.ProductStatus, limit, offset int) ([]*domain.Product, error) {
	var pos []ProductPO
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	products := make([]*domain.Product, len(pos))
	for i, po := range pos {
		products[i] = po.ToDomain()
	}
	return products, nil
}

// Count 统计数量
func (r *ProductRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ProductPO{}).Count(&count).Error
	return count, err
}