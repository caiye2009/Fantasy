package application

import (
	"context"

	"back/internal/product_material/domain"
	"back/internal/product_material/infra"
)

type ProductMaterialService struct {
	repo *infra.ProductMaterialRepo
}

func NewProductMaterialService(repo *infra.ProductMaterialRepo) *ProductMaterialService {
	return &ProductMaterialService{repo: repo}
}

func (s *ProductMaterialService) Create(ctx context.Context, req *CreateProductMaterialRequest) (*ProductMaterialResponse, error) {
	productMaterial := &domain.ProductMaterial{
		ProductCode: req.ProductCode,
		MaterialCode: req.MaterialCode,
		RequiredQty: req.RequiredQty,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, productMaterial); err != nil {
		return nil, err
	}

	return toProductMaterialResponse(productMaterial), nil
}

func (s *ProductMaterialService) Get(ctx context.Context, id uint) (*ProductMaterialResponse, error) {
	productMaterial, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductMaterialNotFound
	}
	return toProductMaterialResponse(productMaterial), nil
}

func (s *ProductMaterialService) List(ctx context.Context, limit, offset int) ([]*ProductMaterialResponse, int64, error) {
	productMaterials, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*ProductMaterialResponse, len(productMaterials))
	for i, productMaterial := range productMaterials {
		responses[i] = toProductMaterialResponse(&productMaterial)
	}

	return responses, count, nil
}

func (s *ProductMaterialService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductMaterialNotFound
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

func (s *ProductMaterialService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductMaterialNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toProductMaterialResponse(productMaterial *domain.ProductMaterial) *ProductMaterialResponse {
	return &ProductMaterialResponse{
		ID: productMaterial.ID,
		ProductCode: productMaterial.ProductCode,
		MaterialCode: productMaterial.MaterialCode,
		RequiredQty: productMaterial.RequiredQty,
		CreatedBy: productMaterial.CreatedBy,
		CreatedAt: productMaterial.CreatedAt,
		UpdatedAt: productMaterial.UpdatedAt,
	}
}
