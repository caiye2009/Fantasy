package application

import (
	"context"

	"back/internal/inventory/domain"
	"back/internal/inventory/infra"
)

// InventoryService 库存服务
type InventoryService struct {
	repo *infra.InventoryRepo
}

// NewInventoryService 创建库存服务
func NewInventoryService(repo *infra.InventoryRepo) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

// CreateInventory 创建库存
func (s *InventoryService) CreateInventory(ctx context.Context, req *CreateInventoryRequest) (*InventoryResponse, error) {
	inventory := &domain.Inventory{
		ProductID: req.ProductID,
		Category:  req.Category,
		BatchID:   req.BatchID,
		Quantity:  req.Quantity,
		Unit:      req.Unit,
		UnitCost:  req.UnitCost,
		Remark:    req.Remark,
	}

	// 计算总成本
	inventory.CalculateTotalCost()

	// 验证
	if err := inventory.Validate(); err != nil {
		return nil, err
	}

	// 保存
	if err := s.repo.Save(ctx, inventory); err != nil {
		return nil, err
	}

	return s.toResponse(inventory), nil
}

// GetInventory 获取库存详情
func (s *InventoryService) GetInventory(ctx context.Context, id uint) (*InventoryResponse, error) {
	inventory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(inventory), nil
}

// GetInventoriesByProductID 根据产品ID获取库存列表
func (s *InventoryService) GetInventoriesByProductID(ctx context.Context, productID uint) (*InventoryListResponse, error) {
	inventories, err := s.repo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	responses := make([]*InventoryResponse, len(inventories))
	for i, inv := range inventories {
		responses[i] = s.toResponse(inv)
	}

	return &InventoryListResponse{
		Total:       len(responses),
		Inventories: responses,
	}, nil
}

// GetInventoryByBatchID 根据批次ID获取库存
func (s *InventoryService) GetInventoryByBatchID(ctx context.Context, batchID string) (*InventoryResponse, error) {
	inventory, err := s.repo.FindByBatchID(ctx, batchID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(inventory), nil
}

// ListInventories 获取库存列表
func (s *InventoryService) ListInventories(ctx context.Context, limit, offset int) (*InventoryListResponse, error) {
	inventories, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*InventoryResponse, len(inventories))
	for i, inv := range inventories {
		responses[i] = s.toResponse(inv)
	}

	return &InventoryListResponse{
		Total:       len(responses),
		Inventories: responses,
	}, nil
}

// ListInventoriesByCategory 根据类别获取库存列表
func (s *InventoryService) ListInventoriesByCategory(ctx context.Context, category string, limit, offset int) (*InventoryListResponse, error) {
	inventories, err := s.repo.FindByCategory(ctx, category, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*InventoryResponse, len(inventories))
	for i, inv := range inventories {
		responses[i] = s.toResponse(inv)
	}

	return &InventoryListResponse{
		Total:       len(responses),
		Inventories: responses,
	}, nil
}

// UpdateInventory 更新库存
func (s *InventoryService) UpdateInventory(ctx context.Context, req *UpdateInventoryRequest) (*InventoryResponse, error) {
	// 查询现有库存
	inventory, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	inventory.ProductID = req.ProductID
	inventory.Category = req.Category
	inventory.BatchID = req.BatchID
	inventory.Quantity = req.Quantity
	inventory.Unit = req.Unit
	inventory.UnitCost = req.UnitCost
	inventory.Remark = req.Remark

	// 计算总成本
	inventory.CalculateTotalCost()

	// 验证
	if err := inventory.Validate(); err != nil {
		return nil, err
	}

	// 保存
	if err := s.repo.Update(ctx, inventory); err != nil {
		return nil, err
	}

	return s.toResponse(inventory), nil
}

// UpdateQuantity 更新数量
func (s *InventoryService) UpdateQuantity(ctx context.Context, req *UpdateQuantityRequest) (*InventoryResponse, error) {
	inventory, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if err := inventory.UpdateQuantity(req.Quantity); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, inventory); err != nil {
		return nil, err
	}

	return s.toResponse(inventory), nil
}

// DeductInventory 扣减库存
func (s *InventoryService) DeductInventory(ctx context.Context, req *DeductInventoryRequest) (*InventoryResponse, error) {
	inventory, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if err := inventory.Deduct(req.Quantity); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, inventory); err != nil {
		return nil, err
	}

	return s.toResponse(inventory), nil
}

// AddInventory 增加库存
func (s *InventoryService) AddInventory(ctx context.Context, req *AddInventoryRequest) (*InventoryResponse, error) {
	inventory, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if err := inventory.Add(req.Quantity); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, inventory); err != nil {
		return nil, err
	}

	return s.toResponse(inventory), nil
}

// DeleteInventory 删除库存
func (s *InventoryService) DeleteInventory(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// toResponse 转换为响应
func (s *InventoryService) toResponse(inventory *domain.Inventory) *InventoryResponse {
	return &InventoryResponse{
		ID:        inventory.ID,
		ProductID: inventory.ProductID,
		Category:  inventory.Category,
		BatchID:   inventory.BatchID,
		Quantity:  inventory.Quantity,
		Unit:      inventory.Unit,
		UnitCost:  inventory.UnitCost,
		TotalCost: inventory.TotalCost,
		Remark:    inventory.Remark,
		CreatedAt: inventory.CreatedAt,
		UpdatedAt: inventory.UpdatedAt,
	}
}
