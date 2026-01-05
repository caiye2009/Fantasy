package application

import (
	"context"

	"back/internal/fabric_material/domain"
	"back/internal/fabric_material/infra"
	"back/pkg/auth"
)

type FabricMaterialService struct {
	repo   *infra.FabricMaterialRepo
}

func NewFabricMaterialService(repo *infra.FabricMaterialRepo) *FabricMaterialService {
	return &FabricMaterialService{
		repo:   repo,
	}
}

func (s *FabricMaterialService) Create(ctx context.Context, req *CreateFabricMaterialRequest) (*FabricMaterialResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	fm := &domain.FabricMaterial{
		FabricCode:   req.FabricCode,
		MaterialCode: req.MaterialCode,
		Ratio:        req.Ratio,
		CreatedBy:    createdBy,
	}

	if err := s.repo.Create(ctx, fm); err != nil {
		return nil, err
	}

	return toFabricMaterialResponse(fm), nil
}

func (s *FabricMaterialService) Get(ctx context.Context, id uint) (*FabricMaterialResponse, error) {
	fm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrFabricMaterialNotFound
	}
	return toFabricMaterialResponse(fm), nil
}

func (s *FabricMaterialService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func toFabricMaterialResponse(fm *domain.FabricMaterial) *FabricMaterialResponse {
	return &FabricMaterialResponse{
		ID:           fm.ID,
		FabricCode:   fm.FabricCode,
		MaterialCode: fm.MaterialCode,
		Ratio:        fm.Ratio,
		CreatedBy:    fm.CreatedBy,
		CreatedAt:    fm.CreatedAt,
		UpdatedAt:    fm.UpdatedAt,
	}
}