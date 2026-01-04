package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/order_material/domain"
	"back/pkg/repo"
)

type OrderMaterialRepo struct {
	*repo.Repo[domain.OrderMaterial]
	db *gorm.DB
}

func NewOrderMaterialRepo(db *gorm.DB) *OrderMaterialRepo {
	return &OrderMaterialRepo{
		Repo: repo.NewRepo[domain.OrderMaterial](db),
		db:   db,
	}
}

func (r *OrderMaterialRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *OrderMaterialRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
