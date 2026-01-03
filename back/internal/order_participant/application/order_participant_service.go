package application

import (
	"context"

	"back/internal/order_participant/domain"
	"back/internal/order_participant/infra"
)

type OrderParticipantService struct {
	repo *infra.OrderParticipantRepo
}

func NewOrderParticipantService(repo *infra.OrderParticipantRepo) *OrderParticipantService {
	return &OrderParticipantService{repo: repo}
}

func (s *OrderParticipantService) Create(ctx context.Context, req *CreateOrderParticipantRequest) (*OrderParticipantResponse, error) {
	orderParticipant := &domain.OrderParticipant{
		OrderCode: req.OrderCode,
		Username: req.Username,
		ParticipantRole: req.ParticipantRole,
	}

	if err := s.repo.Create(ctx, orderParticipant); err != nil {
		return nil, err
	}

	return toOrderParticipantResponse(orderParticipant), nil
}

func (s *OrderParticipantService) Get(ctx context.Context, id uint) (*OrderParticipantResponse, error) {
	orderParticipant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrOrderParticipantNotFound
	}
	return toOrderParticipantResponse(orderParticipant), nil
}

func (s *OrderParticipantService) List(ctx context.Context, limit, offset int) ([]*OrderParticipantResponse, int64, error) {
	orderParticipants, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*OrderParticipantResponse, len(orderParticipants))
	for i, orderParticipant := range orderParticipants {
		responses[i] = toOrderParticipantResponse(&orderParticipant)
	}

	return responses, count, nil
}

func (s *OrderParticipantService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderParticipantNotFound
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

func (s *OrderParticipantService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderParticipantNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toOrderParticipantResponse(orderParticipant *domain.OrderParticipant) *OrderParticipantResponse {
	return &OrderParticipantResponse{
		ID: orderParticipant.ID,
		OrderCode: orderParticipant.OrderCode,
		Username: orderParticipant.Username,
		ParticipantRole: orderParticipant.ParticipantRole,
		CreatedAt: orderParticipant.CreatedAt,
		UpdatedAt: orderParticipant.UpdatedAt,
	}
}
