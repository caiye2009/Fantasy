package application

import "time"

type CreateMaterialRequest struct {
	MaterialCode string `json:"materialCode" validate:"required,max=50"`
	MaterialName string `json:"materialName" validate:"required,max=100"`
	MaterialCategory string `json:"materialCategory" validate:"max=50"`
}

type MaterialResponse struct {
	ID uint `json:"id"`
	MaterialCode string `json:"materialCode"`
	MaterialName string `json:"materialName"`
	MaterialCategory string `json:"materialCategory"`
	CreatedBy string `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
