package application

import "time"

type CreateFabricProcessRequest struct {
	FabricCode  string `json:"fabricCode" validate:"required,max=50"`
	ProcessCode string `json:"processCode" validate:"required,max=50"`
	StepOrder   int    `json:"stepOrder" validate:"required,gt=0"`
}

type FabricProcessResponse struct {
	ID          uint      `json:"id"`
	FabricCode  string    `json:"fabricCode"`
	ProcessCode string    `json:"processCode"`
	StepOrder   int       `json:"stepOrder"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}