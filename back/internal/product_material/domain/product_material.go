package domain

import (
	"time"

	"gorm.io/gorm"
)

// ProductMaterial 产品原料关联
type ProductMaterial struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductCode string `gorm:"size:50;not null;index" json:"productCode"`
	MaterialCode string `gorm:"size:50;not null;index" json:"materialCode"`
	RequiredQty float64 `gorm:"type:decimal(10,2);not null" json:"requiredQty"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (ProductMaterial) TableName() string {
	return "product_materials"
}
