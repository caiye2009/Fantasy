package application

import (
	"context"

	"back/internal/product/domain"
	"back/internal/product/infra"
	"back/pkg/auth"
	"back/pkg/es"
)

type ProductService struct {
	repo   *infra.ProductRepo
	syncer *es.EntitySyncer[*domain.Product]
}

func NewProductService(repo *infra.ProductRepo, esSync *es.ESSync) *ProductService {
	return &ProductService{
		repo:   repo,
		syncer: es.NewEntitySyncer[*domain.Product](esSync),
	}
}

func (s *ProductService) Create(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	product := &domain.Product{
		ProductCode: req.ProductCode,
		ProductName: req.ProductName,
		CreatedBy:   createdBy,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	s.syncer.SyncCreate(product)

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

	if err := s.repo.UpdateFields(ctx, id, updates); err != nil {
		return err
	}

	s.syncer.SyncUpdateWithFetch(ctx, id, s.repo.GetByID)

	return nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductNotFound
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.syncer.SyncDelete(product)

	return nil
}

func toProductResponse(product *domain.Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID,
		ProductCode: product.ProductCode,
		ProductName: product.ProductName,
		CreatedBy:   product.CreatedBy,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}