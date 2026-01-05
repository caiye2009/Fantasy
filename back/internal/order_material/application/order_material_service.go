package application

import (
	"context"

	"back/internal/order_material/domain"
	"back/internal/order_material/infra"
	"back/pkg/auth"
)

type OrderMaterialService struct {
	repo *infra.OrderMaterialRepo
}

func NewOrderMaterialService(repo *infra.OrderMaterialRepo) *OrderMaterialService {
	return &OrderMaterialService{repo: repo}
}

func (s *OrderMaterialService) Create(ctx context.Context, req *CreateOrderMaterialRequest) (*OrderMaterialResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	orderMaterial := &domain.OrderMaterial{
		OrderCode: req.OrderCode,
		MaterialCode: req.MaterialCode,
		RequiredQty: req.RequiredQty,
		CreatedBy: createdBy,
	}

	if err := s.repo.Create(ctx, orderMaterial); err != nil {
		return nil, err
	}

	return toOrderMaterialResponse(orderMaterial), nil
}

func (s *OrderMaterialService) Get(ctx context.Context, id uint) (*OrderMaterialResponse, error) {
	orderMaterial, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrOrderMaterialNotFound
	}
	return toOrderMaterialResponse(orderMaterial), nil
}

func (s *OrderMaterialService) List(ctx context.Context, limit, offset int) ([]*OrderMaterialResponse, int64, error) {
	orderMaterials, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*OrderMaterialResponse, len(orderMaterials))
	for i, orderMaterial := range orderMaterials {
		responses[i] = toOrderMaterialResponse(&orderMaterial)
	}

	return responses, count, nil
}

func (s *OrderMaterialService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderMaterialNotFound
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

func (s *OrderMaterialService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderMaterialNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toOrderMaterialResponse(orderMaterial *domain.OrderMaterial) *OrderMaterialResponse {
	return &OrderMaterialResponse{
		ID: orderMaterial.ID,
		OrderCode: orderMaterial.OrderCode,
		MaterialCode: orderMaterial.MaterialCode,
		RequiredQty: orderMaterial.RequiredQty,
		CreatedBy: orderMaterial.CreatedBy,
		CreatedAt: orderMaterial.CreatedAt,
		UpdatedAt: orderMaterial.UpdatedAt,
	}
}
