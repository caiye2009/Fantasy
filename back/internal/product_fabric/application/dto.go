package application

import "time"

type CreateProductFabricRequest struct {
	ProductCode string  `json:"productCode" validate:"required,max=50"`
	FabricCode  string  `json:"fabricCode" validate:"required,max=50"`
	RequiredQty float64 `json:"requiredQty" validate:"required,gt=0"`
}

type ProductFabricResponse struct {
	ID          uint      `json:"id"`
	ProductCode string    `json:"productCode"`
	FabricCode  string    `json:"fabricCode"`
	RequiredQty float64   `json:"requiredQty"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}