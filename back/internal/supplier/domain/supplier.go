package domain

import (
	"time"

	"gorm.io/gorm"
)

// Supplier 供应商
type Supplier struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	SupplierCode     string         `gorm:"size:50;uniqueIndex;not null" json:"supplierCode"`
	SupplierName     string         `gorm:"size:100;not null" json:"supplierName"`
	SupplierCategory string         `gorm:"size:50" json:"supplierCategory"`
	CreatedBy        uint           `gorm:"not null" json:"createdBy"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Supplier) TableName() string {
	return "suppliers"
}
