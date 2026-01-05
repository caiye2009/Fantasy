package application

import "time"

type CreateOrderMaterialRequest struct {
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	MaterialCode string `json:"materialCode" validate:"required,max=50"`
	RequiredQty float64 `json:"requiredQty" validate:"required,gt=0"`
}

type OrderMaterialResponse struct {
	ID uint `json:"id"`
	OrderCode string `json:"orderCode"`
	MaterialCode string `json:"materialCode"`
	RequiredQty float64 `json:"requiredQty"`
	CreatedBy string `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
