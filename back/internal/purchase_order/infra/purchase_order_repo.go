package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/purchase_order/domain"
	"back/pkg/repo"
)

type PurchaseOrderRepo struct {
	*repo.Repo[domain.PurchaseOrder]
	db *gorm.DB
}

func NewPurchaseOrderRepo(db *gorm.DB) *PurchaseOrderRepo {
	return &PurchaseOrderRepo{
		Repo: repo.NewRepo[domain.PurchaseOrder](db),
		db:   db,
	}
}

func (r *PurchaseOrderRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *PurchaseOrderRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
