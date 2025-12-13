package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"back/internal/product/domain"
	"back/pkg/repo"
)

// ProductRepo 产品仓储实现
type ProductRepo struct {
	*repo.Repo[domain.Product]
	db *gorm.DB
}

// NewProductRepo 创建仓储
func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		Repo: repo.NewRepo[domain.Product](db),
		db:   db,
	}
}

// Save 保存产品
func (r *ProductRepo) Save(ctx context.Context, product *domain.Product) error {
	return r.Create(ctx, product)
}

// FindByID 根据 ID 查询
func (r *ProductRepo) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	product, err := r.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return product, nil
}

// FindAll 查询所有
func (r *ProductRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	products, err := r.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Product, len(products))
	for i := range products {
		result[i] = &products[i]
	}
	return result, nil
}

// Update 更新产品
func (r *ProductRepo) Update(ctx context.Context, product *domain.Product) error {
	return r.Repo.Update(ctx, product)
}

// Delete 删除产品
func (r *ProductRepo) Delete(ctx context.Context, id uint) error {
	return r.Repo.Delete(ctx, id)
}

// ExistsByID 检查是否存在
func (r *ProductRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"id": id})
}

// FindByName 根据名称查询
func (r *ProductRepo) FindByName(ctx context.Context, name string) (*domain.Product, error) {
	product, err := r.First(ctx, map[string]interface{}{"name": name})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return product, nil
}

// FindByStatus 根据状态查询
func (r *ProductRepo) FindByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.Product, error) {
	var result []*domain.Product
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&result).Error
	return result, err
}

// Count 统计数量
func (r *ProductRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
