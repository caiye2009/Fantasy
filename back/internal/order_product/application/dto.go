package application

import "time"

type CreateOrderProductRequest struct {
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	ProductCode string `json:"productCode" validate:"required,max=50"`
	OrderedQty float64 `json:"orderedQty" validate:"required,gt=0"`
	CreatedBy uint `json:"createdBy" validate:"required"`
}

type OrderProductResponse struct {
	ID uint `json:"id"`
	OrderCode string `json:"orderCode"`
	ProductCode string `json:"productCode"`
	OrderedQty float64 `json:"orderedQty"`
	CreatedBy uint `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
