package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/product_fabric/domain"
	"back/pkg/repo"
)

type ProductFabricRepo struct {
	*repo.Repo[domain.ProductFabric]
	db *gorm.DB
}

func NewProductFabricRepo(db *gorm.DB) *ProductFabricRepo {
	return &ProductFabricRepo{
		Repo: repo.NewRepo[domain.ProductFabric](db),
		db:   db,
	}
}

func (r *ProductFabricRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *ProductFabricRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}