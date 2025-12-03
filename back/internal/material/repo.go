package material

import (
	"gorm.io/gorm"
	"back/pkg/repo"
)

type MaterialRepo struct {
	*repo.Repo[Material]
}

func NewMaterialRepo(db *gorm.DB) *MaterialRepo {
	return &MaterialRepo{
		Repo: repo.NewRepo[Material](db),
	}
}