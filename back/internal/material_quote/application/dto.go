package application

import "time"

type CreateMaterialQuoteRequest struct {
	QuoteID string `json:"quoteId" validate:"required,max=50"`
	MaterialCode string `json:"materialCode" validate:"required,max=50"`
	SupplierCode string `json:"supplierCode" validate:"required,max=50"`
	QuotePrice float64 `json:"quotePrice" validate:"required,gt=0"`
	CreatedBy uint `json:"createdBy" validate:"required"`
}

type MaterialQuoteResponse struct {
	ID uint `json:"id"`
	QuoteID string `json:"quoteId"`
	MaterialCode string `json:"materialCode"`
	SupplierCode string `json:"supplierCode"`
	QuotePrice float64 `json:"quotePrice"`
	CreatedBy uint `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
