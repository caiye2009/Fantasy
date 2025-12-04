package infra

import (
	"context"
	"errors"
	
	"gorm.io/gorm"
	
	"back/internal/order/domain"
)

// OrderRepoImpl 订单仓储实现
type OrderRepoImpl struct {
	db *gorm.DB
}

// NewOrderRepoImpl 创建仓储实现
func NewOrderRepoImpl(db *gorm.DB) domain.OrderRepository {
	return &OrderRepoImpl{db: db}
}

// Save 保存订单
func (r *OrderRepoImpl) Save(ctx context.Context, order *domain.Order) error {
	po := FromDomain(order)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return err
	}
	order.ID = po.ID
	order.CreatedAt = po.CreatedAt
	order.UpdatedAt = po.UpdatedAt
	return nil
}

// FindByID 根据 ID 查询
func (r *OrderRepoImpl) FindByID(ctx context.Context, id uint) (*domain.Order, error) {
	var po OrderPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindAll 查询所有
func (r *OrderRepoImpl) FindAll(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	var pos []OrderPO
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	orders := make([]*domain.Order, len(pos))
	for i, po := range pos {
		orders[i] = po.ToDomain()
	}
	return orders, nil
}

// Update 更新订单
func (r *OrderRepoImpl) Update(ctx context.Context, order *domain.Order) error {
	po := FromDomain(order)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除订单
func (r *OrderRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&OrderPO{}, id).Error
}

// ExistsByID 检查是否存在
func (r *OrderRepoImpl) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrderPO{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByOrderNo 根据订单编号查询
func (r *OrderRepoImpl) FindByOrderNo(ctx context.Context, orderNo string) (*domain.Order, error) {
	var po OrderPO
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ExistsByOrderNo 检查订单编号是否存在
func (r *OrderRepoImpl) ExistsByOrderNo(ctx context.Context, orderNo string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrderPO{}).Where("order_no = ?", orderNo).Count(&count).Error
	return count > 0, err
}

// FindByClientID 根据客户 ID 查询
func (r *OrderRepoImpl) FindByClientID(ctx context.Context, clientID uint) ([]*domain.Order, error) {
	var pos []OrderPO
	err := r.db.WithContext(ctx).Where("client_id = ?", clientID).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	orders := make([]*domain.Order, len(pos))
	for i, po := range pos {
		orders[i] = po.ToDomain()
	}
	return orders, nil
}

// FindByProductID 根据产品 ID 查询
func (r *OrderRepoImpl) FindByProductID(ctx context.Context, productID uint) ([]*domain.Order, error) {
	var pos []OrderPO
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	orders := make([]*domain.Order, len(pos))
	for i, po := range pos {
		orders[i] = po.ToDomain()
	}
	return orders, nil
}

// FindByStatus 根据状态查询
func (r *OrderRepoImpl) FindByStatus(ctx context.Context, status domain.OrderStatus, limit, offset int) ([]*domain.Order, error) {
	var pos []OrderPO
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	
	orders := make([]*domain.Order, len(pos))
	for i, po := range pos {
		orders[i] = po.ToDomain()
	}
	return orders, nil
}

// Count 统计数量
func (r *OrderRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrderPO{}).Count(&count).Error
	return count, err
}