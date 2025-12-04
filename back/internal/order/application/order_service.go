package application

import (
	"context"
	"strconv"
	
	"back/internal/order/domain"
)

// ESSync ES 同步接口
type ESSync interface {
	Index(doc interface{}) error
	Update(doc interface{}) error
	Delete(indexName, docID string) error
}

// OrderService 订单应用服务
type OrderService struct {
	repo   domain.OrderRepository
	esSync ESSync
}

// NewOrderService 创建订单服务
func NewOrderService(repo domain.OrderRepository, esSync ESSync) *OrderService {
	return &OrderService{
		repo:   repo,
		esSync: esSync,
	}
}

// Create 创建订单
func (s *OrderService) Create(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error) {
	// 1. 检查订单编号是否重复
	exists, err := s.repo.ExistsByOrderNo(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrOrderNoDuplicate
	}
	
	// 2. DTO → Domain Model
	order := ToOrder(req)
	
	// 3. 领域验证
	if err := order.Validate(); err != nil {
		return nil, err
	}
	
	// 4. 保存到数据库
	if err := s.repo.Save(ctx, order); err != nil {
		return nil, err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Index(order)
	}
	
	// 6. Domain Model → DTO
	return ToOrderResponse(order), nil
}

// Get 获取订单
func (s *OrderService) Get(ctx context.Context, id uint) (*OrderResponse, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return ToOrderResponse(order), nil
}

// List 订单列表
func (s *OrderService) List(ctx context.Context, limit, offset int) (*OrderListResponse, error) {
	orders, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	
	return ToOrderListResponse(orders, total), nil
}

// Update 更新订单
func (s *OrderService) Update(ctx context.Context, id uint, req *UpdateOrderRequest) error {
	// 1. 查询订单
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 更新字段（通过领域方法）
	if req.Quantity > 0 {
		if err := order.UpdateQuantity(req.Quantity); err != nil {
			return err
		}
	}
	
	if req.UnitPrice >= 0 {
		if err := order.UpdateUnitPrice(req.UnitPrice); err != nil {
			return err
		}
	}
	
	// 状态更新
	if req.Status != "" {
		targetStatus := domain.OrderStatus(req.Status)
		
		switch targetStatus {
		case domain.OrderStatusConfirmed:
			if err := order.Confirm(); err != nil {
				return err
			}
		case domain.OrderStatusProduction:
			if err := order.StartProduction(); err != nil {
				return err
			}
		case domain.OrderStatusCompleted:
			if err := order.Complete(); err != nil {
				return err
			}
		case domain.OrderStatusCancelled:
			if err := order.Cancel(); err != nil {
				return err
			}
		}
	}
	
	// 3. 验证
	if err := order.Validate(); err != nil {
		return err
	}
	
	// 4. 保存
	if err := s.repo.Update(ctx, order); err != nil {
		return err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Update(order)
	}
	
	return nil
}

// Delete 删除订单
func (s *OrderService) Delete(ctx context.Context, id uint) error {
	// 1. 查询订单
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 检查是否可以删除（领域规则）
	if !order.CanDelete() {
		return domain.ErrCannotDeleteCompleted
	}
	
	// 3. 删除
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	// 4. 异步删除 ES 文档
	if s.esSync != nil {
		s.esSync.Delete("orders", strconv.Itoa(int(id)))
	}
	
	return nil
}

// GetByClientID 根据客户 ID 查询订单
func (s *OrderService) GetByClientID(ctx context.Context, clientID uint) ([]*OrderResponse, error) {
	orders, err := s.repo.FindByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	
	responses := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = ToOrderResponse(o)
	}
	
	return responses, nil
}

// GetByProductID 根据产品 ID 查询订单
func (s *OrderService) GetByProductID(ctx context.Context, productID uint) ([]*OrderResponse, error) {
	orders, err := s.repo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	
	responses := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = ToOrderResponse(o)
	}
	
	return responses, nil
}

// GetByStatus 根据状态查询订单
func (s *OrderService) GetByStatus(ctx context.Context, status string, limit, offset int) (*OrderListResponse, error) {
	orderStatus := domain.OrderStatus(status)
	
	orders, err := s.repo.FindByStatus(ctx, orderStatus, limit, offset)
	if err != nil {
		return nil, err
	}
	
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	
	return ToOrderListResponse(orders, total), nil
}