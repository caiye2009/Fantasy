package application

import (
	"context"

	"back/internal/order_event/domain"
	"back/internal/order_event/infra"
)

type OrderEventService struct {
	repo *infra.OrderEventRepo
}

func NewOrderEventService(repo *infra.OrderEventRepo) *OrderEventService {
	return &OrderEventService{repo: repo}
}

func (s *OrderEventService) Create(ctx context.Context, req *CreateOrderEventRequest) (*OrderEventResponse, error) {
	orderEvent := &domain.OrderEvent{
		EventID: req.EventID,
		OrderCode: req.OrderCode,
		EventType: req.EventType,
		EventContent: req.EventContent,
		Operator: req.Operator,
		OperatedAt: req.OperatedAt,
	}

	if err := s.repo.Create(ctx, orderEvent); err != nil {
		return nil, err
	}

	return toOrderEventResponse(orderEvent), nil
}

func (s *OrderEventService) Get(ctx context.Context, id uint) (*OrderEventResponse, error) {
	orderEvent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrOrderEventNotFound
	}
	return toOrderEventResponse(orderEvent), nil
}

func (s *OrderEventService) List(ctx context.Context, limit, offset int) ([]*OrderEventResponse, int64, error) {
	orderEvents, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*OrderEventResponse, len(orderEvents))
	for i, orderEvent := range orderEvents {
		responses[i] = toOrderEventResponse(&orderEvent)
	}

	return responses, count, nil
}

func (s *OrderEventService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderEventNotFound
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

func (s *OrderEventService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderEventNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toOrderEventResponse(orderEvent *domain.OrderEvent) *OrderEventResponse {
	return &OrderEventResponse{
		ID: orderEvent.ID,
		EventID: orderEvent.EventID,
		OrderCode: orderEvent.OrderCode,
		EventType: orderEvent.EventType,
		EventContent: orderEvent.EventContent,
		Operator: orderEvent.Operator,
		OperatedAt: orderEvent.OperatedAt,
		CreatedAt: orderEvent.CreatedAt,
		UpdatedAt: orderEvent.UpdatedAt,
	}
}
