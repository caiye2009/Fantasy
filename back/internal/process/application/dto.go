package application

import "time"

type CreateProcessRequest struct {
	ProcessCode string `json:"processCode" validate:"required,max=50"`
	ProcessName string `json:"processName" validate:"required,max=100"`
	ProcessCategory string `json:"processCategory" validate:"max=50"`
}

type ProcessResponse struct {
	ID uint `json:"id"`
	ProcessCode string `json:"processCode"`
	ProcessName string `json:"processName"`
	ProcessCategory string `json:"processCategory"`
	CreatedBy string `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
