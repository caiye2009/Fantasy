package application

import "time"

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username   string `json:"username" binding:"required,min=2,max=50"`
	Department string `json:"department" binding:"required,max=100"`
	Role       string `json:"role" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username   string `json:"username" binding:"omitempty,min=2,max=50"`
	Department string `json:"department" binding:"omitempty,max=100"`
	Email      string `json:"email" binding:"omitempty,email"`
	Role       string `json:"role" binding:"omitempty"`
	Status     string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required_without=IsInitPass"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID          uint      `json:"id"`
	LoginID     string    `json:"login_id"`
	Username    string    `json:"username"`
	Department  string    `json:"department"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	HasInitPass bool      `json:"has_init_pass"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64           `json:"total"`
	Users []*UserResponse `json:"users"`
}

// CreateUserResponse 创建用户响应（只返回 login_id）
type CreateUserResponse struct {
	LoginID string `json:"login_id"`
}