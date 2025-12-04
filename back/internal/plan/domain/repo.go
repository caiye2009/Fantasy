package domain

import "context"

// PlanRepository 计划仓储接口
type PlanRepository interface {
	// 基础 CRUD
	Save(ctx context.Context, plan *Plan) error
	FindByID(ctx context.Context, id uint) (*Plan, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Plan, error)
	Update(ctx context.Context, plan *Plan) error
	Delete(ctx context.Context, id uint) error
	
	// 业务查询
	ExistsByID(ctx context.Context, id uint) (bool, error)
	FindByPlanNo(ctx context.Context, planNo string) (*Plan, error)
	ExistsByPlanNo(ctx context.Context, planNo string) (bool, error)
	FindByOrderID(ctx context.Context, orderID uint) ([]*Plan, error)
	FindByProductID(ctx context.Context, productID uint) ([]*Plan, error)
	FindByStatus(ctx context.Context, status PlanStatus, limit, offset int) ([]*Plan, error)
	Count(ctx context.Context) (int64, error)
}