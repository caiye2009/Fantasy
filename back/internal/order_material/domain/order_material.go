package domain

import (
	"time"

	"gorm.io/gorm"
)

// OrderMaterial 订单原料
type OrderMaterial struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderCode string `gorm:"size:50;not null;index" json:"orderCode"`
	MaterialCode string `gorm:"size:50;not null;index" json:"materialCode"`
	RequiredQty float64 `gorm:"type:decimal(10,2);not null" json:"requiredQty"`
	CreatedBy string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (OrderMaterial) TableName() string {
	return "order_materials"
}
