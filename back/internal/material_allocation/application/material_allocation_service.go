package application

import (
	"context"

	"back/internal/material_allocation/domain"
	"back/internal/material_allocation/infra"
)

type MaterialAllocationService struct {
	repo *infra.MaterialAllocationRepo
}

func NewMaterialAllocationService(repo *infra.MaterialAllocationRepo) *MaterialAllocationService {
	return &MaterialAllocationService{repo: repo}
}

func (s *MaterialAllocationService) Create(ctx context.Context, req *CreateMaterialAllocationRequest) (*MaterialAllocationResponse, error) {
	materialAllocation := &domain.MaterialAllocation{
		AllocationCode: req.AllocationCode,
		OrderCode: req.OrderCode,
		MaterialCode: req.MaterialCode,
		AllocationQty: req.AllocationQty,
		ActualQty: req.ActualQty,
		AllocationStatus: req.AllocationStatus,
		AllocatedBy: req.AllocatedBy,
		AllocatedAt: req.AllocatedAt,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, materialAllocation); err != nil {
		return nil, err
	}

	return toMaterialAllocationResponse(materialAllocation), nil
}

func (s *MaterialAllocationService) Get(ctx context.Context, id uint) (*MaterialAllocationResponse, error) {
	materialAllocation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrMaterialAllocationNotFound
	}
	return toMaterialAllocationResponse(materialAllocation), nil
}

func (s *MaterialAllocationService) List(ctx context.Context, limit, offset int) ([]*MaterialAllocationResponse, int64, error) {
	materialAllocations, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*MaterialAllocationResponse, len(materialAllocations))
	for i, materialAllocation := range materialAllocations {
		responses[i] = toMaterialAllocationResponse(&materialAllocation)
	}

	return responses, count, nil
}

func (s *MaterialAllocationService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrMaterialAllocationNotFound
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

func (s *MaterialAllocationService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrMaterialAllocationNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toMaterialAllocationResponse(materialAllocation *domain.MaterialAllocation) *MaterialAllocationResponse {
	return &MaterialAllocationResponse{
		ID: materialAllocation.ID,
		AllocationCode: materialAllocation.AllocationCode,
		OrderCode: materialAllocation.OrderCode,
		MaterialCode: materialAllocation.MaterialCode,
		AllocationQty: materialAllocation.AllocationQty,
		ActualQty: materialAllocation.ActualQty,
		AllocationStatus: materialAllocation.AllocationStatus,
		AllocatedBy: materialAllocation.AllocatedBy,
		AllocatedAt: materialAllocation.AllocatedAt,
		CreatedBy: materialAllocation.CreatedBy,
		CreatedAt: materialAllocation.CreatedAt,
		UpdatedAt: materialAllocation.UpdatedAt,
	}
}
