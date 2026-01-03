package domain

import (
	"time"

	"gorm.io/gorm"
)

// Order 订单
type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderCode string `gorm:"size:50;uniqueIndex;not null" json:"orderCode"`
	ClientCode string `gorm:"size:50;not null;index" json:"clientCode"`
	OrderDate time.Time `gorm:"type:date;not null" json:"orderDate"`
	DeliveryDate time.Time `gorm:"type:date" json:"deliveryDate"`
	OrderStatus string `gorm:"size:50;not null;default:'pending'" json:"orderStatus"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Order) TableName() string {
	return "orders"
}
