package domain

import (
	"time"

	"gorm.io/gorm"
)

// OrderProduct 订单产品
type OrderProduct struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderCode string `gorm:"size:50;not null;index" json:"orderCode"`
	ProductCode string `gorm:"size:50;not null;index" json:"productCode"`
	OrderedQty float64 `gorm:"type:decimal(10,2);not null" json:"orderedQty"`
	CreatedBy string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (OrderProduct) TableName() string {
	return "order_products"
}
