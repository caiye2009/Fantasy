package application

import "time"

// CreatePlanRequest 创建计划请求
type CreatePlanRequest struct {
	PlanNo      string     `json:"plan_no" binding:"required,max=50"`
	OrderID     uint       `json:"order_id" binding:"required"`
	ProductID   uint       `json:"product_id" binding:"required"`
	Quantity    float64    `json:"quantity" binding:"required,gt=0"`
	ScheduledAt *time.Time `json:"scheduled_at" binding:"omitempty"`
	CreatedBy   uint       `json:"created_by" binding:"required"`
}

// UpdatePlanRequest 更新计划请求
type UpdatePlanRequest struct {
	Status      string     `json:"status" binding:"omitempty,oneof=planned in_progress completed cancelled"`
	Quantity    float64    `json:"quantity" binding:"omitempty,gt=0"`
	CompletedAt *time.Time `json:"completed_at" binding:"omitempty"`
}

// PlanResponse 计划响应
type PlanResponse struct {
	ID          uint       `json:"id"`
	PlanNo      string     `json:"plan_no"`
	OrderID     uint       `json:"order_id"`
	ProductID   uint       `json:"product_id"`
	Quantity    float64    `json:"quantity"`
	Status      string     `json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedBy   uint       `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// PlanListResponse 计划列表响应
type PlanListResponse struct {
	Total int64           `json:"total"`
	Plans []*PlanResponse `json:"plans"`
}