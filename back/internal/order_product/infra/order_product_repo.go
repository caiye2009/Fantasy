package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/order_product/domain"
	"back/pkg/repo"
)

type OrderProductRepo struct {
	*repo.Repo[domain.OrderProduct]
	db *gorm.DB
}

func NewOrderProductRepo(db *gorm.DB) *OrderProductRepo {
	return &OrderProductRepo{
		Repo: repo.NewRepo[domain.OrderProduct](db),
		db:   db,
	}
}

func (r *OrderProductRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *OrderProductRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
