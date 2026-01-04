package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/product_allocation/domain"
	"back/pkg/repo"
)

type ProductAllocationRepo struct {
	*repo.Repo[domain.ProductAllocation]
	db *gorm.DB
}

func NewProductAllocationRepo(db *gorm.DB) *ProductAllocationRepo {
	return &ProductAllocationRepo{
		Repo: repo.NewRepo[domain.ProductAllocation](db),
		db:   db,
	}
}

func (r *ProductAllocationRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *ProductAllocationRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
