package application

import (
	"context"

	"back/internal/order_process/domain"
	"back/internal/order_process/infra"
)

type OrderProcessService struct {
	repo *infra.OrderProcessRepo
}

func NewOrderProcessService(repo *infra.OrderProcessRepo) *OrderProcessService {
	return &OrderProcessService{repo: repo}
}

func (s *OrderProcessService) Create(ctx context.Context, req *CreateOrderProcessRequest) (*OrderProcessResponse, error) {
	orderProcess := &domain.OrderProcess{
		OrderCode: req.OrderCode,
		ProcessCode: req.ProcessCode,
		ProcessSeq: req.ProcessSeq,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, orderProcess); err != nil {
		return nil, err
	}

	return toOrderProcessResponse(orderProcess), nil
}

func (s *OrderProcessService) Get(ctx context.Context, id uint) (*OrderProcessResponse, error) {
	orderProcess, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrOrderProcessNotFound
	}
	return toOrderProcessResponse(orderProcess), nil
}

func (s *OrderProcessService) List(ctx context.Context, limit, offset int) ([]*OrderProcessResponse, int64, error) {
	orderProcesss, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*OrderProcessResponse, len(orderProcesss))
	for i, orderProcess := range orderProcesss {
		responses[i] = toOrderProcessResponse(&orderProcess)
	}

	return responses, count, nil
}

func (s *OrderProcessService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderProcessNotFound
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

func (s *OrderProcessService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderProcessNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toOrderProcessResponse(orderProcess *domain.OrderProcess) *OrderProcessResponse {
	return &OrderProcessResponse{
		ID: orderProcess.ID,
		OrderCode: orderProcess.OrderCode,
		ProcessCode: orderProcess.ProcessCode,
		ProcessSeq: orderProcess.ProcessSeq,
		CreatedBy: orderProcess.CreatedBy,
		CreatedAt: orderProcess.CreatedAt,
		UpdatedAt: orderProcess.UpdatedAt,
	}
}
