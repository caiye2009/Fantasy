package infra

import (
	"gorm.io/gorm"

	"back/internal/fabric_process/domain"
	"back/pkg/repo"
)

type FabricProcessRepo struct {
	*repo.Repo[domain.FabricProcess]
	db *gorm.DB
}

func NewFabricProcessRepo(db *gorm.DB) *FabricProcessRepo {
	return &FabricProcessRepo{
		Repo: repo.NewRepo[domain.FabricProcess](db),
		db:   db,
	}
}