package domain

import (
	"time"

	"gorm.io/gorm"
)

// Material 原料
type Material struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MaterialCode string `gorm:"size:50;uniqueIndex;not null" json:"materialCode"`
	MaterialName string `gorm:"size:100;not null" json:"materialName"`
	MaterialCountry string `gorm:"size:50" json:"materialCountry"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Material) TableName() string {
	return "materials"
}
