package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/production_order/domain"
	"back/pkg/repo"
)

type ProductionOrderRepo struct {
	*repo.Repo[domain.ProductionOrder]
	db *gorm.DB
}

func NewProductionOrderRepo(db *gorm.DB) *ProductionOrderRepo {
	return &ProductionOrderRepo{
		Repo: repo.NewRepo[domain.ProductionOrder](db),
		db:   db,
	}
}

func (r *ProductionOrderRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *ProductionOrderRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
