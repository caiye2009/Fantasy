package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"back/internal/inventory/domain"
	"back/pkg/repo"
)

// InventoryRepo 库存仓储实现
type InventoryRepo struct {
	*repo.Repo[domain.Inventory]
	db *gorm.DB
}

// NewInventoryRepo 创建仓储
func NewInventoryRepo(db *gorm.DB) *InventoryRepo {
	return &InventoryRepo{
		Repo: repo.NewRepo[domain.Inventory](db),
		db:   db,
	}
}

// Save 保存库存
func (r *InventoryRepo) Save(ctx context.Context, inventory *domain.Inventory) error {
	return r.Create(ctx, inventory)
}

// FindByID 根据 ID 查询
func (r *InventoryRepo) FindByID(ctx context.Context, id uint) (*domain.Inventory, error) {
	inventory, err := r.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrInventoryNotFound
	}
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

// FindByProductID 根据产品 ID 查询库存列表
func (r *InventoryRepo) FindByProductID(ctx context.Context, productID uint) ([]*domain.Inventory, error) {
	var inventories []domain.Inventory
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Find(&inventories).Error

	if err != nil {
		return nil, err
	}

	result := make([]*domain.Inventory, len(inventories))
	for i := range inventories {
		result[i] = &inventories[i]
	}
	return result, nil
}

// FindByBatchID 根据批次 ID 查询
func (r *InventoryRepo) FindByBatchID(ctx context.Context, batchID string) (*domain.Inventory, error) {
	inventory, err := r.First(ctx, map[string]interface{}{"batch_id": batchID})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrInventoryNotFound
	}
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

// FindAll 查询所有
func (r *InventoryRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.Inventory, error) {
	inventories, err := r.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Inventory, len(inventories))
	for i := range inventories {
		result[i] = &inventories[i]
	}
	return result, nil
}

// Update 更新库存
func (r *InventoryRepo) Update(ctx context.Context, inventory *domain.Inventory) error {
	return r.Repo.Update(ctx, inventory)
}

// Delete 删除库存
func (r *InventoryRepo) Delete(ctx context.Context, id uint) error {
	return r.Repo.Delete(ctx, id)
}

// ExistsByID 检查是否存在
func (r *InventoryRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"id": id})
}

// FindByCategory 根据类别查询
func (r *InventoryRepo) FindByCategory(ctx context.Context, category string, limit, offset int) ([]*domain.Inventory, error) {
	var inventories []domain.Inventory
	err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Limit(limit).
		Offset(offset).
		Find(&inventories).Error

	if err != nil {
		return nil, err
	}

	result := make([]*domain.Inventory, len(inventories))
	for i := range inventories {
		result[i] = &inventories[i]
	}
	return result, nil
}
