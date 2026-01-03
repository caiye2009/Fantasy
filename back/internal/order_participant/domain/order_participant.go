package domain

import (
	"time"

	"gorm.io/gorm"
)

// OrderParticipant 订单参与者
type OrderParticipant struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderCode string `gorm:"size:50;not null;index" json:"orderCode"`
	Username string `gorm:"size:50;not null;index" json:"username"`
	ParticipantRole string `gorm:"size:50;not null" json:"participantRole"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (OrderParticipant) TableName() string {
	return "order_participants"
}
