package domain

import (
	"time"

	"gorm.io/gorm"
)

// OrderEvent 订单事件
type OrderEvent struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	EventID string `gorm:"size:50;uniqueIndex;not null" json:"eventId"`
	OrderCode string `gorm:"size:50;not null;index" json:"orderCode"`
	EventType string `gorm:"size:50;not null" json:"eventType"`
	EventContent string `gorm:"type:text" json:"eventContent"`
	Operator string `gorm:"size:100;not null" json:"operator"`
	OperatedAt time.Time `json:"operatedAt"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (OrderEvent) TableName() string {
	return "order_events"
}
