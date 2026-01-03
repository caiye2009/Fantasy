package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/process/domain"
	"back/pkg/repo"
)

// ProcessRepo 客户仓储
type ProcessRepo struct {
	*repo.Repo[domain.Process]
	db *gorm.DB
}

// NewProcessRepo 创建仓储
func NewProcessRepo(db *gorm.DB) *ProcessRepo {
	return &ProcessRepo{
		Repo: repo.NewRepo[domain.Process](db),
		db:   db,
	}
}

// Exists 检查是否存在
func (r *ProcessRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

// Count 统计数量
func (r *ProcessRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
