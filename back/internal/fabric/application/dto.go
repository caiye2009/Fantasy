package application

import "time"

type CreateFabricRequest struct {
	FabricCode string `json:"fabricCode" validate:"required,max=50"`
	FabricName string `json:"fabricName" validate:"required,max=100"`
}

type FabricResponse struct {
	ID         uint      `json:"id"`
	FabricCode string    `json:"fabricCode"`
	FabricName string    `json:"fabricName"`
	CreatedBy  string    `json:"createdBy"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}