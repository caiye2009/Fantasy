package application

import (
	"context"
	"time"
)

// QuoteRequest 报价请求
type QuoteRequest struct {
	TargetID   uint    `json:"target_id" binding:"required"`
	SupplierID uint    `json:"supplier_id" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
}

// PriceData 价格数据
type PriceData struct {
	Price        float64   `json:"price"`
	SupplierID   uint      `json:"supplier_id"`
	SupplierName string    `json:"supplier_name"`
	QuotedAt     time.Time `json:"quoted_at"`
}

// SupplierServiceInterface Supplier 服务接口
type SupplierServiceInterface interface {
	GetSupplierInfo(ctx context.Context, id uint) (*SupplierInfo, error)
}

// SupplierInfo 供应商基本信息
type SupplierInfo struct {
	ID   uint
	Name string
}