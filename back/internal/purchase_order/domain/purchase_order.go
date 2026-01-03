package domain

import (
	"time"

	"gorm.io/gorm"
)

// PurchaseOrder 采购单
type PurchaseOrder struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PurchaseCode string `gorm:"size:50;uniqueIndex;not null" json:"purchaseCode"`
	OrderCode string `gorm:"size:50;not null;index" json:"orderCode"`
	SupplierCode string `gorm:"size:50;not null;index" json:"supplierCode"`
	PurchaseDate time.Time `gorm:"type:date;not null" json:"purchaseDate"`
	PurchaseStatus string `gorm:"size:50;not null;default:'pending'" json:"purchaseStatus"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}
