package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/material_quote/domain"
	"back/pkg/repo"
)

// MaterialQuoteRepo ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·ÃÂÃÂ¤ÃÂÃÂ»ÃÂÃÂÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¨
type MaterialQuoteRepo struct {
	*repo.Repo[domain.MaterialQuote]
	db *gorm.DB
}

// NewMaterialQuoteRepo ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ»ÃÂÃÂºÃÂÃÂ¤ÃÂÃÂ»ÃÂÃÂÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¨
func NewMaterialQuoteRepo(db *gorm.DB) *MaterialQuoteRepo {
	return &MaterialQuoteRepo{
		Repo: repo.NewRepo[domain.MaterialQuote](db),
		db:   db,
	}
}

// Exists ÃÂÃÂ¦ÃÂÃÂ£ÃÂÃÂÃÂÃÂ¦ÃÂÃÂÃÂÃÂ¥ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ¯ÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¦ÃÂÃÂ¥ÃÂÃÂ­ÃÂÃÂÃÂÃÂ¥ÃÂÃÂÃÂÃÂ¨
func (r *MaterialQuoteRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

// Count ÃÂÃÂ§ÃÂÃÂ»ÃÂÃÂÃÂÃÂ¨ÃÂÃÂ®ÃÂÃÂ¡ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ°ÃÂÃÂ©ÃÂÃÂÃÂÃÂ
func (r *MaterialQuoteRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
