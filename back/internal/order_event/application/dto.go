package application

import "time"

type CreateOrderEventRequest struct {
	EventID string `json:"eventId" validate:"required,max=50"`
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	EventType string `json:"eventType" validate:"required,max=50"`
	EventContent string `json:"eventContent"`
	Operator string `json:"operator" validate:"required,max=100"`
	OperatedAt time.Time `json:"operatedAt" validate:"required"`
}

type OrderEventResponse struct {
	ID uint `json:"id"`
	EventID string `json:"eventId"`
	OrderCode string `json:"orderCode"`
	EventType string `json:"eventType"`
	EventContent string `json:"eventContent"`
	Operator string `json:"operator"`
	OperatedAt time.Time `json:"operatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
