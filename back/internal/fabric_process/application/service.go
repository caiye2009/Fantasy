package application

import (
	"context"

	"back/internal/fabric_process/domain"
	"back/internal/fabric_process/infra"
	"back/pkg/auth"
)

type FabricProcessService struct {
	repo   *infra.FabricProcessRepo
}

func NewFabricProcessService(repo *infra.FabricProcessRepo) *FabricProcessService {
	return &FabricProcessService{
		repo:   repo,
	}
}

func (s *FabricProcessService) Create(ctx context.Context, req *CreateFabricProcessRequest) (*FabricProcessResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	fp := &domain.FabricProcess{
		FabricCode:  req.FabricCode,
		ProcessCode: req.ProcessCode,
		StepOrder:   req.StepOrder,
		CreatedBy:   createdBy,
	}

	if err := s.repo.Create(ctx, fp); err != nil {
		return nil, err
	}

	return toFabricProcessResponse(fp), nil
}

func (s *FabricProcessService) Get(ctx context.Context, id uint) (*FabricProcessResponse, error) {
	fp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrFabricProcessNotFound
	}
	return toFabricProcessResponse(fp), nil
}

func (s *FabricProcessService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func toFabricProcessResponse(fp *domain.FabricProcess) *FabricProcessResponse {
	return &FabricProcessResponse{
		ID:          fp.ID,
		FabricCode:  fp.FabricCode,
		ProcessCode: fp.ProcessCode,
		StepOrder:   fp.StepOrder,
		CreatedBy:   fp.CreatedBy,
		CreatedAt:   fp.CreatedAt,
		UpdatedAt:   fp.UpdatedAt,
	}
}