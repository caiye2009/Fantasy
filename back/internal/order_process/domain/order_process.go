package domain

import (
	"time"

	"gorm.io/gorm"
)

// OrderProcess 订单工艺
type OrderProcess struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderCode string `gorm:"size:50;not null;index" json:"orderCode"`
	ProcessCode string `gorm:"size:50;not null;index" json:"processCode"`
	ProcessSeq int `gorm:"not null" json:"processSeq"`
	CreatedBy string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (OrderProcess) TableName() string {
	return "order_processes"
}
