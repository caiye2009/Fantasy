package application

import (
	"context"

	"back/internal/order_product/domain"
	"back/internal/order_product/infra"
	"back/pkg/auth"
)

type OrderProductService struct {
	repo *infra.OrderProductRepo
}

func NewOrderProductService(repo *infra.OrderProductRepo) *OrderProductService {
	return &OrderProductService{repo: repo}
}

func (s *OrderProductService) Create(ctx context.Context, req *CreateOrderProductRequest) (*OrderProductResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	orderProduct := &domain.OrderProduct{
		OrderCode: req.OrderCode,
		ProductCode: req.ProductCode,
		OrderedQty: req.OrderedQty,
		CreatedBy: createdBy,
	}

	if err := s.repo.Create(ctx, orderProduct); err != nil {
		return nil, err
	}

	return toOrderProductResponse(orderProduct), nil
}

func (s *OrderProductService) Get(ctx context.Context, id uint) (*OrderProductResponse, error) {
	orderProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrOrderProductNotFound
	}
	return toOrderProductResponse(orderProduct), nil
}

func (s *OrderProductService) List(ctx context.Context, limit, offset int) ([]*OrderProductResponse, int64, error) {
	orderProducts, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*OrderProductResponse, len(orderProducts))
	for i, orderProduct := range orderProducts {
		responses[i] = toOrderProductResponse(&orderProduct)
	}

	return responses, count, nil
}

func (s *OrderProductService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderProductNotFound
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

func (s *OrderProductService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrOrderProductNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toOrderProductResponse(orderProduct *domain.OrderProduct) *OrderProductResponse {
	return &OrderProductResponse{
		ID: orderProduct.ID,
		OrderCode: orderProduct.OrderCode,
		ProductCode: orderProduct.ProductCode,
		OrderedQty: orderProduct.OrderedQty,
		CreatedBy: orderProduct.CreatedBy,
		CreatedAt: orderProduct.CreatedAt,
		UpdatedAt: orderProduct.UpdatedAt,
	}
}
