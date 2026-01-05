package application

import (
	"context"

	"back/internal/product_fabric/domain"
	"back/internal/product_fabric/infra"
	"back/pkg/auth"
)

type ProductFabricService struct {
	repo *infra.ProductFabricRepo
}

func NewProductFabricService(repo *infra.ProductFabricRepo) *ProductFabricService {
	return &ProductFabricService{repo: repo}
}

func (s *ProductFabricService) Create(ctx context.Context, req *CreateProductFabricRequest) (*ProductFabricResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	productFabric := &domain.ProductFabric{
		ProductCode: req.ProductCode,
		FabricCode:  req.FabricCode,
		RequiredQty: req.RequiredQty,
		CreatedBy:   createdBy,
	}

	if err := s.repo.Create(ctx, productFabric); err != nil {
		return nil, err
	}

	return toProductFabricResponse(productFabric), nil
}

func (s *ProductFabricService) Get(ctx context.Context, id uint) (*ProductFabricResponse, error) {
	productFabric, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductFabricNotFound
	}
	return toProductFabricResponse(productFabric), nil
}

func (s *ProductFabricService) List(ctx context.Context, limit, offset int) ([]*ProductFabricResponse, int64, error) {
	productFabrics, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*ProductFabricResponse, len(productFabrics))
	for i, productFabric := range productFabrics {
		responses[i] = toProductFabricResponse(&productFabric)
	}

	return responses, count, nil
}

func (s *ProductFabricService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductFabricNotFound
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

func (s *ProductFabricService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductFabricNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toProductFabricResponse(productFabric *domain.ProductFabric) *ProductFabricResponse {
	return &ProductFabricResponse{
		ID:          productFabric.ID,
		ProductCode: productFabric.ProductCode,
		FabricCode:  productFabric.FabricCode,
		RequiredQty: productFabric.RequiredQty,
		CreatedBy:   productFabric.CreatedBy,
		CreatedAt:   productFabric.CreatedAt,
		UpdatedAt:   productFabric.UpdatedAt,
	}
}