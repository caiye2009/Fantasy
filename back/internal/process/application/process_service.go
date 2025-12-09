package application

import (
	"context"
	"strconv"
	
	"back/internal/process/domain"
)

// ESSync ES 同步接口
type ESSync interface {
	Index(doc interface{}) error
	Update(doc interface{}) error
	Delete(indexName, docID string) error
}

// ProcessService 工序应用服务
type ProcessService struct {
	repo   domain.ProcessRepository
	esSync ESSync
}

// NewProcessService 创建工序服务
func NewProcessService(repo domain.ProcessRepository, esSync ESSync) *ProcessService {
	return &ProcessService{
		repo:   repo,
		esSync: esSync,
	}
}

// Create 创建工序
func (s *ProcessService) Create(ctx context.Context, req *CreateProcessRequest) (*ProcessResponse, error) {
	// 1. DTO → Domain Model
	process := ToProcess(req)
	
	// 2. 领域验证
	if err := process.Validate(); err != nil {
		return nil, err
	}
	
	// 3. 保存到数据库
	if err := s.repo.Save(ctx, process); err != nil {
		return nil, err
	}
	
	// 4. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Index(process)
	}
	
	// 5. Domain Model → DTO
	return ToProcessResponse(process), nil
}

// Get 获取工序
func (s *ProcessService) Get(ctx context.Context, id uint) (*ProcessResponse, error) {
	process, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return ToProcessResponse(process), nil
}

// List 功能已移至 search 模块，通过 ES 实现
// 使用 POST /api/v1/search 并指定 indices: ["processes"]

// Update 更新工序
func (s *ProcessService) Update(ctx context.Context, id uint, req *UpdateProcessRequest) error {
	// 1. 查询工序
	process, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 更新字段（通过领域方法）
	if req.Name != "" {
		if err := process.UpdateName(req.Name); err != nil {
			return err
		}
	}
	
	if req.Description != "" {
		process.Description = req.Description
	}
	
	// 3. 验证
	if err := process.Validate(); err != nil {
		return err
	}
	
	// 4. 保存
	if err := s.repo.Update(ctx, process); err != nil {
		return err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Update(process)
	}
	
	return nil
}

// Delete 删除工序
func (s *ProcessService) Delete(ctx context.Context, id uint) error {
	// 1. 检查是否存在
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrProcessNotFound
	}
	
	// 2. 删除
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	// 3. 异步删除 ES 文档
	if s.esSync != nil {
		s.esSync.Delete("processes", strconv.Itoa(int(id)))
	}
	
	return nil
}

// Exists 检查工序是否存在（供其他模块调用）
func (s *ProcessService) Exists(ctx context.Context, id uint) (bool, error) {
	return s.repo.ExistsByID(ctx, id)
}