package application

import (
	"context"

	"back/internal/material/domain"
	"back/internal/material/infra"
	"back/pkg/es"
)

type MaterialService struct {
	repo   *infra.MaterialRepo
	syncer *es.EntitySyncer[*domain.Material]
}

func NewMaterialService(repo *infra.MaterialRepo, esSync *es.ESSync) *MaterialService {
	return &MaterialService{
		repo:   repo,
		syncer: es.NewEntitySyncer[*domain.Material](esSync),
	}
}

func (s *MaterialService) Create(ctx context.Context, req *CreateMaterialRequest) (*MaterialResponse, error) {
	material := &domain.Material{
		MaterialCode: req.MaterialCode,
		MaterialName: req.MaterialName,
		MaterialCategory: req.MaterialCategory,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, material); err != nil {
		return nil, err
	}

	// 同步到 ES
	s.syncer.SyncCreate(material)

	return toMaterialResponse(material), nil
}

func (s *MaterialService) Get(ctx context.Context, id uint) (*MaterialResponse, error) {
	material, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrMaterialNotFound
	}
	return toMaterialResponse(material), nil
}

func (s *MaterialService) List(ctx context.Context, limit, offset int) ([]*MaterialResponse, int64, error) {
	materials, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*MaterialResponse, len(materials))
	for i, material := range materials {
		responses[i] = toMaterialResponse(&material)
	}

	return responses, count, nil
}

func (s *MaterialService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrMaterialNotFound
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

	// 同步到 ES - 更新后需要获取完整的实体
	s.syncer.SyncUpdateWithFetch(ctx, id, s.repo.GetByID)

	return nil
}

func (s *MaterialService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrMaterialNotFound
	}

	// 先获取material以获取索引名称
	material, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 从 ES 删除
	s.syncer.SyncDelete(material)

	return nil
}

func toMaterialResponse(material *domain.Material) *MaterialResponse {
	return &MaterialResponse{
		ID: material.ID,
		MaterialCode: material.MaterialCode,
		MaterialName: material.MaterialName,
		MaterialCategory: material.MaterialCategory,
		CreatedBy: material.CreatedBy,
		CreatedAt: material.CreatedAt,
		UpdatedAt: material.UpdatedAt,
	}
}
