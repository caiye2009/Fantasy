package order

import (
	"strconv"
	"back/pkg/fields"
)

type Order struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	OrderNo    string  `gorm:"size:50;uniqueIndex;not null" json:"order_no"`
	ClientID   uint    `gorm:"not null;index" json:"client_id"`
	ProductID  uint    `gorm:"not null;index" json:"product_id"`
	Quantity   float64 `gorm:"type:decimal(10,2);not null" json:"quantity"`
	UnitPrice  float64 `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     string  `gorm:"size:20;default:'pending';index" json:"status"`
	CreatedBy  uint    `gorm:"not null;index" json:"created_by"`
	fields.DTFields
}

func (Order) TableName() string {
	return "orders"
}

// ========== Indexable 接口实现 ==========

func (o *Order) GetIndexName() string {
	return "orders"
}

func (o *Order) GetDocumentID() string {
	return strconv.Itoa(int(o.ID))
}

func (o *Order) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":          o.ID,
		"order_no":    o.OrderNo,
		"client_id":   o.ClientID,
		"product_id":  o.ProductID,
		"quantity":    o.Quantity,
		"unit_price":  o.UnitPrice,
		"total_price": o.TotalPrice,
		"status":      o.Status,
		"created_by":  o.CreatedBy,
		"created_at":  o.CreatedAt,
		"updated_at":  o.UpdatedAt,
	}
}

// ========== 其他结构体保持不变 ==========