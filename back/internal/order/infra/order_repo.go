package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/order/domain"
	"back/pkg/repo"
)

// OrderRepo ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·ÃÂÃÂ¤ÃÂÃÂ»ÃÂÃÂÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¨
type OrderRepo struct {
	*repo.Repo[domain.Order]
	db *gorm.DB
}

// NewOrderRepo ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ»ÃÂÃÂºÃÂÃÂ¤ÃÂÃÂ»ÃÂÃÂÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¨
func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{
		Repo: repo.NewRepo[domain.Order](db),
		db:   db,
	}
}

// Exists ÃÂÃÂ¦ÃÂÃÂ£ÃÂÃÂÃÂÃÂ¦ÃÂÃÂÃÂÃÂ¥ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ¯ÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¦ÃÂÃÂ¥ÃÂÃÂ­ÃÂÃÂÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¨
func (r *OrderRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

// Count ÃÂÃÂ§ÃÂÃÂ»ÃÂÃÂÃÂÃÂ¨ÃÂÃÂ®ÃÂÃÂ¡ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ°ÃÂÃÂ©ÃÂÃÂÃÂÃÂ
func (r *OrderRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
