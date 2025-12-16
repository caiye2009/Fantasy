package application

// QuoteRequest 报价请求
type QuoteRequest struct {
	TargetID   uint    `json:"target_id" binding:"required"`
	SupplierID uint    `json:"supplier_id" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
}