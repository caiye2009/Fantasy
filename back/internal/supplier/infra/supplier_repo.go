package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"back/internal/supplier/domain"
	"back/pkg/repo"
)

// SupplierRepo 供应商仓储实现
type SupplierRepo struct {
	*repo.Repo[domain.Supplier]
	db *gorm.DB
}

// NewSupplierRepo 创建仓储
func NewSupplierRepo(db *gorm.DB) *SupplierRepo {
	return &SupplierRepo{
		Repo: repo.NewRepo[domain.Supplier](db),
		db:   db,
	}
}

// Save 保存供应商
func (r *SupplierRepo) Save(ctx context.Context, supplier *domain.Supplier) error {
	return r.Create(ctx, supplier)
}

// FindByID 根据 ID 查询
func (r *SupplierRepo) FindByID(ctx context.Context, id uint) (*domain.Supplier, error) {
	supplier, err := r.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSupplierNotFound
	}
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

// FindAll 查询所有
func (r *SupplierRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.Supplier, error) {
	suppliers, err := r.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Supplier, len(suppliers))
	for i := range suppliers {
		result[i] = &suppliers[i]
	}
	return result, nil
}

// Update 更新供应商
func (r *SupplierRepo) Update(ctx context.Context, supplier *domain.Supplier) error {
	return r.Repo.Update(ctx, supplier)
}

// Delete 删除供应商
func (r *SupplierRepo) Delete(ctx context.Context, id uint) error {
	return r.Repo.Delete(ctx, id)
}

// ExistsByID 检查是否存在
func (r *SupplierRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"id": id})
}

// FindByName 根据名称查询
func (r *SupplierRepo) FindByName(ctx context.Context, name string) (*domain.Supplier, error) {
	supplier, err := r.First(ctx, map[string]interface{}{"name": name})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSupplierNotFound
	}
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

// ExistsByName 检查名称是否存在
func (r *SupplierRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"name": name})
}

// FindByPhone 根据电话查询
func (r *SupplierRepo) FindByPhone(ctx context.Context, phone string) (*domain.Supplier, error) {
	supplier, err := r.First(ctx, map[string]interface{}{"phone": phone})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSupplierNotFound
	}
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

// FindByEmail 根据邮箱查询
func (r *SupplierRepo) FindByEmail(ctx context.Context, email string) (*domain.Supplier, error) {
	supplier, err := r.First(ctx, map[string]interface{}{"email": email})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSupplierNotFound
	}
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

// Count 统计数量
func (r *SupplierRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}
