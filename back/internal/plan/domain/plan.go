package domain

import (
	"time"
)

// PlanStatus 计划状态
type PlanStatus string

const (
	PlanStatusPlanned   PlanStatus = "planned"   // 已计划
	PlanStatusInProgress PlanStatus = "in_progress" // 进行中
	PlanStatusCompleted PlanStatus = "completed" // 已完成
	PlanStatusCancelled PlanStatus = "cancelled" // 已取消
)

// Plan 计划聚合根
type Plan struct {
	ID          uint
	PlanNo      string
	OrderID     uint
	ProductID   uint
	Quantity    float64
	Status      PlanStatus
	ScheduledAt *time.Time
	CompletedAt *time.Time
	CreatedBy   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
		"status":     string(p.Status),
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