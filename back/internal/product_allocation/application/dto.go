package application

import "time"

type CreateProductAllocationRequest struct {
	AllocationCode string `json:"allocationCode" validate:"required,max=50"`
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	ProductCode string `json:"productCode" validate:"required,max=50"`
	AllocationQty float64 `json:"allocationQty" validate:"required,gt=0"`
	ActualQty float64 `json:"actualQty" validate:"gte=0"`
	AllocationStatus string `json:"allocationStatus" validate:"required"`
	AllocatedBy uint `json:"allocatedBy" validate:"required"`
	AllocatedAt time.Time `json:"allocatedAt"`
}

type ProductAllocationResponse struct {
	ID uint `json:"id"`
	AllocationCode string `json:"allocationCode"`
	OrderCode string `json:"orderCode"`
	ProductCode string `json:"productCode"`
	AllocationQty float64 `json:"allocationQty"`
	ActualQty float64 `json:"actualQty"`
	AllocationStatus string `json:"allocationStatus"`
	AllocatedBy uint `json:"allocatedBy"`
	AllocatedAt time.Time `json:"allocatedAt"`
	CreatedBy string `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
