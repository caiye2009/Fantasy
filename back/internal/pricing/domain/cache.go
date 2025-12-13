package domain

import (
	"context"
	"time"
)

// PriceData 价格数据
type PriceData struct {
	Price        float64   `json:"price"`
	SupplierID   uint      `json:"supplier_id"`
	SupplierName string    `json:"supplier_name"`
	QuotedAt     time.Time `json:"quoted_at"`
}

// PriceCache 价格缓存接口
type PriceCache interface {
	GetMin(ctx context.Context, targetType string, targetID uint) (*PriceData, error)
	SetMin(ctx context.Context, targetType string, targetID uint, data *PriceData) error
	UpdateMin(ctx context.Context, targetType string, targetID uint, newPrice *PriceData) error

	GetMax(ctx context.Context, targetType string, targetID uint) (*PriceData, error)
	SetMax(ctx context.Context, targetType string, targetID uint, data *PriceData) error
	UpdateMax(ctx context.Context, targetType string, targetID uint, newPrice *PriceData) error
}