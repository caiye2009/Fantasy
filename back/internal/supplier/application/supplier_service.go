package application

import (
	"context"

	"back/internal/supplier/domain"
	"back/internal/supplier/infra"
)

// SupplierService 供应商应用服务
type SupplierService struct {
	repo *infra.SupplierRepo
}

// NewSupplierService 创建服务
func NewSupplierService(repo *infra.SupplierRepo) *SupplierService {
	return &SupplierService{repo: repo}
}

// Create 创建供应商
func (s *SupplierService) Create(ctx context.Context, req *CreateSupplierRequest) (*SupplierResponse, error) {
	supplier := &domain.Supplier{
		SupplierCode:     req.SupplierCode,
		SupplierName:     req.SupplierName,
		SupplierCategory: req.SupplierCategory,
		CreatedBy:        req.CreatedBy,
	}

	if err := s.repo.Create(ctx, supplier); err != nil {
		return nil, err
	}

	return toSupplierResponse(supplier), nil
}

// Get 获取供应商
func (s *SupplierService) Get(ctx context.Context, id uint) (*SupplierResponse, error) {
	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrSupplierNotFound
	}
	return toSupplierResponse(supplier), nil
}

// List 列表查询
func (s *SupplierService) List(ctx context.Context, limit, offset int) ([]*SupplierResponse, int64, error) {
	suppliers, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		responses[i] = toSupplierResponse(&supplier)
	}

	return responses, count, nil
}

// Update 部分更新
func (s *SupplierService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrSupplierNotFound
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

// Delete 删除供应商
func (s *SupplierService) Delete(ctx context.Context, id uint) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrSupplierNotFound
	}

	return s.repo.Delete(ctx, id)
}

// DTO 转换
func toSupplierResponse(supplier *domain.Supplier) *SupplierResponse {
	return &SupplierResponse{
		ID:               supplier.ID,
		SupplierCode:     supplier.SupplierCode,
		SupplierName:     supplier.SupplierName,
		SupplierCategory: supplier.SupplierCategory,
		CreatedBy:        supplier.CreatedBy,
		CreatedAt:        supplier.CreatedAt,
		UpdatedAt:        supplier.UpdatedAt,
	}
}
