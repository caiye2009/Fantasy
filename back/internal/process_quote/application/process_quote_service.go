package application

import (
	"context"

	"back/internal/process_quote/domain"
	"back/internal/process_quote/infra"
	"back/pkg/auth"
)

type ProcessQuoteService struct {
	repo *infra.ProcessQuoteRepo
}

func NewProcessQuoteService(repo *infra.ProcessQuoteRepo) *ProcessQuoteService {
	return &ProcessQuoteService{repo: repo}
}

func (s *ProcessQuoteService) Create(ctx context.Context, req *CreateProcessQuoteRequest) (*ProcessQuoteResponse, error) {
	createdBy := auth.GetLoginIDFromContext(ctx)

	processQuote := &domain.ProcessQuote{
		QuoteID: req.QuoteID,
		ProcessCode: req.ProcessCode,
		SupplierCode: req.SupplierCode,
		QuotePrice: req.QuotePrice,
		CreatedBy: createdBy,
	}

	if err := s.repo.Create(ctx, processQuote); err != nil {
		return nil, err
	}

	return toProcessQuoteResponse(processQuote), nil
}

func (s *ProcessQuoteService) Get(ctx context.Context, id uint) (*ProcessQuoteResponse, error) {
	processQuote, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProcessQuoteNotFound
	}
	return toProcessQuoteResponse(processQuote), nil
}

func (s *ProcessQuoteService) List(ctx context.Context, limit, offset int) ([]*ProcessQuoteResponse, int64, error) {
	processQuotes, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*ProcessQuoteResponse, len(processQuotes))
	for i, processQuote := range processQuotes {
		responses[i] = toProcessQuoteResponse(&processQuote)
	}

	return responses, count, nil
}

func (s *ProcessQuoteService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProcessQuoteNotFound
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

func (s *ProcessQuoteService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProcessQuoteNotFound
	}

	return s.repo.Delete(ctx, id)
}

func toProcessQuoteResponse(processQuote *domain.ProcessQuote) *ProcessQuoteResponse {
	return &ProcessQuoteResponse{
		ID: processQuote.ID,
		QuoteID: processQuote.QuoteID,
		ProcessCode: processQuote.ProcessCode,
		SupplierCode: processQuote.SupplierCode,
		QuotePrice: processQuote.QuotePrice,
		CreatedBy: processQuote.CreatedBy,
		CreatedAt: processQuote.CreatedAt,
		UpdatedAt: processQuote.UpdatedAt,
	}
}
