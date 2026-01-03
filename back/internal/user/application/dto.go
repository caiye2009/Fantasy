package application

import "time"

// UserCreateRequest 创建用户请求（仅暴露给前端的字段）
type UserCreateRequest struct {
	Username   string `json:"username" validate:"required"`
	Role       string `json:"role" validate:"required"`
	Department string `json:"department" validate:"required"`
}

// UserCreateResponse 创建用户响应
type UserCreateResponse struct {
	LoginID  string `json:"loginId"`
	Username string `json:"username"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID          uint      `json:"id"`
	LoginID     string    `json:"loginId"`
	Username    string    `json:"username"`
	Password    string    `json:"password,omitempty"`
	Role        string    `json:"role"`
	Department  string    `json:"department"`
	IsActive    bool      `json:"isActive"`
	HasInitPass bool      `json:"hasInitPass"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}