package application

import "time"

// CreateUserDeptRequest 创建部门请求
type CreateUserDeptRequest struct {
	DeptCode    string `json:"deptCode" validate:"required"`
	DeptName    string `json:"deptName" validate:"required"`
	Description string `json:"description"`
	CreatedBy   uint   `json:"createdBy" validate:"required"`
}

// UserDeptResponse 部门响应
type UserDeptResponse struct {
	ID          uint      `json:"id"`
	DeptCode    string    `json:"deptCode"`
	DeptName    string    `json:"deptName"`
	Description string    `json:"description"`
	CreatedBy   uint      `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
