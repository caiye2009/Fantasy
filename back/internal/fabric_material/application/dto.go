package application

import "time"

type CreateFabricMaterialRequest struct {
	FabricCode   string  `json:"fabricCode" validate:"required,max=50"`
	MaterialCode string  `json:"materialCode" validate:"required,max=50"`
	Ratio        float64 `json:"ratio" validate:"required,gt=0"`
}

type FabricMaterialResponse struct {
	ID          uint      `json:"id"`
	FabricCode  string    `json:"fabricCode"`
	MaterialCode string   `json:"materialCode"`
	Ratio       float64   `json:"ratio"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}