package client

import (
	"gorm.io/gorm"
	"back/pkg/repo"
)

type ClientRepo struct {
	*repo.Repo[Client]
}

func NewClientRepo(db *gorm.DB) *ClientRepo {
	return &ClientRepo{
		Repo: repo.NewRepo[Client](db),
	}
}