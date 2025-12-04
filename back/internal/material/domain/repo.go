package domain

import "context"

// MaterialRepository 材料仓储接口
type MaterialRepository interface {
	// 基础 CRUD
	Save(ctx context.Context, material *Material) error
	FindByID(ctx context.Context, id uint) (*Material, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Material, error)
	Update(ctx context.Context, material *Material) error
	Delete(ctx context.Context, id uint) error
	
	// 业务查询
	ExistsByID(ctx context.Context, id uint) (bool, error)
	FindByName(ctx context.Context, name string) (*Material, error)
	Count(ctx context.Context) (int64, error)
}