package application

import "time"

type CreatePurchaseOrderRequest struct {
	PurchaseCode string `json:"purchaseCode" validate:"required,max=50"`
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	SupplierCode string `json:"supplierCode" validate:"required,max=50"`
	PurchaseDate time.Time `json:"purchaseDate" validate:"required"`
	PurchaseStatus string `json:"purchaseStatus" validate:"required"`
	CreatedBy uint `json:"createdBy" validate:"required"`
}

type PurchaseOrderResponse struct {
	ID uint `json:"id"`
	PurchaseCode string `json:"purchaseCode"`
	OrderCode string `json:"orderCode"`
	SupplierCode string `json:"supplierCode"`
	PurchaseDate time.Time `json:"purchaseDate"`
	PurchaseStatus string `json:"purchaseStatus"`
	CreatedBy uint `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
