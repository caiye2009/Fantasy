package application

import (
	"context"
	"strconv"
	
	"back/internal/supplier/domain"
)

// ESSync ES 同步接口
type ESSync interface {
	Index(doc interface{}) error
	Update(doc interface{}) error
	Delete(indexName, docID string) error
}

// SupplierService 供应商应用服务
type SupplierService struct {
	repo   domain.SupplierRepository
	esSync ESSync
}

// NewSupplierService 创建供应商服务
func NewSupplierService(repo domain.SupplierRepository, esSync ESSync) *SupplierService {
	return &SupplierService{
		repo:   repo,
		esSync: esSync,
	}
}

// Create 创建供应商
func (s *SupplierService) Create(ctx context.Context, req *CreateSupplierRequest) (*SupplierResponse, error) {
	// 1. 检查名称是否重复
	exists, err := s.repo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrSupplierNameExists
	}
	
	// 2. DTO → Domain Model
	supplier := ToSupplier(req)
	
	// 3. 领域验证
	if err := supplier.Validate(); err != nil {
		return nil, err
	}
	
	// 4. 保存到数据库
	if err := s.repo.Save(ctx, supplier); err != nil {
		return nil, err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Index(supplier)
	}
	
	// 6. Domain Model → DTO
	return ToSupplierResponse(supplier), nil
}

// Get 获取供应商
func (s *SupplierService) Get(ctx context.Context, id uint) (*SupplierResponse, error) {
	supplier, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return ToSupplierResponse(supplier), nil
}

// List 供应商列表
func (s *SupplierService) List(ctx context.Context, limit, offset int) (*SupplierListResponse, error) {
	suppliers, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	
	return ToSupplierListResponse(suppliers, total), nil
}

// Update 更新供应商
func (s *SupplierService) Update(ctx context.Context, id uint, req *UpdateSupplierRequest) error {
	// 1. 查询供应商
	supplier, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 更新字段（通过领域方法）
	if req.Name != "" {
		// 检查新名称是否重复
		if req.Name != supplier.Name {
			exists, err := s.repo.ExistsByName(ctx, req.Name)
			if err != nil {
				return err
			}
			if exists {
				return domain.ErrSupplierNameExists
			}
		}
		
		if err := supplier.UpdateName(req.Name); err != nil {
			return err
		}
	}
	
	if req.Contact != "" {
		if err := supplier.UpdateContact(req.Contact); err != nil {
			return err
		}
	}
	
	if req.Phone != "" {
		if err := supplier.UpdatePhone(req.Phone); err != nil {
			return err
		}
	}
	
	if req.Email != "" {
		if err := supplier.UpdateEmail(req.Email); err != nil {
			return err
		}
	}
	
	if req.Address != "" {
		if err := supplier.UpdateAddress(req.Address); err != nil {
			return err
		}
	}
	
	// 3. 验证
	if err := supplier.Validate(); err != nil {
		return err
	}
	
	// 4. 保存
	if err := s.repo.Update(ctx, supplier); err != nil {
		return err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Update(supplier)
	}
	
	return nil
}

// Delete 删除供应商
func (s *SupplierService) Delete(ctx context.Context, id uint) error {
	// 1. 检查是否存在
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrSupplierNotFound
	}
	
	// 2. 删除
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	// 3. 异步删除 ES 文档
	if s.esSync != nil {
		s.esSync.Delete("suppliers", strconv.Itoa(int(id)))
	}
	
	return nil
}

// Exists 检查供应商是否存在（供其他模块调用）
func (s *SupplierService) Exists(ctx context.Context, id uint) (bool, error) {
	return s.repo.ExistsByID(ctx, id)
}

// GetSupplierInfo 获取供应商基本信息（供 Pricing 模块调用）
func (s *SupplierService) GetSupplierInfo(ctx context.Context, id uint) (*SupplierInfo, error) {
	supplier, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &SupplierInfo{
		ID:   supplier.ID,
		Name: supplier.Name,
	}, nil
}

// SupplierInfo 供应商基本信息（供其他模块使用）
type SupplierInfo struct {
	ID   uint
	Name string
}