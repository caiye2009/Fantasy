package application

import "time"

// CreateProcessRequest 创建工序请求
type CreateProcessRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty"`
}

// UpdateProcessRequest 更新工序请求
type UpdateProcessRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"omitempty"`
}

// ProcessResponse 工序响应
type ProcessResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProcessListResponse 工序列表响应
type ProcessListResponse struct {
	Total     int64              `json:"total"`
	Processes []*ProcessResponse `json:"processes"`
}