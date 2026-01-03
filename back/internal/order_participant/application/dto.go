package application

import "time"

type CreateOrderParticipantRequest struct {
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	Username string `json:"username" validate:"required,max=50"`
	ParticipantRole string `json:"participantRole" validate:"required,max=50"`
}

type OrderParticipantResponse struct {
	ID uint `json:"id"`
	OrderCode string `json:"orderCode"`
	Username string `json:"username"`
	ParticipantRole string `json:"participantRole"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
