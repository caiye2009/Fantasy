package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"back/internal/order/domain"
	"back/pkg/repo"
)

// OrderRepo 订单仓储实现
type OrderRepo struct {
	*repo.Repo[domain.Order]
	db *gorm.DB
}

// NewOrderRepo 创建仓储
func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{
		Repo: repo.NewRepo[domain.Order](db),
		db:   db,
	}
}

// Save 保存订单
func (r *OrderRepo) Save(ctx context.Context, order *domain.Order) error {
	return r.Create(ctx, order)
}

// FindByID 根据 ID 查询
func (r *OrderRepo) FindByID(ctx context.Context, id uint) (*domain.Order, error) {
	order, err := r.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return order, nil
}

// FindAll 查询所有
func (r *OrderRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	orders, err := r.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Order, len(orders))
	for i := range orders {
		result[i] = &orders[i]
	}
	return result, nil
}

// Update 更新订单
func (r *OrderRepo) Update(ctx context.Context, order *domain.Order) error {
	return r.Repo.Update(ctx, order)
}

// Delete 删除订单
func (r *OrderRepo) Delete(ctx context.Context, id uint) error {
	return r.Repo.Delete(ctx, id)
}

// ExistsByID 检查是否存在
func (r *OrderRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"id": id})
}

// FindByOrderNo 根据订单编号查询
func (r *OrderRepo) FindByOrderNo(ctx context.Context, orderNo string) (*domain.Order, error) {
	order, err := r.First(ctx, map[string]interface{}{"order_no": orderNo})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return order, nil
}

// ExistsByOrderNo 检查订单编号是否存在
func (r *OrderRepo) ExistsByOrderNo(ctx context.Context, orderNo string) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"order_no": orderNo})
}

// FindByClientID 根据客户 ID 查询
func (r *OrderRepo) FindByClientID(ctx context.Context, clientID uint) ([]*domain.Order, error) {
	var result []*domain.Order
	err := r.db.WithContext(ctx).
		Where("client_id = ?", clientID).
		Order("created_at DESC").
		Find(&result).Error
	return result, err
}

// FindByProductID 根据产品 ID 查询
func (r *OrderRepo) FindByProductID(ctx context.Context, productID uint) ([]*domain.Order, error) {
	var result []*domain.Order
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&result).Error
	return result, err
}

// FindByStatus 根据状态查询
func (r *OrderRepo) FindByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.Order, error) {
	var result []*domain.Order
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&result).Error
	return result, err
}

// Count 统计数量
func (r *OrderRepo) Count(ctx context.Context) (int64, error) {
	return r.Repo.Count(ctx, map[string]interface{}{})
}