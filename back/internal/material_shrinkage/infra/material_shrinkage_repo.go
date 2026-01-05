package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/material_shrinkage/domain"
	"back/pkg/repo"
)

type MaterialShrinkageRepo struct {
	*repo.Repo[domain.MaterialShrinkage]
	db *gorm.DB
}

func NewMaterialShrinkageRepo(db *gorm.DB) *MaterialShrinkageRepo {
	return &MaterialShrinkageRepo{
		Repo: repo.NewRepo[domain.MaterialShrinkage](db),
		db:   db,
	}
}

func (r *MaterialShrinkageRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *MaterialShrinkageRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}