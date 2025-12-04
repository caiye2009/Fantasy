package infra

import (
	"time"
	
	"gorm.io/gorm"
	
	"back/internal/order/domain"
)

// OrderPO 订单持久化对象
type OrderPO struct {
	ID         uint           `gorm:"primaryKey"`
	OrderNo    string         `gorm:"size:50;uniqueIndex;not null"`
	ClientID   uint           `gorm:"not null;index"`
	ProductID  uint           `gorm:"not null;index"`
	Quantity   float64        `gorm:"type:decimal(10,2);not null"`
	UnitPrice  float64        `gorm:"type:decimal(10,2);not null"`
	TotalPrice float64        `gorm:"type:decimal(10,2);not null"`
	Status     string         `gorm:"size:20;default:'pending';index"`
	CreatedBy  uint           `gorm:"not null;index"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// TableName 表名
func (OrderPO) TableName() string {
	return "orders"
}

// ToDomain 转换为领域模型
func (po *OrderPO) ToDomain() *domain.Order {
	return &domain.Order{
		ID:         po.ID,
		OrderNo:    po.OrderNo,
		ClientID:   po.ClientID,
		ProductID:  po.ProductID,
		Quantity:   po.Quantity,
		UnitPrice:  po.UnitPrice,
		TotalPrice: po.TotalPrice,
		Status:     domain.OrderStatus(po.Status),
		CreatedBy:  po.CreatedBy,
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
	}
}

// FromDomain 从领域模型转换
func FromDomain(o *domain.Order) *OrderPO {
	return &OrderPO{
		ID:         o.ID,
		OrderNo:    o.OrderNo,
		ClientID:   o.ClientID,
		ProductID:  o.ProductID,
		Quantity:   o.Quantity,
		UnitPrice:  o.UnitPrice,
		TotalPrice: o.TotalPrice,
		Status:     string(o.Status),
		CreatedBy:  o.CreatedBy,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
}