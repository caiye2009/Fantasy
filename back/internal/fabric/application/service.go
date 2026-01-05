package application

import (
	"context"

	"back/internal/fabric/domain"
	"back/internal/fabric/infra"
	"back/pkg/auth"
	"back/pkg/es"
)

type FabricService struct {
	repo   *infra.FabricRepo
	syncer *es.EntitySyncer[*domain.Fabric]
}

func NewFabricService(repo *infra.FabricRepo, esSync *es.ESSync) *FabricService {
	return &FabricService{
		repo:   repo,
		syncer: es.NewEntitySyncer[*domain.Fabric](esSync),
	}
}

func (s *FabricService) Create(ctx context.Context, req *CreateFabricRequest) (*FabricResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	fabric := &domain.Fabric{
		FabricCode: req.FabricCode,
		FabricName: req.FabricName,
		CreatedBy:  createdBy,
	}

	if err := s.repo.Create(ctx, fabric); err != nil {
		return nil, err
	}

	s.syncer.SyncCreate(fabric)

	return toFabricResponse(fabric), nil
}

func (s *FabricService) Get(ctx context.Context, id uint) (*FabricResponse, error) {
	fabric, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrFabricNotFound
	}
	return toFabricResponse(fabric), nil
}

func (s *FabricService) List(ctx context.Context, limit, offset int) ([]*FabricResponse, int64, error) {
	fabrics, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*FabricResponse, len(fabrics))
	for i, f := range fabrics {
		responses[i] = toFabricResponse(&f)
	}

	return responses, count, nil
}

func (s *FabricService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrFabricNotFound
	}

	delete(updates, "id")
	delete(updates, "createdAt")
	delete(updates, "createdBy")
	delete(updates, "deletedAt")

	if len(updates) == 0 {
		return nil
	}

	if err := s.repo.UpdateFields(ctx, id, updates); err != nil {
		return err
	}

	s.syncer.SyncUpdateWithFetch(ctx, id, s.repo.GetByID)

	return nil
}

func (s *FabricService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrFabricNotFound
	}

	fabric, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.syncer.SyncDelete(fabric)

	return nil
}

func toFabricResponse(fabric *domain.Fabric) *FabricResponse {
	return &FabricResponse{
		ID:         fabric.ID,
		FabricCode: fabric.FabricCode,
		FabricName: fabric.FabricName,
		CreatedBy:  fabric.CreatedBy,
		CreatedAt:  fabric.CreatedAt,
		UpdatedAt:  fabric.UpdatedAt,
	}
}