package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/order_participant/domain"
	"back/pkg/repo"
)

// OrderParticipantRepo ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·ÃÂÃÂ¤ÃÂÃÂ»ÃÂÃÂÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¨
type OrderParticipantRepo struct {
	*repo.Repo[domain.OrderParticipant]
	db *gorm.DB
}

// NewOrderParticipantRepo ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ»ÃÂÃÂºÃÂÃÂ¤ÃÂÃÂ»ÃÂÃÂÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¨
func NewOrderParticipantRepo(db *gorm.DB) *OrderParticipantRepo {
	return &OrderParticipantRepo{
		Repo: repo.NewRepo[domain.OrderParticipant](db),
		db:   db,
	}
}

// Exists ÃÂÃÂ¦ÃÂÃÂ£ÃÂÃÂÃÂÃÂ¦ÃÂÃÂÃÂÃÂ¥ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ¯ÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¦ÃÂÃÂ¥ÃÂÃÂ­ÃÂÃÂÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¨
func (r *OrderParticipantRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

// Count ÃÂÃÂ§ÃÂÃÂ»ÃÂÃÂÃÂÃÂ¨ÃÂÃÂ®ÃÂÃÂ¡ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ°ÃÂÃÂ©ÃÂÃÂÃÂÃÂ
func (r *OrderParticipantRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
