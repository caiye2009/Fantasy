package application

import (
	"context"

	"back/internal/product_process/domain"
	"back/internal/product_process/infra"
)

type ProductProcessService struct {
	repo *infra.ProductProcessRepo
}

func NewProductProcessService(repo *infra.ProductProcessRepo) *ProductProcessService {
	return &ProductProcessService{repo: repo}
}

func (s *ProductProcessService) Create(ctx context.Context, req *CreateProductProcessRequest) (*ProductProcessResponse, error) {
	productProcess := &domain.ProductProcess{
		ProductCode: req.ProductCode,
		ProcessCode: req.ProcessCode,
		ProcessSeq: req.ProcessSeq,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, productProcess); err != nil {
		return nil, err
	}

	return toProductProcessResponse(productProcess), nil
}

func (s *ProductProcessService) Get(ctx context.Context, id uint) (*ProductProcessResponse, error) {
	productProcess, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProductProcessNotFound
	}
	return toProductProcessResponse(productProcess), nil
}

func (s *ProductProcessService) List(ctx context.Context, limit, offset int) ([]*ProductProcessResponse, int64, error) {
	productProcesss, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*ProductProcessResponse, len(productProcesss))
	for i, productProcess := range productProcesss {
		responses[i] = toProductProcessResponse(&productProcess)
	}

	return responses, count, nil
}

func (s *ProductProcessService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductProcessNotFound
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

func (s *ProductProcessService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProductProcessNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toProductProcessResponse(productProcess *domain.ProductProcess) *ProductProcessResponse {
	return &ProductProcessResponse{
		ID: productProcess.ID,
		ProductCode: productProcess.ProductCode,
		ProcessCode: productProcess.ProcessCode,
		ProcessSeq: productProcess.ProcessSeq,
		CreatedBy: productProcess.CreatedBy,
		CreatedAt: productProcess.CreatedAt,
		UpdatedAt: productProcess.UpdatedAt,
	}
}
