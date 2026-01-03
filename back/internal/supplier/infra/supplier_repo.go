package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/supplier/domain"
	"back/pkg/repo"
)

// SupplierRepo 供应商仓储
type SupplierRepo struct {
	*repo.Repo[domain.Supplier]
	db *gorm.DB
}

// NewSupplierRepo 创建仓储
func NewSupplierRepo(db *gorm.DB) *SupplierRepo {
	return &SupplierRepo{
		Repo: repo.NewRepo[domain.Supplier](db),
		db:   db,
	}
}

// Exists 检查是否存在
func (r *SupplierRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

// Count 统计数量
func (r *SupplierRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
