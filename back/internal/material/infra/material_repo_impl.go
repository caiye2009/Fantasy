package infra

import (
	"context"
	"errors"
	
	"gorm.io/gorm"
	
	"back/internal/material/domain"
)

// MaterialRepoImpl 材料仓储实现
type MaterialRepoImpl struct {
	db *gorm.DB
}

// NewMaterialRepoImpl 创建仓储实现
func NewMaterialRepoImpl(db *gorm.DB) domain.MaterialRepository {
	return &MaterialRepoImpl{db: db}
}

// Save 保存材料
func (r *MaterialRepoImpl) Save(ctx context.Context, material *domain.Material) error {
	po := FromDomain(material)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return err
	}
	material.ID = po.ID
	material.CreatedAt = po.CreatedAt
	material.UpdatedAt = po.UpdatedAt
	return nil
}

// FindByID 根据 ID 查询
func (r *MaterialRepoImpl) FindByID(ctx context.Context, id uint) (*domain.Material, error) {
	var po MaterialPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrMaterialNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindAll 查询所有
func (r *MaterialRepoImpl) FindAll(ctx context.Context, limit, offset int) ([]*domain.Material, error) {
	var pos []MaterialPO
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	materials := make([]*domain.Material, len(pos))
	for i, po := range pos {
		materials[i] = po.ToDomain()
	}
	return materials, nil
}

// Update 更新材料
func (r *MaterialRepoImpl) Update(ctx context.Context, material *domain.Material) error {
	po := FromDomain(material)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除材料
func (r *MaterialRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&MaterialPO{}, id).Error
}

// ExistsByID 检查是否存在
func (r *MaterialRepoImpl) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&MaterialPO{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByName 根据名称查询
func (r *MaterialRepoImpl) FindByName(ctx context.Context, name string) (*domain.Material, error) {
	var po MaterialPO
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrMaterialNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// Count 统计数量
func (r *MaterialRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&MaterialPO{}).Count(&count).Error
	return count, err
}