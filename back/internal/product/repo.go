package product

import (
	"gorm.io/gorm"
	"back/pkg/repo"
)

type ProductRepo struct {
	*repo.Repo[Product]
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		Repo: repo.NewRepo[Product](db),
	}
}