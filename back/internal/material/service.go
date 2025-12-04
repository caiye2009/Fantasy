package material

import (
	"context"
	"strconv"
	
	"back/pkg/es"
	"back/pkg/repo"
	"gorm.io/gorm"
)

type MaterialService struct {
	db     *gorm.DB
	esSync *es.ESSync
}

func NewMaterialService(db *gorm.DB, esSync *es.ESSync) *MaterialService {
	return &MaterialService{
		db:     db,
		esSync: esSync,
	}
}

// Create 创建材料
func (s *MaterialService) Create(ctx context.Context, req *CreateMaterialRequest) (*Material, error) {
	material := &Material{
		Name:        req.Name,
		Spec:        req.Spec,
		Unit:        req.Unit,
		Description: req.Description,
	}

	// 使用通用 repo
	materialRepo := repo.NewRepo[Material](s.db)
	if err := materialRepo.Create(ctx, material); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(material)

	return material, nil
}

// Get 获取材料
func (s *MaterialService) Get(ctx context.Context, id uint) (*Material, error) {
	materialRepo := repo.NewRepo[Material](s.db)
	return materialRepo.GetByID(ctx, id)
}

// List 材料列表
func (s *MaterialService) List(ctx context.Context, limit, offset int) ([]Material, error) {
	materialRepo := repo.NewRepo[Material](s.db)
	return materialRepo.List(ctx, limit, offset)
}

// Update 更新材料
func (s *MaterialService) Update(ctx context.Context, id uint, req *UpdateMaterialRequest) error {
	materialRepo := repo.NewRepo[Material](s.db)
	
	// 1. 查询材料
	material, err := materialRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 构造更新字段
	fields := make(map[string]interface{})
	if req.Name != "" {
		fields["name"] = req.Name
		material.Name = req.Name
	}
	if req.Spec != "" {
		fields["spec"] = req.Spec
		material.Spec = req.Spec
	}
	if req.Unit != "" {
		fields["unit"] = req.Unit
		material.Unit = req.Unit
	}
	if req.Description != "" {
		fields["description"] = req.Description
		material.Description = req.Description
	}

	// 3. 如果没有要更新的字段，直接返回
	if len(fields) == 0 {
		return nil
	}

	// 4. 更新数据库
	if err := materialRepo.UpdateFields(ctx, id, fields); err != nil {
		return err
	}

	// 5. 异步同步到 ES
	s.esSync.Update(material)

	return nil
}

// Delete 删除材料
func (s *MaterialService) Delete(ctx context.Context, id uint) error {
	materialRepo := repo.NewRepo[Material](s.db)
	
	// 1. 删除数据库记录
	if err := materialRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 2. 异步删除 ES 文档
	s.esSync.Delete("materials", strconv.Itoa(int(id)))

	return nil
}