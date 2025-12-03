package price

import (
	"time"
	"back/pkg/fields"
)

type ProcessPrice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VendorID  uint      `gorm:"not null;index:idx_vendor_process" json:"vendor_id"`
	ProcessID uint      `gorm:"not null;index:idx_vendor_process;index:idx_process_time" json:"process_id"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	QuotedAt  time.Time `gorm:"not null;index:idx_process_time" json:"quoted_at"`
	fields.DTFields
}


type PriceData struct {
	ID         uint      `json:"id"`
	Price      float64   `json:"price"`
	VendorID   uint      `json:"vendor_id"`
	VendorName string    `json:"vendor_name"`
	QuotedAt   time.Time `json:"quoted_at"`
}