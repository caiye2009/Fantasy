package application

import "time"

// CreateSupplierRequest 创建供应商请求
type CreateSupplierRequest struct {
	SupplierCode     string `json:"supplierCode" binding:"required,max=50"`
	SupplierName     string `json:"supplierName" binding:"required,max=100"`
	SupplierCategory string `json:"supplierCategory" binding:"max=50"`
}

// SupplierResponse 供应商响应
type SupplierResponse struct {
	ID               uint      `json:"id"`
	SupplierCode     string    `json:"supplierCode"`
	SupplierName     string    `json:"supplierName"`
	SupplierCategory string    `json:"supplierCategory"`
	CreatedBy        string    `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
