package domain

import (
	"time"

	"gorm.io/gorm"
)

// ProductionOrder 生产指示
type ProductionOrder struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductionCode string `gorm:"size:50;uniqueIndex;not null" json:"productionCode"`
	OrderCode string `gorm:"size:50;not null;index" json:"orderCode"`
	ProcessCode string `gorm:"size:50;not null;index" json:"processCode"`
	SupplierCode string `gorm:"size:50;not null;index" json:"supplierCode"`
	ProductionQty float64 `gorm:"type:decimal(10,2);not null" json:"productionQty"`
	ProductionStatus string `gorm:"size:50;not null;default:'pending'" json:"productionStatus"`
	CreatedBy string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (ProductionOrder) TableName() string {
	return "production_orders"
}
