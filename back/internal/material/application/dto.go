package application

import "time"

// CreateMaterialRequest 创建材料请求
type CreateMaterialRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Spec        string `json:"spec" binding:"omitempty,max=200"`
	Unit        string `json:"unit" binding:"omitempty,max=20"`
	Description string `json:"description" binding:"omitempty"`
}

// UpdateMaterialRequest 更新材料请求
type UpdateMaterialRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Spec        string `json:"spec" binding:"omitempty,max=200"`
	Unit        string `json:"unit" binding:"omitempty,max=20"`
	Description string `json:"description" binding:"omitempty"`
}

// MaterialResponse 材料响应
type MaterialResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Spec        string    `json:"spec"`
	Unit        string    `json:"unit"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MaterialListResponse 材料列表响应
type MaterialListResponse struct {
	Total     int64               `json:"total"`
	Materials []*MaterialResponse `json:"materials"`
}