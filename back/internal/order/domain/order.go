package domain

import (
	"time"
	"gorm.io/gorm"
)

// 订单状态常量
const (
	OrderStatusPending    = "pending"     // 待处理
	OrderStatusConfirmed  = "confirmed"   // 已确认
	OrderStatusProduction = "production"  // 生产中
	OrderStatusCompleted  = "completed"   // 已完成
	OrderStatusCancelled  = "cancelled"   // 已取消
)

// Order 订单聚合根
type Order struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	OrderNo    string         `gorm:"size:50;uniqueIndex;not null" json:"order_no"`
	ClientID   uint           `gorm:"not null;index" json:"client_id"`
	ProductID  uint           `gorm:"not null;index" json:"product_id"`
	Quantity   float64        `gorm:"type:decimal(10,2);not null" json:"quantity"`
	UnitPrice  float64        `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalPrice float64        `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     string         `gorm:"size:20;default:pending;index" json:"status"`
	CreatedBy  uint           `gorm:"not null;index" json:"created_by"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Order) TableName() string {
	return "orders"
}

// Validate 验证订单数据
func (o *Order) Validate() error {
	if o.OrderNo == "" {
		return ErrOrderNoEmpty
	}
	
	if len(o.OrderNo) > 50 {
		return ErrOrderNoTooLong
	}
	
	if o.ClientID == 0 {
		return ErrClientIDRequired
	}
	
	if o.ProductID == 0 {
		return ErrProductIDRequired
	}
	
	if o.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	
	if o.UnitPrice < 0 {
		return ErrInvalidUnitPrice
	}
	
	if o.CreatedBy == 0 {
		return ErrCreatedByRequired
	}
	
	return nil
}

// CalculateTotalPrice 计算总价
func (o *Order) CalculateTotalPrice() {
	o.TotalPrice = o.Quantity * o.UnitPrice
}

// Confirm 确认订单
func (o *Order) Confirm() error {
	if o.Status != OrderStatusPending {
		return ErrCannotConfirm
	}
	
	o.Status = OrderStatusConfirmed
	return nil
}

// StartProduction 开始生产
func (o *Order) StartProduction() error {
	if o.Status != OrderStatusConfirmed {
		return ErrCannotStartProduction
	}
	
	o.Status = OrderStatusProduction
	return nil
}

// Complete 完成订单
func (o *Order) Complete() error {
	if o.Status != OrderStatusProduction {
		return ErrCannotComplete
	}
	
	o.Status = OrderStatusCompleted
	return nil
}

// Cancel 取消订单
func (o *Order) Cancel() error {
	if o.Status == OrderStatusCompleted {
		return ErrCannotCancelCompleted
	}
	
	o.Status = OrderStatusCancelled
	return nil
}

// UpdateQuantity 更新数量
func (o *Order) UpdateQuantity(newQuantity float64) error {
	if newQuantity <= 0 {
		return ErrInvalidQuantity
	}
	
	if o.Status == OrderStatusCompleted || o.Status == OrderStatusCancelled {
		return ErrCannotUpdateCompleted
	}
	
	o.Quantity = newQuantity
	o.CalculateTotalPrice()
	return nil
}

// UpdateUnitPrice 更新单价
func (o *Order) UpdateUnitPrice(newPrice float64) error {
	if newPrice < 0 {
		return ErrInvalidUnitPrice
	}
	
	if o.Status == OrderStatusCompleted || o.Status == OrderStatusCancelled {
		return ErrCannotUpdateCompleted
	}
	
	o.UnitPrice = newPrice
	o.CalculateTotalPrice()
	return nil
}

// CanDelete 是否可以删除
func (o *Order) CanDelete() bool {
	return o.Status != OrderStatusCompleted
}

// IsCompleted 是否已完成
func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusCompleted
}

// ToDocument 转换为 ES 文档
func (o *Order) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":          o.ID,
		"order_no":    o.OrderNo,
		"client_id":   o.ClientID,
		"product_id":  o.ProductID,
		"quantity":    o.Quantity,
		"unit_price":  o.UnitPrice,
		"total_price": o.TotalPrice,
		"status":      o.Status,
		"created_by":  o.CreatedBy,
		"created_at":  o.CreatedAt,
		"updated_at":  o.UpdatedAt,
	}
}

// GetIndexName ES 索引名称
func (o *Order) GetIndexName() string {
	return "orders"
}

// GetDocumentID ES 文档 ID
func (o *Order) GetDocumentID() string {
	return string(rune(o.ID))
} 