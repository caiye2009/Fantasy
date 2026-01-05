package application

import (
	"context"

	"back/internal/production_order/domain"
	"back/internal/production_order/infra"
	"back/pkg/auth"
)

type ProductionOrderService struct {
	repo *infra.ProductionOrderRepo
}

func NewProductionOrderService(repo *infra.ProductionOrderRepo) *ProductionOrderService {
	return &ProductionOrderService{repo: repo}
}

func (s *ProductionOrderService) Create(ctx context.Context, req *CreateProductionOrderRequest) (*ProductionOrderResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	productionOrder := &domain.ProductionOrder{
		ProductionCode: req.ProductionCode,
		OrderCode: req.OrderCode,
		ProcessCode: req.ProcessCode,
		SupplierCode: req.SupplierCode,
		ProductionQty: req.ProductionQty,
		ProductionStatus: req.ProductionStatus,
		CreatedBy: createdBy,
	}

	if err := s.repo.Create(ctx, productionOrder); err != nil {
		return nil, err
	}

	return toProductionOrderResponse(productionOrder), nil
}

func (s *ProductionOrderService) Get(ctx context.Context, id uint) (*ProductionOrderResponse, error) {
	productionOrder, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductionOrderNotFound
	}
	return toProductionOrderResponse(productionOrder), nil
}

func (s *ProductionOrderService) List(ctx context.Context, limit, offset int) ([]*ProductionOrderResponse, int64, error) {
	productionOrders, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*ProductionOrderResponse, len(productionOrders))
	for i, productionOrder := range productionOrders {
		responses[i] = toProductionOrderResponse(&productionOrder)
	}

	return responses, count, nil
}

func (s *ProductionOrderService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductionOrderNotFound
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

func (s *ProductionOrderService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductionOrderNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toProductionOrderResponse(productionOrder *domain.ProductionOrder) *ProductionOrderResponse {
	return &ProductionOrderResponse{
		ID: productionOrder.ID,
		ProductionCode: productionOrder.ProductionCode,
		OrderCode: productionOrder.OrderCode,
		ProcessCode: productionOrder.ProcessCode,
		SupplierCode: productionOrder.SupplierCode,
		ProductionQty: productionOrder.ProductionQty,
		ProductionStatus: productionOrder.ProductionStatus,
		CreatedBy: productionOrder.CreatedBy,
		CreatedAt: productionOrder.CreatedAt,
		UpdatedAt: productionOrder.UpdatedAt,
	}
}
