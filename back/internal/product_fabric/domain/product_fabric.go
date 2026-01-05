package domain

import (
	"time"

	"gorm.io/gorm"
)

// ProductFabric 产品面料关联
type ProductFabric struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProductCode string         `gorm:"size:50;not null;index" json:"productCode"`
	FabricCode  string         `gorm:"size:50;not null;index" json:"fabricCode"`
	RequiredQty float64        `gorm:"type:decimal(10,2);not null" json:"requiredQty"`
	CreatedBy   string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (ProductFabric) TableName() string {
	return "product_fabrics"
}