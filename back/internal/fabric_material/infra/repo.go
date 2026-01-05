package infra

import (
	"gorm.io/gorm"

	"back/internal/fabric_material/domain"
	"back/pkg/repo"
)

type FabricMaterialRepo struct {
	*repo.Repo[domain.FabricMaterial]
	db *gorm.DB
}

func NewFabricMaterialRepo(db *gorm.DB) *FabricMaterialRepo {
	return &FabricMaterialRepo{
		Repo: repo.NewRepo[domain.FabricMaterial](db),
		db:   db,
	}
}