package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"back/internal/material/domain"
	"back/pkg/repo"
)

// MaterialRepo 材料仓储
type MaterialRepo struct {
	*repo.Repo[domain.Material]
	db *gorm.DB
}

// NewMaterialRepo 创建仓储
func NewMaterialRepo(db *gorm.DB) *MaterialRepo {
	return &MaterialRepo{
		Repo: repo.NewRepo[domain.Material](db),
		db:   db,
	}
}

// Save 保存材料
func (r *MaterialRepo) Save(ctx context.Context, material *domain.Material) error {
	return r.Create(ctx, material)
}

// FindByID 根据 ID 查询
func (r *MaterialRepo) FindByID(ctx context.Context, id uint) (*domain.Material, error) {
	material, err := r.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrMaterialNotFound
	}
	if err != nil {
		return nil, err
	}
	return material, nil
}

// FindAll 查询所有
func (r *MaterialRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.Material, error) {
	materials, err := r.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Material, len(materials))
	for i := range materials {
		result[i] = &materials[i]
	}
	return result, nil
}

// Update 更新材料
func (r *MaterialRepo) Update(ctx context.Context, material *domain.Material) error {
	return r.Repo.Update(ctx, material)
}

// Delete 删除材料
func (r *MaterialRepo) Delete(ctx context.Context, id uint) error {
	return r.Repo.Delete(ctx, id)
}

// ExistsByID 检查是否存在
func (r *MaterialRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"id": id})
}

// FindByName 根据名称查询
func (r *MaterialRepo) FindByName(ctx context.Context, name string) (*domain.Material, error) {
	material, err := r.First(ctx, map[string]interface{}{"name": name})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrMaterialNotFound
	}
	if err != nil {
		return nil, err
	}
	return material, nil
}

// Count 统计数量
func (r *MaterialRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
