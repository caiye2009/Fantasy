package application

import "time"

type CreateProductMaterialRequest struct {
	ProductCode string `json:"productCode" validate:"required,max=50"`
	MaterialCode string `json:"materialCode" validate:"required,max=50"`
	RequiredQty float64 `json:"requiredQty" validate:"required,gt=0"`
	CreatedBy uint `json:"createdBy" validate:"required"`
}

type ProductMaterialResponse struct {
	ID uint `json:"id"`
	ProductCode string `json:"productCode"`
	MaterialCode string `json:"materialCode"`
	RequiredQty float64 `json:"requiredQty"`
	CreatedBy uint `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
