package application

import "time"

type CreateProcessQuoteRequest struct {
	QuoteID string `json:"quoteId" validate:"required,max=50"`
	ProcessCode string `json:"processCode" validate:"required,max=50"`
	SupplierCode string `json:"supplierCode" validate:"required,max=50"`
	QuotePrice float64 `json:"quotePrice" validate:"required,gt=0"`
	CreatedBy uint `json:"createdBy" validate:"required"`
}

type ProcessQuoteResponse struct {
	ID uint `json:"id"`
	QuoteID string `json:"quoteId"`
	ProcessCode string `json:"processCode"`
	SupplierCode string `json:"supplierCode"`
	QuotePrice float64 `json:"quotePrice"`
	CreatedBy uint `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
