package price

import (
	"time"
	"back/pkg/fields"
)

type MaterialPrice struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VendorID   uint      `gorm:"not null;index:idx_vendor_material" json:"vendor_id"`
	MaterialID uint      `gorm:"not null;index:idx_vendor_material;index:idx_material_time" json:"material_id"`
	Price      float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	QuotedAt   time.Time `gorm:"not null;index:idx_material_time" json:"quoted_at"`
	fields.DTFields
}

type PriceData struct {
	ID         uint      `json:"id"`
	Price      float64   `json:"price"`
	VendorID   uint      `json:"vendor_id"`
	VendorName string    `json:"vendor_name"`
	QuotedAt   time.Time `json:"quoted_at"`
}