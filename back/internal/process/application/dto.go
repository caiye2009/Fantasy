package application

import "time"

type CreateProcessRequest struct {
	ProcessCode string `json:"processCode" validate:"required,max=50"`
	ProcessName string `json:"processName" validate:"required,max=100"`
	ProcessCountry string `json:"processCountry" validate:"max=50"`
	CreatedBy uint `json:"createdBy" validate:"required"`
}

type ProcessResponse struct {
	ID uint `json:"id"`
	ProcessCode string `json:"processCode"`
	ProcessName string `json:"processName"`
	ProcessCountry string `json:"processCountry"`
	CreatedBy uint `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
