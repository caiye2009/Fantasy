package application

import (
	"context"

	"back/internal/product/domain"
	"back/internal/product/infra"
)

type ProductService struct {
	repo *infra.ProductRepo
}

func NewProductService(repo *infra.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	product := &domain.Product{
		ProductCode: req.ProductCode,
		ProductName: req.ProductName,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return toProductResponse(product), nil
}

func (s *ProductService) Get(ctx context.Context, id uint) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductNotFound
	}
	return toProductResponse(product), nil
}

func (s *ProductService) List(ctx context.Context, limit, offset int) ([]*ProductResponse, int64, error) {
	products, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = toProductResponse(&product)
	}

	return responses, count, nil
}

func (s *ProductService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductNotFound
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

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toProductResponse(product *domain.Product) *ProductResponse {
	return &ProductResponse{
		ID: product.ID,
		ProductCode: product.ProductCode,
		ProductName: product.ProductName,
		CreatedBy: product.CreatedBy,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
