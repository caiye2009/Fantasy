package order

import (
	"context"
	"strconv"
	
	"back/pkg/es"
	"back/pkg/repo"
	"gorm.io/gorm"
)

type OrderService struct {
	db     *gorm.DB
	esSync *es.ESSync
}

func NewOrderService(db *gorm.DB, esSync *es.ESSync) *OrderService {
	return &OrderService{
		db:     db,
		esSync: esSync,
	}
}

// Create 创建订单
func (s *OrderService) Create(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
	order := &Order{
		OrderNo:    req.OrderNo,
		ClientID:   req.ClientID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		UnitPrice:  req.UnitPrice,
		TotalPrice: req.Quantity * req.UnitPrice,
		Status:     "pending",
		CreatedBy:  req.CreatedBy,
	}

	orderRepo := repo.NewRepo[Order](s.db)
	if err := orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(order)

	return order, nil
}

// Get 获取订单
func (s *OrderService) Get(ctx context.Context, id uint) (*Order, error) {
	orderRepo := repo.NewRepo[Order](s.db)
	return orderRepo.GetByID(ctx, id)
}

// List 订单列表
func (s *OrderService) List(ctx context.Context, limit, offset int) ([]Order, error) {
	orderRepo := repo.NewRepo[Order](s.db)
	return orderRepo.List(ctx, limit, offset)
}

// Update 更新订单
func (s *OrderService) Update(ctx context.Context, id uint, req *UpdateOrderRequest) error {
	orderRepo := repo.NewRepo[Order](s.db)
	
	// 1. 查询订单
	order, err := orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. 构造更新字段
	fields := make(map[string]interface{})
	if req.Status != "" {
		fields["status"] = req.Status
		order.Status = req.Status
	}
	if req.Quantity > 0 {
		fields["quantity"] = req.Quantity
		order.Quantity = req.Quantity
		order.TotalPrice = req.Quantity * order.UnitPrice
		fields["total_price"] = order.TotalPrice
	}

	// 3. 如果没有要更新的字段，直接返回
	if len(fields) == 0 {
		return nil
	}

	// 4. 更新数据库
	if err := orderRepo.UpdateFields(ctx, id, fields); err != nil {
		return err
	}

	// 5. 异步同步到 ES
	s.esSync.Update(order)

	return nil
}

// Delete 删除订单
func (s *OrderService) Delete(ctx context.Context, id uint) error {
	orderRepo := repo.NewRepo[Order](s.db)
	
	if err := orderRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("orders", strconv.Itoa(int(id)))

	return nil
}