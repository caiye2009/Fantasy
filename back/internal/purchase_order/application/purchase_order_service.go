package application

import (
	"context"

	"back/internal/purchase_order/domain"
	"back/internal/purchase_order/infra"
)

type PurchaseOrderService struct {
	repo *infra.PurchaseOrderRepo
}

func NewPurchaseOrderService(repo *infra.PurchaseOrderRepo) *PurchaseOrderService {
	return &PurchaseOrderService{repo: repo}
}

func (s *PurchaseOrderService) Create(ctx context.Context, req *CreatePurchaseOrderRequest) (*PurchaseOrderResponse, error) {
	purchaseOrder := &domain.PurchaseOrder{
		PurchaseCode: req.PurchaseCode,
		OrderCode: req.OrderCode,
		SupplierCode: req.SupplierCode,
		PurchaseDate: req.PurchaseDate,
		PurchaseStatus: req.PurchaseStatus,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, purchaseOrder); err != nil {
		return nil, err
	}

	return toPurchaseOrderResponse(purchaseOrder), nil
}

func (s *PurchaseOrderService) Get(ctx context.Context, id uint) (*PurchaseOrderResponse, error) {
	purchaseOrder, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrPurchaseOrderNotFound
	}
	return toPurchaseOrderResponse(purchaseOrder), nil
}

func (s *PurchaseOrderService) List(ctx context.Context, limit, offset int) ([]*PurchaseOrderResponse, int64, error) {
	purchaseOrders, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*PurchaseOrderResponse, len(purchaseOrders))
	for i, purchaseOrder := range purchaseOrders {
		responses[i] = toPurchaseOrderResponse(&purchaseOrder)
	}

	return responses, count, nil
}

func (s *PurchaseOrderService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrPurchaseOrderNotFound
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

func (s *PurchaseOrderService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrPurchaseOrderNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toPurchaseOrderResponse(purchaseOrder *domain.PurchaseOrder) *PurchaseOrderResponse {
	return &PurchaseOrderResponse{
		ID: purchaseOrder.ID,
		PurchaseCode: purchaseOrder.PurchaseCode,
		OrderCode: purchaseOrder.OrderCode,
		SupplierCode: purchaseOrder.SupplierCode,
		PurchaseDate: purchaseOrder.PurchaseDate,
		PurchaseStatus: purchaseOrder.PurchaseStatus,
		CreatedBy: purchaseOrder.CreatedBy,
		CreatedAt: purchaseOrder.CreatedAt,
		UpdatedAt: purchaseOrder.UpdatedAt,
	}
}
