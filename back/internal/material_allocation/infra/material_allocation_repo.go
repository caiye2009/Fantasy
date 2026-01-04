package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/material_allocation/domain"
	"back/pkg/repo"
)

type MaterialAllocationRepo struct {
	*repo.Repo[domain.MaterialAllocation]
	db *gorm.DB
}

func NewMaterialAllocationRepo(db *gorm.DB) *MaterialAllocationRepo {
	return &MaterialAllocationRepo{
		Repo: repo.NewRepo[domain.MaterialAllocation](db),
		db:   db,
	}
}

func (r *MaterialAllocationRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *MaterialAllocationRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
