package vendor

import (
	"gorm.io/gorm"
	"back/pkg/repo"
)

type VendorRepo struct {
	*repo.Repo[Vendor]
}

func NewVendorRepo(db *gorm.DB) *VendorRepo {
	return &VendorRepo{
		Repo: repo.NewRepo[Vendor](db),
	}
}