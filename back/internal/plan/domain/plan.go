package domain

import (
	"time"

	"gorm.io/gorm"
)

// 计划状态常量
const (
	PlanStatusPlanned    = "planned"     // 已计划
	PlanStatusInProgress = "in_progress" // 进行中
	PlanStatusCompleted  = "completed"   // 已完成
	PlanStatusCancelled  = "cancelled"   // 已取消
)

// Plan 计划聚合根
type Plan struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PlanNo      string         `gorm:"size:50;uniqueIndex;not null" json:"plan_no"`
	OrderID     uint           `gorm:"not null;index" json:"order_id"`
	ProductID   uint           `gorm:"not null;index" json:"product_id"`
	Quantity    float64        `gorm:"type:decimal(10,2);not null" json:"quantity"`
	Status      string         `gorm:"size:20;default:planned;index" json:"status"`
	ScheduledAt *time.Time     `gorm:"type:timestamp" json:"scheduled_at,omitempty"`
	CompletedAt *time.Time     `gorm:"type:timestamp" json:"completed_at,omitempty"`
	CreatedBy   uint           `gorm:"not null;index" json:"created_by"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Plan) TableName() string {
	return "plans"
}

// Validate 验证计划数据
func (p *Plan) Validate() error {
	if p.PlanNo == "" {
		return ErrPlanNoEmpty
	}
	
	if len(p.PlanNo) > 50 {
		return ErrPlanNoTooLong
	}
	
	if p.OrderID == 0 {
		return ErrOrderIDRequired
	}
	
	if p.ProductID == 0 {
		return ErrProductIDRequired
	}
	
	if p.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	
	if p.CreatedBy == 0 {
		return ErrCreatedByRequired
	}
	
	return nil
}

// Start 开始计划
func (p *Plan) Start() error {
	if p.Status != PlanStatusPlanned {
		return ErrCannotStartPlan
	}
	
	p.Status = PlanStatusInProgress
	return nil
}

// Complete 完成计划
func (p *Plan) Complete() error {
	if p.Status != PlanStatusInProgress {
		return ErrCannotCompletePlan
	}
	
	now := time.Now()
	p.Status = PlanStatusCompleted
	p.CompletedAt = &now
	return nil
}

// Cancel 取消计划
func (p *Plan) Cancel() error {
	if p.Status == PlanStatusCompleted {
		return ErrCannotCancelCompletedPlan
	}
	
	p.Status = PlanStatusCancelled
	return nil
}

// UpdateQuantity 更新数量
func (p *Plan) UpdateQuantity(newQuantity float64) error {
	if newQuantity <= 0 {
		return ErrInvalidQuantity
	}
	
	if p.Status == PlanStatusCompleted {
		return ErrCannotUpdateCompletedPlan
	}
	
	p.Quantity = newQuantity
	return nil
}

// IsCompleted 是否已完成
func (p *Plan) IsCompleted() bool {
	return p.Status == PlanStatusCompleted
}

// CanDelete 是否可以删除
func (p *Plan) CanDelete() bool {
	// 已完成的计划不允许删除
	return p.Status != PlanStatusCompleted
}

// ToDocument 转换为 ES 文档
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

// GetIndexName ES 索引名称
func (p *Plan) GetIndexName() string {
	return "plans"
}

// GetDocumentID ES 文档 ID
func (p *Plan) GetDocumentID() string {
	return string(rune(p.ID))
}