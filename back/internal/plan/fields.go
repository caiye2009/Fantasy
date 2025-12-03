package plan

import (
	"strconv"
	"back/pkg/fields"
)

type Plan struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	PlanNo      string  `gorm:"size:50;uniqueIndex;not null" json:"plan_no"`
	OrderID     uint    `gorm:"not null;index" json:"order_id"`
	ProductID   uint    `gorm:"not null;index" json:"product_id"`
	Quantity    float64 `gorm:"type:decimal(10,2);not null" json:"quantity"`
	Status      string  `gorm:"size:20;default:'planned';index" json:"status"`
	ScheduledAt *fields.Time `gorm:"type:timestamp" json:"scheduled_at"`
	CompletedAt *fields.Time `gorm:"type:timestamp" json:"completed_at"`
	CreatedBy   uint    `gorm:"not null;index" json:"created_by"`
	fields.DTFields
}

func (Plan) TableName() string {
	return "plans"
}

// ========== Indexable 接口实现 ==========

func (p *Plan) GetIndexName() string {
	return "plans"
}

func (p *Plan) GetDocumentID() string {
	return strconv.Itoa(int(p.ID))
}

func (p *Plan) ToDocument() map[string]interface{} {
	doc := map[string]interface{}{
		"id":         p.ID,
		"plan_no":    p.PlanNo,
		"order_id":   p.OrderID,
		"product_id": p.ProductID,
		"quantity":   p.Quantity,
		"status":     p.Status,
		"created_by": p.CreatedBy,
		"created_at": p.CreatedAt,
		"updated_at": p.UpdatedAt,
	}

	if p.ScheduledAt != nil {
		doc["scheduled_at"] = p.ScheduledAt
	}
	if p.CompletedAt != nil {
		doc["completed_at"] = p.CompletedAt
	}

	return doc
}

// ========== 其他结构体保持不变 ==========