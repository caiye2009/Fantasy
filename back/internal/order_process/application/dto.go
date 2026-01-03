package application

import "time"

type CreateOrderProcessRequest struct {
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	ProcessCode string `json:"processCode" validate:"required,max=50"`
	ProcessSeq int `json:"processSeq" validate:"required,gte=0"`
	CreatedBy uint `json:"createdBy" validate:"required"`
}

type OrderProcessResponse struct {
	ID uint `json:"id"`
	OrderCode string `json:"orderCode"`
	ProcessCode string `json:"processCode"`
	ProcessSeq int `json:"processSeq"`
	CreatedBy uint `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
