package application

import "time"

// CreateInventoryRequest 创建库存请求
type CreateInventoryRequest struct {
	ProductID uint    `json:"productId" binding:"required"`
	Category  string  `json:"category" binding:"required"`
	BatchID   string  `json:"batchId" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	Unit      string  `json:"unit" binding:"required"`
	UnitCost  float64 `json:"unitCost" binding:"required,gte=0"`
	Remark    string  `json:"remark"`
}

// UpdateInventoryRequest 更新库存请求
type UpdateInventoryRequest struct {
	ID        uint    `json:"id" binding:"required"`
	ProductID uint    `json:"productId" binding:"required"`
	Category  string  `json:"category" binding:"required"`
	BatchID   string  `json:"batchId" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	Unit      string  `json:"unit" binding:"required"`
	UnitCost  float64 `json:"unitCost" binding:"required,gte=0"`
	Remark    string  `json:"remark"`
}

// UpdateQuantityRequest 更新数量请求
type UpdateQuantityRequest struct {
	ID       uint    `json:"id" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

// DeductInventoryRequest 扣减库存请求
type DeductInventoryRequest struct {
	ID       uint    `json:"id" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

// AddInventoryRequest 增加库存请求
type AddInventoryRequest struct {
	ID       uint    `json:"id" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

// InventoryResponse 库存响应
type InventoryResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"productId"`
	Category  string    `json:"category"`
	BatchID   string    `json:"batchId"`
	Quantity  float64   `json:"quantity"`
	Unit      string    `json:"unit"`
	UnitCost  float64   `json:"unitCost"`
	TotalCost float64   `json:"totalCost"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// InventoryListResponse 库存列表响应
type InventoryListResponse struct {
	Total       int                  `json:"total"`
	Inventories []*InventoryResponse `json:"inventories"`
}
