package application

import "time"

type CreateMaterialShrinkageRequest struct {
	MaterialCode  string  `json:"materialCode" validate:"required,max=50"`
	SupplierCode  string  `json:"supplierCode" validate:"required,max=50"`
	ShrinkageRate float64 `json:"shrinkageRate" validate:"required,gt=0"`
}

type MaterialShrinkageResponse struct {
	ShrinkageID   string    `json:"shrinkageId"`
	MaterialCode  string    `json:"materialCode"`
	SupplierCode  string    `json:"supplierCode"`
	ShrinkageRate float64   `json:"shrinkageRate"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
