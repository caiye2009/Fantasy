package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/order_process/domain"
	"back/pkg/repo"
)

type OrderProcessRepo struct {
	*repo.Repo[domain.OrderProcess]
	db *gorm.DB
}

func NewOrderProcessRepo(db *gorm.DB) *OrderProcessRepo {
	return &OrderProcessRepo{
		Repo: repo.NewRepo[domain.OrderProcess](db),
		db:   db,
	}
}

func (r *OrderProcessRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *OrderProcessRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
