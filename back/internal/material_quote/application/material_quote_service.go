package application

import (
	"context"

	"back/internal/material_quote/domain"
	"back/internal/material_quote/infra"
)

type MaterialQuoteService struct {
	repo *infra.MaterialQuoteRepo
}

func NewMaterialQuoteService(repo *infra.MaterialQuoteRepo) *MaterialQuoteService {
	return &MaterialQuoteService{repo: repo}
}

func (s *MaterialQuoteService) Create(ctx context.Context, req *CreateMaterialQuoteRequest) (*MaterialQuoteResponse, error) {
	materialQuote := &domain.MaterialQuote{
		QuoteID: req.QuoteID,
		MaterialCode: req.MaterialCode,
		SupplierCode: req.SupplierCode,
		QuotePrice: req.QuotePrice,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, materialQuote); err != nil {
		return nil, err
	}

	return toMaterialQuoteResponse(materialQuote), nil
}

func (s *MaterialQuoteService) Get(ctx context.Context, id uint) (*MaterialQuoteResponse, error) {
	materialQuote, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrMaterialQuoteNotFound
	}
	return toMaterialQuoteResponse(materialQuote), nil
}

func (s *MaterialQuoteService) List(ctx context.Context, limit, offset int) ([]*MaterialQuoteResponse, int64, error) {
	materialQuotes, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*MaterialQuoteResponse, len(materialQuotes))
	for i, materialQuote := range materialQuotes {
		responses[i] = toMaterialQuoteResponse(&materialQuote)
	}

	return responses, count, nil
}

func (s *MaterialQuoteService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrMaterialQuoteNotFound
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

func (s *MaterialQuoteService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrMaterialQuoteNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toMaterialQuoteResponse(materialQuote *domain.MaterialQuote) *MaterialQuoteResponse {
	return &MaterialQuoteResponse{
		ID: materialQuote.ID,
		QuoteID: materialQuote.QuoteID,
		MaterialCode: materialQuote.MaterialCode,
		SupplierCode: materialQuote.SupplierCode,
		QuotePrice: materialQuote.QuotePrice,
		CreatedBy: materialQuote.CreatedBy,
		CreatedAt: materialQuote.CreatedAt,
		UpdatedAt: materialQuote.UpdatedAt,
	}
}
