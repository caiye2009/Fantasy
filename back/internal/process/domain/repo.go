package domain

import "context"

// ProcessRepository 工序仓储接口
type ProcessRepository interface {
	// 基础 CRUD
	Save(ctx context.Context, process *Process) error
	FindByID(ctx context.Context, id uint) (*Process, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Process, error)
	Update(ctx context.Context, process *Process) error
	Delete(ctx context.Context, id uint) error
	
	// 业务查询
	ExistsByID(ctx context.Context, id uint) (bool, error)
	FindByName(ctx context.Context, name string) (*Process, error)
	Count(ctx context.Context) (int64, error)
}