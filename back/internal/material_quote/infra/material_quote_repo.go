package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/material_quote/domain"
	"back/pkg/repo"
)

type MaterialQuoteRepo struct {
	*repo.Repo[domain.MaterialQuote]
	db *gorm.DB
}

func NewMaterialQuoteRepo(db *gorm.DB) *MaterialQuoteRepo {
	return &MaterialQuoteRepo{
		Repo: repo.NewRepo[domain.MaterialQuote](db),
		db:   db,
	}
}

func (r *MaterialQuoteRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *MaterialQuoteRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
