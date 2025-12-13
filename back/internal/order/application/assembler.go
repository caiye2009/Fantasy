package application

import "back/internal/order/domain"

// ToOrder DTO → Domain Model
func ToOrder(req *CreateOrderRequest) *domain.Order {
	order := &domain.Order{
		OrderNo:   req.OrderNo,
		ClientID:  req.ClientID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
		Status:    domain.OrderStatusPending,
		CreatedBy: req.CreatedBy,
	}
	order.CalculateTotalPrice()
	return order
}

// ToOrderResponse Domain Model → DTO
func ToOrderResponse(o *domain.Order) *OrderResponse {
	return &OrderResponse{
		ID:         o.ID,
		OrderNo:    o.OrderNo,
		ClientID:   o.ClientID,
		ProductID:  o.ProductID,
		Quantity:   o.Quantity,
		UnitPrice:  o.UnitPrice,
		TotalPrice: o.TotalPrice,
		Status:     o.Status,
		CreatedBy:  o.CreatedBy,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
}

// ToOrderListResponse Domain Models → List DTO
func ToOrderListResponse(orders []*domain.Order, total int64) *OrderListResponse {
	responses := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = ToOrderResponse(o)
	}
	
	return &OrderListResponse{
		Total:  total,
		Orders: responses,
	}
}