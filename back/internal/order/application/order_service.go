package application

import (
	"context"

	"back/internal/order/domain"
	"back/internal/order/infra"
)

type OrderService struct {
	repo *infra.OrderRepo
}

func NewOrderService(repo *infra.OrderRepo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error) {
	order := &domain.Order{
		OrderCode: req.OrderCode,
		ClientCode: req.ClientCode,
		OrderDate: req.OrderDate,
		DeliveryDate: req.DeliveryDate,
		OrderStatus: req.OrderStatus,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return toOrderResponse(order), nil
}

func (s *OrderService) Get(ctx context.Context, id uint) (*OrderResponse, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}
	return toOrderResponse(order), nil
}

func (s *OrderService) List(ctx context.Context, limit, offset int) ([]*OrderResponse, int64, error) {
	orders, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = toOrderResponse(&order)
	}

	return responses, count, nil
}

func (s *OrderService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderNotFound
	}

	delete(updates, "id")
	delete(updates, "createdAt")
	delete(updates, "createdBy")
	delete(updates, "deletedAt")

	if len(updates) == 0 {
		return nil
	}

	return s.repo.UpdateFields(ctx, id, updates)
}

func (s *OrderService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toOrderResponse(order *domain.Order) *OrderResponse {
	return &OrderResponse{
		ID: order.ID,
		OrderCode: order.OrderCode,
		ClientCode: order.ClientCode,
		OrderDate: order.OrderDate,
		DeliveryDate: order.DeliveryDate,
		OrderStatus: order.OrderStatus,
		CreatedBy: order.CreatedBy,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}
