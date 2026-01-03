package domain

import (
	"time"

	"gorm.io/gorm"
)

// ProcessQuote 工艺报价
type ProcessQuote struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	QuoteID string `gorm:"size:50;uniqueIndex;not null" json:"quoteId"`
	ProcessCode string `gorm:"size:50;not null;index" json:"processCode"`
	SupplierCode string `gorm:"size:50;not null;index" json:"supplierCode"`
	QuotePrice float64 `gorm:"type:decimal(10,2);not null" json:"quotePrice"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (ProcessQuote) TableName() string {
	return "process_quotes"
}
