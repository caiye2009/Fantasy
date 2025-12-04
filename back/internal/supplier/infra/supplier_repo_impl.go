package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"back/internal/supplier/domain"
)

// SupplierRepoImpl 供应商仓储实现
type SupplierRepoImpl struct {
	db *gorm.DB
}

// NewSupplierRepoImpl 创建仓储实现
func NewSupplierRepoImpl(db *gorm.DB) domain.SupplierRepository {
	return &SupplierRepoImpl{db: db}
}

// Save 保存供应商
func (r *SupplierRepoImpl) Save(ctx context.Context, supplier *domain.Supplier) error {
	po := FromDomain(supplier)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return err
	}
	supplier.ID = po.ID
	supplier.CreatedAt = po.CreatedAt
	supplier.UpdatedAt = po.UpdatedAt
	return nil
}

// FindByID 根据 ID 查询
func (r *SupplierRepoImpl) FindByID(ctx context.Context, id uint) (*domain.Supplier, error) {
	var po SupplierPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSupplierNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindAll 查询所有
func (r *SupplierRepoImpl) FindAll(ctx context.Context, limit, offset int) ([]*domain.Supplier, error) {
	var pos []SupplierPO
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}

	suppliers := make([]*domain.Supplier, len(pos))
	for i, po := range pos {
		suppliers[i] = po.ToDomain()
	}
	return suppliers, nil
}

// Update 更新供应商
func (r *SupplierRepoImpl) Update(ctx context.Context, supplier *domain.Supplier) error {
	po := FromDomain(supplier)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除供应商
func (r *SupplierRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&SupplierPO{}, id).Error
}

// ExistsByID 检查是否存在
func (r *SupplierRepoImpl) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SupplierPO{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByName 根据名称查询
func (r *SupplierRepoImpl) FindByName(ctx context.Context, name string) (*domain.Supplier, error) {
	var po SupplierPO
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSupplierNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ExistsByName 检查名称是否存在
func (r *SupplierRepoImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SupplierPO{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// FindByPhone 根据电话查询
func (r *SupplierRepoImpl) FindByPhone(ctx context.Context, phone string) (*domain.Supplier, error) {
	var po SupplierPO
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSupplierNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindByEmail 根据邮箱查询
func (r *SupplierRepoImpl) FindByEmail(ctx context.Context, email string) (*domain.Supplier, error) {
	var po SupplierPO
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSupplierNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// Count 统计数量
func (r *SupplierRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&SupplierPO{}).Count(&count).Error
	return count, err
}
