package infra

import (
	"context"

	"gorm.io/gorm"

	"back/internal/fabric/domain"
	"back/pkg/repo"
)

type FabricRepo struct {
	*repo.Repo[domain.Fabric]
	db *gorm.DB
}

func NewFabricRepo(db *gorm.DB) *FabricRepo {
	return &FabricRepo{
		Repo: repo.NewRepo[domain.Fabric](db),
		db:   db,
	}
}

func (r *FabricRepo) Exists(ctx context.Context, id uint) (bool, error) {
	return r.Repo.Exists(ctx, map[string]interface{}{"id": id})
}

func (r *FabricRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}