package domain

import (
	"time"

	"gorm.io/gorm"
)

// Client 客户
type Client struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ClientCode    string         `gorm:"size:50;uniqueIndex;not null" json:"clientCode"`
	ClientName    string         `gorm:"size:100;not null" json:"clientName"`
	ClientCountry string         `gorm:"size:50" json:"clientCountry"`
	CreatedBy     uint           `gorm:"not null" json:"createdBy"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Client) TableName() string {
	return "clients"
}
