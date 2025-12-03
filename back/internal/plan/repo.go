package plan

import (
	"gorm.io/gorm"
	"back/pkg/repo"
)

type PlanRepo struct {
	*repo.Repo[Plan]
}

func NewPlanRepo(db *gorm.DB) *PlanRepo {
	return &PlanRepo{
		Repo: repo.NewRepo[Plan](db),
	}
}