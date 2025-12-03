package process

import (
	"gorm.io/gorm"
	"back/pkg/repo"
)

type ProcessRepo struct {
	*repo.Repo[Process]
}

func NewProcessRepo(db *gorm.DB) *ProcessRepo {
	return &ProcessRepo{
		Repo: repo.NewRepo[Process](db),
	}
}