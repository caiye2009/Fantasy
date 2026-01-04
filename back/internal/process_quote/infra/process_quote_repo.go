package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/process_quote/domain"
	"back/pkg/repo"
)

type ProcessQuoteRepo struct {
	*repo.Repo[domain.ProcessQuote]
	db *gorm.DB
}

func NewProcessQuoteRepo(db *gorm.DB) *ProcessQuoteRepo {
	return &ProcessQuoteRepo{
		Repo: repo.NewRepo[domain.ProcessQuote](db),
		db:   db,
	}
}

func (r *ProcessQuoteRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *ProcessQuoteRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
