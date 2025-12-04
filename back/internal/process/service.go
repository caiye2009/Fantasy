package process

import (
	"context"
	"strconv"
	
	"back/pkg/es"
	"back/pkg/repo"
	"gorm.io/gorm"
)

type ProcessService struct {
	db     *gorm.DB
	esSync *es.ESSync
}

func NewProcessService(db *gorm.DB, esSync *es.ESSync) *ProcessService {
	return &ProcessService{
		db:     db,
		esSync: esSync,
	}
}

// Create 创建工序
func (s *ProcessService) Create(ctx context.Context, req *CreateProcessRequest) (*Process, error) {
	process := &Process{
		Name:        req.Name,
		Description: req.Description,
	}

	processRepo := repo.NewRepo[Process](s.db)
	if err := processRepo.Create(ctx, process); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(process)

	return process, nil
}

// Get 获取工序
func (s *ProcessService) Get(ctx context.Context, id uint) (*Process, error) {
	processRepo := repo.NewRepo[Process](s.db)
	return processRepo.GetByID(ctx, id)
}

// List 工序列表
func (s *ProcessService) List(ctx context.Context, limit, offset int) ([]Process, error) {
	processRepo := repo.NewRepo[Process](s.db)
	return processRepo.List(ctx, limit, offset)
}

// Update 更新工序
func (s *ProcessService) Update(ctx context.Context, id uint, req *UpdateProcessRequest) error {
	processRepo := repo.NewRepo[Process](s.db)
	
	// 1. 查询工序
	process, err := processRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 构造更新字段
	fields := make(map[string]interface{})
	if req.Name != "" {
		fields["name"] = req.Name
		process.Name = req.Name
	}
	if req.Description != "" {
		fields["description"] = req.Description
		process.Description = req.Description
	}

	// 3. 如果没有要更新的字段，直接返回
	if len(fields) == 0 {
		return nil
	}

	// 4. 更新数据库
	if err := processRepo.UpdateFields(ctx, id, fields); err != nil {
		return err
	}

	// 5. 异步同步到 ES
	s.esSync.Update(process)

	return nil
}

// Delete 删除工序
func (s *ProcessService) Delete(ctx context.Context, id uint) error {
	processRepo := repo.NewRepo[Process](s.db)
	
	if err := processRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("processes", strconv.Itoa(int(id)))

	return nil
}