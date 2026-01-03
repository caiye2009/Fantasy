package domain

import (
	"time"

	"gorm.io/gorm"
)

// MaterialQuote 原料报价
type MaterialQuote struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	QuoteID string `gorm:"size:50;uniqueIndex;not null" json:"quoteId"`
	MaterialCode string `gorm:"size:50;not null;index" json:"materialCode"`
	SupplierCode string `gorm:"size:50;not null;index" json:"supplierCode"`
	QuotePrice float64 `gorm:"type:decimal(10,2);not null" json:"quotePrice"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (MaterialQuote) TableName() string {
	return "material_quotes"
}
