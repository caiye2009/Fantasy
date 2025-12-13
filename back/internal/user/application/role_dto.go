package application

import "time"

// CreateRoleRequest 创建职位请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Level       int    `json:"level" binding:"omitempty,min=0"`
}

// UpdateRoleRequest 更新职位请求
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Code        string `json:"code" binding:"omitempty,min=2,max=50"`
	Description string `json:"description" binding:"omitempty"`
	Level       int    `json:"level" binding:"omitempty,min=0"`
}

// RoleResponse 职位响应
type RoleResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Level       int        `json:"level"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// RoleListResponse 职位列表响应
type RoleListResponse struct {
	Total int64           `json:"total"`
	Roles []*RoleResponse `json:"roles"`
}
