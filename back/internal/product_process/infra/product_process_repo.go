package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/product_process/domain"
	"back/pkg/repo"
)

type ProductProcessRepo struct {
	*repo.Repo[domain.ProductProcess]
	db *gorm.DB
}

func NewProductProcessRepo(db *gorm.DB) *ProductProcessRepo {
	return &ProductProcessRepo{
		Repo: repo.NewRepo[domain.ProductProcess](db),
		db:   db,
	}
}

func (r *ProductProcessRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *ProductProcessRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
