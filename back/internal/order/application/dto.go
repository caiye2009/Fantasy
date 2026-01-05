package application

import "time"

type CreateOrderRequest struct {
	OrderCode string `json:"orderCode" validate:"required,max=50"`
	ClientCode string `json:"clientCode" validate:"required,max=50"`
	OrderDate time.Time `json:"orderDate" validate:"required"`
	DeliveryDate time.Time `json:"deliveryDate"`
	OrderStatus string `json:"orderStatus" validate:"required"`
}

type OrderResponse struct {
	ID uint `json:"id"`
	OrderCode string `json:"orderCode"`
	ClientCode string `json:"clientCode"`
	OrderDate time.Time `json:"orderDate"`
	DeliveryDate time.Time `json:"deliveryDate"`
	OrderStatus string `json:"orderStatus"`
	CreatedBy string `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
