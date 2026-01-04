package application

import (
	"context"

	"back/internal/material/domain"
	"back/internal/material/infra"
)

type MaterialService struct {
	repo *infra.MaterialRepo
}

func NewMaterialService(repo *infra.MaterialRepo) *MaterialService {
	return &MaterialService{repo: repo}
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

	return s.repo.UpdateFields(ctx, id, updates)
}

func (s *MaterialService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrMaterialNotFound
	}

	return s.repo.Delete(ctx, id)
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
