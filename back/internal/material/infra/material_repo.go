package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/material/domain"
	"back/pkg/repo"
)

// MaterialRepo 客户仓储
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

// Exists 检查是否存在
func (r *MaterialRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

// Count 统计数量
func (r *MaterialRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
