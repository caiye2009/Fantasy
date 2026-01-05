package application

import "time"

type CreateOrderProcessRequest struct {
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	ProcessCode string `json:"processCode" validate:"required,max=50"`
	ProcessSeq int `json:"processSeq" validate:"required,gte=0"`
}

type OrderProcessResponse struct {
	ID uint `json:"id"`
	OrderCode string `json:"orderCode"`
	ProcessCode string `json:"processCode"`
	ProcessSeq int `json:"processSeq"`
	CreatedBy string `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
