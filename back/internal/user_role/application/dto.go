package application

import "time"

// CreateUserRoleRequest 创建角色请求
type CreateUserRoleRequest struct {
	RoleCode    string `json:"roleCode" validate:"required"`
	RoleName    string `json:"roleName" validate:"required"`
	Description string `json:"description"`
	CreatedBy   uint   `json:"createdBy" validate:"required"`
}

// UserRoleResponse 角色响应
type UserRoleResponse struct {
	ID          uint      `json:"id"`
	RoleCode    string    `json:"roleCode"`
	RoleName    string    `json:"roleName"`
	Description string    `json:"description"`
	CreatedBy   uint      `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
