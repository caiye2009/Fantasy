package domain

import (
	"time"

	"gorm.io/gorm"
)

// Product 产品
type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductCode string `gorm:"size:50;uniqueIndex;not null" json:"productCode"`
	ProductName string `gorm:"size:100;not null" json:"productName"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Product) TableName() string {
	return "products"
}
