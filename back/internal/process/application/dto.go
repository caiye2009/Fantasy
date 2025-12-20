package application

import "time"

// CreateProcessRequest 创建工序请求
type CreateProcessRequest struct {
	Name         string  `json:"name" binding:"required,min=2,max=100"`
	Description  string  `json:"description" binding:"omitempty"`
	CurrentPrice float64 `json:"currentPrice" binding:"omitempty,min=0"` // 当前价格
}

// UpdateProcessRequest 更新工序请求
type UpdateProcessRequest struct {
	Name         string  `json:"name" binding:"omitempty,min=2,max=100"`
	Description  string  `json:"description" binding:"omitempty"`
	CurrentPrice float64 `json:"currentPrice" binding:"omitempty,min=0"`
}

// ProcessResponse 工序响应
type ProcessResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CurrentPrice float64   `json:"currentPrice"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ProcessListResponse 工序列表响应
type ProcessListResponse struct {
	Total     int64              `json:"total"`
	Processes []*ProcessResponse `json:"processes"`
}