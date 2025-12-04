package application

import "time"

// CreateClientRequest 创建客户请求
type CreateClientRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Contact string `json:"contact" binding:"omitempty,max=50"`
	Phone   string `json:"phone" binding:"omitempty,max=20"`
	Email   string `json:"email" binding:"omitempty,email"`
	Address string `json:"address" binding:"omitempty,max=200"`
}

// UpdateClientRequest 更新客户请求
type UpdateClientRequest struct {
	Name    string `json:"name" binding:"omitempty,min=2,max=100"`
	Contact string `json:"contact" binding:"omitempty,max=50"`
	Phone   string `json:"phone" binding:"omitempty,max=20"`
	Email   string `json:"email" binding:"omitempty,email"`
	Address string `json:"address" binding:"omitempty,max=200"`
}

// ClientResponse 客户响应
type ClientResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Contact   string    `json:"contact"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ClientListResponse 客户列表响应
type ClientListResponse struct {
	Total   int64             `json:"total"`
	Clients []*ClientResponse `json:"clients"`
}