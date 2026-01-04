package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/order_event/domain"
	"back/pkg/repo"
)

type OrderEventRepo struct {
	*repo.Repo[domain.OrderEvent]
	db *gorm.DB
}

func NewOrderEventRepo(db *gorm.DB) *OrderEventRepo {
	return &OrderEventRepo{
		Repo: repo.NewRepo[domain.OrderEvent](db),
		db:   db,
	}
}

func (r *OrderEventRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *OrderEventRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
