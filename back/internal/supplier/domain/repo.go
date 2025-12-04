package domain

import "context"

// SupplierRepository 供应商仓储接口
type SupplierRepository interface {
	// 基础 CRUD
	Save(ctx context.Context, supplier *Supplier) error
	FindByID(ctx context.Context, id uint) (*Supplier, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Supplier, error)
	Update(ctx context.Context, supplier *Supplier) error
	Delete(ctx context.Context, id uint) error

	// 业务查询
	ExistsByID(ctx context.Context, id uint) (bool, error)
	FindByName(ctx context.Context, name string) (*Supplier, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	FindByPhone(ctx context.Context, phone string) (*Supplier, error)
	FindByEmail(ctx context.Context, email string) (*Supplier, error)
	Count(ctx context.Context) (int64, error)
}