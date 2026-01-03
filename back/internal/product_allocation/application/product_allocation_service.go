package application

import (
	"context"

	"back/internal/product_allocation/domain"
	"back/internal/product_allocation/infra"
)

type ProductAllocationService struct {
	repo *infra.ProductAllocationRepo
}

func NewProductAllocationService(repo *infra.ProductAllocationRepo) *ProductAllocationService {
	return &ProductAllocationService{repo: repo}
}

func (s *ProductAllocationService) Create(ctx context.Context, req *CreateProductAllocationRequest) (*ProductAllocationResponse, error) {
	productAllocation := &domain.ProductAllocation{
		AllocationCode: req.AllocationCode,
		OrderCode: req.OrderCode,
		ProductCode: req.ProductCode,
		AllocationQty: req.AllocationQty,
		ActualQty: req.ActualQty,
		AllocationStatus: req.AllocationStatus,
		AllocatedBy: req.AllocatedBy,
		AllocatedAt: req.AllocatedAt,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, productAllocation); err != nil {
		return nil, err
	}

	return toProductAllocationResponse(productAllocation), nil
}

func (s *ProductAllocationService) Get(ctx context.Context, id uint) (*ProductAllocationResponse, error) {
	productAllocation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductAllocationNotFound
	}
	return toProductAllocationResponse(productAllocation), nil
}

func (s *ProductAllocationService) List(ctx context.Context, limit, offset int) ([]*ProductAllocationResponse, int64, error) {
	productAllocations, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*ProductAllocationResponse, len(productAllocations))
	for i, productAllocation := range productAllocations {
		responses[i] = toProductAllocationResponse(&productAllocation)
	}

	return responses, count, nil
}

func (s *ProductAllocationService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductAllocationNotFound
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

func (s *ProductAllocationService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductAllocationNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toProductAllocationResponse(productAllocation *domain.ProductAllocation) *ProductAllocationResponse {
	return &ProductAllocationResponse{
		ID: productAllocation.ID,
		AllocationCode: productAllocation.AllocationCode,
		OrderCode: productAllocation.OrderCode,
		ProductCode: productAllocation.ProductCode,
		AllocationQty: productAllocation.AllocationQty,
		ActualQty: productAllocation.ActualQty,
		AllocationStatus: productAllocation.AllocationStatus,
		AllocatedBy: productAllocation.AllocatedBy,
		AllocatedAt: productAllocation.AllocatedAt,
		CreatedBy: productAllocation.CreatedBy,
		CreatedAt: productAllocation.CreatedAt,
		UpdatedAt: productAllocation.UpdatedAt,
	}
}
