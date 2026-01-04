package application

import (
	"context"
	"log"

	"back/internal/process/domain"
	"back/internal/process/infra"
	"back/pkg/es"
)

type ProcessService struct {
	repo   *infra.ProcessRepo
	esSync *es.ESSync
}

func NewProcessService(repo *infra.ProcessRepo, esSync *es.ESSync) *ProcessService {
	return &ProcessService{
		repo:   repo,
		esSync: esSync,
	}
}

func (s *ProcessService) Create(ctx context.Context, req *CreateProcessRequest) (*ProcessResponse, error) {
	process := &domain.Process{
		ProcessCode: req.ProcessCode,
		ProcessName: req.ProcessName,
		ProcessCategory: req.ProcessCategory,
		CreatedBy: req.CreatedBy,
	}

	if err := s.repo.Create(ctx, process); err != nil {
		return nil, err
	}

	// 同步到 ES
	if s.esSync != nil {
		if err := s.esSync.Index(process); err != nil {
			log.Printf("Failed to sync process to ES: %v", err)
			// 不影响主流程，只记录错误
		}
	}

	return toProcessResponse(process), nil
}

func (s *ProcessService) Get(ctx context.Context, id uint) (*ProcessResponse, error) {
	process, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrProcessNotFound
	}
	return toProcessResponse(process), nil
}

func (s *ProcessService) List(ctx context.Context, limit, offset int) ([]*ProcessResponse, int64, error) {
	processs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*ProcessResponse, len(processs))
	for i, process := range processs {
		responses[i] = toProcessResponse(&process)
	}

	return responses, count, nil
}

func (s *ProcessService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProcessNotFound
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
	if s.esSync != nil {
		process, err := s.repo.GetByID(ctx, id)
		if err == nil {
			if err := s.esSync.Update(process); err != nil {
				log.Printf("Failed to sync process update to ES: %v", err)
				// 不影响主流程，只记录错误
			}
		}
	}

	return nil
}

func (s *ProcessService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProcessNotFound
	}

	// 先获取process以获取索引名称
	process, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 从 ES 删除
	if s.esSync != nil {
		if err := s.esSync.Delete(process.GetIndexName(), process.GetDocumentID()); err != nil {
			log.Printf("Failed to delete process from ES: %v", err)
			// 不影响主流程，只记录错误
		}
	}

	return nil
}

func toProcessResponse(process *domain.Process) *ProcessResponse {
	return &ProcessResponse{
		ID: process.ID,
		ProcessCode: process.ProcessCode,
		ProcessName: process.ProcessName,
		ProcessCategory: process.ProcessCategory,
		CreatedBy: process.CreatedBy,
		CreatedAt: process.CreatedAt,
		UpdatedAt: process.UpdatedAt,
	}
}
