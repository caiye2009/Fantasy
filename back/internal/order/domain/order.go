package domain

import (
	"fmt"
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

// ES Indexable 接口实现

// GetIndexName 返回 ES 索引名称
func (o *Order) GetIndexName() string {
	return "order"
}

// GetDocumentID 返回文档 ID
func (o *Order) GetDocumentID() string {
	return fmt.Sprintf("%d", o.ID)
}

// ToDocument 转换为 ES 文档
func (o *Order) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":           o.ID,
		"orderCode":    o.OrderCode,
		"clientCode":   o.ClientCode,
		"orderDate":    o.OrderDate,
		"deliveryDate": o.DeliveryDate,
		"orderStatus":  o.OrderStatus,
		"createdBy":    o.CreatedBy,
		"createdAt":    o.CreatedAt,
		"updatedAt":    o.UpdatedAt,
	}
}