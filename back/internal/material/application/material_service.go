package application

import (
	"context"
	"strconv"

	"back/internal/material/domain"
	"back/internal/material/infra"
)

// ESSync ES 同步接口
type ESSync interface {
	Index(doc interface{}) error
	Update(doc interface{}) error
	Delete(indexName, docID string) error
}

// MaterialService 材料应用服务
type MaterialService struct {
	repo   *infra.MaterialRepo
	esSync ESSync
}

// NewMaterialService 创建材料服务
func NewMaterialService(repo *infra.MaterialRepo, esSync ESSync) *MaterialService {
	return &MaterialService{
		repo:   repo,
		esSync: esSync,
	}
}

// Create 创建材料
func (s *MaterialService) Create(ctx context.Context, req *CreateMaterialRequest) (*MaterialResponse, error) {
	// 1. DTO → Domain Model
	material := ToMaterial(req)
	
	// 2. 领域验证
	if err := material.Validate(); err != nil {
		return nil, err
	}
	
	// 3. 保存到数据库
	if err := s.repo.Save(ctx, material); err != nil {
		return nil, err
	}
	
	// 4. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Index(material)
	}
	
	// 5. Domain Model → DTO
	return ToMaterialResponse(material), nil
}

// Get 获取材料
func (s *MaterialService) Get(ctx context.Context, id uint) (*MaterialResponse, error) {
	material, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return ToMaterialResponse(material), nil
}

// Update 更新材料
func (s *MaterialService) Update(ctx context.Context, id uint, req *UpdateMaterialRequest) error {
	// 1. 查询材料
	material, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 更新字段（通过领域方法）
	if req.Name != "" {
		if err := material.UpdateName(req.Name); err != nil {
			return err
		}
	}
	
	if req.Spec != "" {
		if err := material.UpdateSpec(req.Spec); err != nil {
			return err
		}
	}
	
	if req.Unit != "" {
		material.Unit = req.Unit
	}
	
	if req.Description != "" {
		material.Description = req.Description
	}
	
	// 3. 验证
	if err := material.Validate(); err != nil {
		return err
	}
	
	// 4. 保存
	if err := s.repo.Update(ctx, material); err != nil {
		return err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Update(material)
	}
	
	return nil
}

// Delete 删除材料
func (s *MaterialService) Delete(ctx context.Context, id uint) error {
	// 1. 检查是否存在
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrMaterialNotFound
	}
	
	// 2. 删除
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	// 3. 异步删除 ES 文档
	if s.esSync != nil {
		s.esSync.Delete("materials", strconv.Itoa(int(id)))
	}
	
	return nil
}

// Exists 检查材料是否存在（供其他模块调用）
func (s *MaterialService) Exists(ctx context.Context, id uint) (bool, error) {
	return s.repo.ExistsByID(ctx, id)
}