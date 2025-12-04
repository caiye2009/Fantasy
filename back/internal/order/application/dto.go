package application

import "time"

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	OrderNo   string  `json:"order_no" binding:"required,max=50"`
	ClientID  uint    `json:"client_id" binding:"required"`
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"required,gte=0"`
	CreatedBy uint    `json:"created_by" binding:"required"`
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	Status    string  `json:"status" binding:"omitempty,oneof=pending confirmed production completed cancelled"`
	Quantity  float64 `json:"quantity" binding:"omitempty,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"omitempty,gte=0"`
}

// OrderResponse 订单响应
type OrderResponse struct {
	ID         uint      `json:"id"`
	OrderNo    string    `json:"order_no"`
	ClientID   uint      `json:"client_id"`
	ProductID  uint      `json:"product_id"`
	Quantity   float64   `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedBy  uint      `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Total  int64            `json:"total"`
	Orders []*OrderResponse `json:"orders"`
}