package application

import "time"

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Code        string `json:"code" binding:"omitempty,min=2,max=50"`
	Description string `json:"description" binding:"omitempty"`
	ParentID    *uint  `json:"parent_id" binding:"omitempty"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Code        string `json:"code" binding:"omitempty,min=2,max=50"`
	Description string `json:"description" binding:"omitempty"`
	ParentID    *uint  `json:"parent_id" binding:"omitempty"`
}

// DepartmentResponse 部门响应
type DepartmentResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ParentID    *uint      `json:"parent_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// DepartmentListResponse 部门列表响应
type DepartmentListResponse struct {
	Total       int64                 `json:"total"`
	Departments []*DepartmentResponse `json:"departments"`
}
