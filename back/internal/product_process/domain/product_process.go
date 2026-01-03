package domain

import (
	"time"

	"gorm.io/gorm"
)

// ProductProcess 产品工艺关联
type ProductProcess struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductCode string `gorm:"size:50;not null;index" json:"productCode"`
	ProcessCode string `gorm:"size:50;not null;index" json:"processCode"`
	ProcessSeq int `gorm:"not null" json:"processSeq"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (ProductProcess) TableName() string {
	return "product_processes"
}
