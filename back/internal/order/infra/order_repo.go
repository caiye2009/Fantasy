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

// ==================== 参与者管理 ====================

// CreateParticipant 创建参与者
func (r *OrderRepo) CreateParticipant(ctx context.Context, participant *domain.OrderParticipant) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

// GetParticipants 获取订单的所有参与者
func (r *OrderRepo) GetParticipants(ctx context.Context, orderID uint) ([]domain.OrderParticipant, error) {
	var participants []domain.OrderParticipant
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&participants).Error
	return participants, err
}

// GetActiveParticipants 获取订单的当前激活参与者
func (r *OrderRepo) GetActiveParticipants(ctx context.Context, orderID uint) ([]domain.OrderParticipant, error) {
	var participants []domain.OrderParticipant
	err := r.db.WithContext(ctx).
		Where("order_id = ? AND is_active = ?", orderID, true).
		Order("created_at ASC").
		Find(&participants).Error
	return participants, err
}

// DeactivateParticipant 停用参与者
func (r *OrderRepo) DeactivateParticipant(ctx context.Context, participantID uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.OrderParticipant{}).
		Where("id = ?", participantID).
		Update("is_active", false).Error
}

// FindParticipantByRole 根据角色查找参与者
func (r *OrderRepo) FindParticipantByRole(ctx context.Context, orderID uint, role string) (*domain.OrderParticipant, error) {
	var participant domain.OrderParticipant
	err := r.db.WithContext(ctx).
		Where("order_id = ? AND role = ? AND is_active = ?", orderID, role, true).
		First(&participant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrParticipantNotFound
	}
	return &participant, err
}

// ==================== 进度管理 ====================

// CreateProgress 创建进度项
func (r *OrderRepo) CreateProgress(ctx context.Context, progress *domain.OrderProgress) error {
	return r.db.WithContext(ctx).Create(progress).Error
}

// UpdateProgress 更新进度项
func (r *OrderRepo) UpdateProgress(ctx context.Context, progress *domain.OrderProgress) error {
	return r.db.WithContext(ctx).Save(progress).Error
}

// GetProgresses 获取订单的所有进度项
func (r *OrderRepo) GetProgresses(ctx context.Context, orderID uint) ([]domain.OrderProgress, error) {
	var progresses []domain.OrderProgress
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&progresses).Error
	return progresses, err
}

// FindProgressByType 根据类型查找进度项
func (r *OrderRepo) FindProgressByType(ctx context.Context, orderID uint, progressType string) (*domain.OrderProgress, error) {
	var progress domain.OrderProgress
	err := r.db.WithContext(ctx).
		Where("order_id = ? AND type = ?", orderID, progressType).
		First(&progress).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrProgressNotFound
	}
	return &progress, err
}

// ==================== 事件管理 ====================

// CreateEvent 创建事件
func (r *OrderRepo) CreateEvent(ctx context.Context, event *domain.OrderEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetEvents 获取订单的所有事件（时间降序）
func (r *OrderRepo) GetEvents(ctx context.Context, orderID uint) ([]domain.OrderEvent, error) {
	var events []domain.OrderEvent
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

// ==================== 关联查询 ====================

// GetOrderWithDetails 获取订单及其所有关联数据
func (r *OrderRepo) GetOrderWithDetails(ctx context.Context, orderID uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Participants").
		Preload("Progresses").
		Preload("Events", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&order, orderID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrOrderNotFound
	}
	return &order, err
}

// GetOrderWithClientProduct 获取订单详情（包含客户和产品信息）
func (r *OrderRepo) GetOrderWithClientProduct(ctx context.Context, orderID uint) (map[string]interface{}, error) {
	var result map[string]interface{}

	err := r.db.WithContext(ctx).
		Table("orders").
		Select(`
			orders.*,
			clients.name as client_name,
			products.name as product_name,
			products.code as product_code
		`).
		Joins("LEFT JOIN clients ON orders.client_id = clients.id").
		Joins("LEFT JOIN products ON orders.product_id = products.id").
		Where("orders.id = ?", orderID).
		First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrOrderNotFound
	}
	return result, err
}

// GetOrdersWithClientProduct 获取订单列表（包含客户和产品信息）
func (r *OrderRepo) GetOrdersWithClientProduct(ctx context.Context, limit, offset int) ([]map[string]interface{}, int64, error) {
	var results []map[string]interface{}
	var total int64

	// 获取总数
	r.db.WithContext(ctx).Table("orders").Count(&total)

	// 获取数据
	err := r.db.WithContext(ctx).
		Table("orders").
		Select(`
			orders.id,
			orders.order_no,
			orders.client_id,
			clients.name as client_name,
			orders.product_id,
			products.name as product_name,
			products.code as product_code,
			orders.required_quantity,
			orders.product_history_shrinkage,
			orders.unit_price,
			orders.total_price,
			orders.status,
			orders.assigned_department,
			orders.created_at,
			orders.updated_at
		`).
		Joins("LEFT JOIN clients ON orders.client_id = clients.id").
		Joins("LEFT JOIN products ON orders.product_id = products.id").
		Order("orders.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&results).Error

	return results, total, err
}

// GetProgressesByOrderIDs 批量获取订单进度
func (r *OrderRepo) GetProgressesByOrderIDs(ctx context.Context, orderIDs []uint) ([]domain.OrderProgress, error) {
	var progresses []domain.OrderProgress
	err := r.db.WithContext(ctx).
		Where("order_id IN ?", orderIDs).
		Find(&progresses).Error
	return progresses, err
}

// GetEventsByOrderIDs 批量获取订单事件
func (r *OrderRepo) GetEventsByOrderIDs(ctx context.Context, orderIDs []uint) ([]domain.OrderEvent, error) {
	var events []domain.OrderEvent
	err := r.db.WithContext(ctx).
		Where("order_id IN ?", orderIDs).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

// Transaction 执行事务
func (r *OrderRepo) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建新的上下文，包含事务
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}