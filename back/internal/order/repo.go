package order

import (
	"gorm.io/gorm"
	"back/pkg/repo"
)

type OrderRepo struct {
	*repo.Repo[Order]
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{
		Repo: repo.NewRepo[Order](db),
	}
}