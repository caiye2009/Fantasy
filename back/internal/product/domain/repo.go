package domain

import "context"

// ProductRepository 产品仓储接口
type ProductRepository interface {
	// 基础 CRUD
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id uint) (*Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uint) error
	
	// 业务查询
	ExistsByID(ctx context.Context, id uint) (bool, error)
	FindByName(ctx context.Context, name string) (*Product, error)
	FindByStatus(ctx context.Context, status ProductStatus, limit, offset int) ([]*Product, error)
	Count(ctx context.Context) (int64, error)
}