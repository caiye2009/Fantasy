package application

import "time"

type CreateProductRequest struct {
	ProductCode string `json:"productCode" validate:"required,max=50"`
	ProductName string `json:"productName" validate:"required,max=100"`
	CreatedBy uint `json:"createdBy" validate:"required"`
}

type ProductResponse struct {
	ID uint `json:"id"`
	ProductCode string `json:"productCode"`
	ProductName string `json:"productName"`
	CreatedBy uint `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
