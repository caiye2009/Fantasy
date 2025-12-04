package domain

import "context"

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// 基础 CRUD
	Save(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id uint) (*Order, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uint) error
	
	// 业务查询
	ExistsByID(ctx context.Context, id uint) (bool, error)
	FindByOrderNo(ctx context.Context, orderNo string) (*Order, error)
	ExistsByOrderNo(ctx context.Context, orderNo string) (bool, error)
	FindByClientID(ctx context.Context, clientID uint) ([]*Order, error)
	FindByProductID(ctx context.Context, productID uint) ([]*Order, error)
	FindByStatus(ctx context.Context, status OrderStatus, limit, offset int) ([]*Order, error)
	Count(ctx context.Context) (int64, error)
}