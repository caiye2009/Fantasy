package application

import "time"

type CreateProductionOrderRequest struct {
	ProductionCode string `json:"productionCode" validate:"required,max=50"`
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	ProcessCode string `json:"processCode" validate:"required,max=50"`
	SupplierCode string `json:"supplierCode" validate:"required,max=50"`
	ProductionQty float64 `json:"productionQty" validate:"required,gt=0"`
	ProductionStatus string `json:"productionStatus" validate:"required"`
	CreatedBy uint `json:"createdBy" validate:"required"`
}

type ProductionOrderResponse struct {
	ID uint `json:"id"`
	ProductionCode string `json:"productionCode"`
	OrderCode string `json:"orderCode"`
	ProcessCode string `json:"processCode"`
	SupplierCode string `json:"supplierCode"`
	ProductionQty float64 `json:"productionQty"`
	ProductionStatus string `json:"productionStatus"`
	CreatedBy uint `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
