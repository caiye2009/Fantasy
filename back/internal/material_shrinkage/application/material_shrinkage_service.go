package application

import (
	"context"

	"github.com/google/uuid"

	"back/internal/material_shrinkage/domain"
	"back/internal/material_shrinkage/infra"
	"back/pkg/auth"
)

type MaterialShrinkageService struct {
	repo *infra.MaterialShrinkageRepo
}

func NewMaterialShrinkageService(repo *infra.MaterialShrinkageRepo) *MaterialShrinkageService {
	return &MaterialShrinkageService{repo: repo}
}

func (s *MaterialShrinkageService) Create(ctx context.Context, req *CreateMaterialShrinkageRequest) (*MaterialShrinkageResponse, error) {
	// 从 context 获取 createdBy
	createdBy := auth.GetLoginIDFromContext(ctx)

	materialShrinkage := &domain.MaterialShrinkage{
		ShrinkageID:   uuid.NewString(),
		MaterialCode:  req.MaterialCode,
		SupplierCode:  req.SupplierCode,
		ShrinkageRate: req.ShrinkageRate,
		CreatedBy:     createdBy,
	}

	if err := s.repo.Create(ctx, materialShrinkage); err != nil {
		return nil, err
	}

	return toMaterialShrinkageResponse(materialShrinkage), nil
}

func toMaterialShrinkageResponse(m *domain.MaterialShrinkage) *MaterialShrinkageResponse {
	return &MaterialShrinkageResponse{
		ShrinkageID:   m.ShrinkageID,
		MaterialCode:  m.MaterialCode,
		SupplierCode:  m.SupplierCode,
		ShrinkageRate: m.ShrinkageRate,
		CreatedBy:     m.CreatedBy,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}