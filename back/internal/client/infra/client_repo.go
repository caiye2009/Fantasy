package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/client/domain"
	"back/pkg/repo"
)

// ClientRepo 客户仓储
type ClientRepo struct {
	*repo.Repo[domain.Client]
	db *gorm.DB
}

// NewClientRepo 创建仓储
func NewClientRepo(db *gorm.DB) *ClientRepo {
	return &ClientRepo{
		Repo: repo.NewRepo[domain.Client](db),
		db:   db,
	}
}

// Exists 检查是否存在
func (r *ClientRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

// Count 统计数量
func (r *ClientRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
