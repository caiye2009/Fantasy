package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/product/domain"
	"back/pkg/repo"
)

// ProductRepo 客户仓储
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

// Exists 检查是否存在
func (r *ProductRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

// Count 统计数量
func (r *ProductRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
