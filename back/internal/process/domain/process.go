package domain

import (
	"time"

	"gorm.io/gorm"
)

// Process 工艺
type Process struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProcessCode string `gorm:"size:50;uniqueIndex;not null" json:"processCode"`
	ProcessName string `gorm:"size:100;not null" json:"processName"`
	ProcessCountry string `gorm:"size:50" json:"processCountry"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Process) TableName() string {
	return "processes"
}
