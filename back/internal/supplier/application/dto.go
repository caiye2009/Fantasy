package application

import "time"

// CreateSupplierRequest 创建供应商请求
type CreateSupplierRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Contact string `json:"contact" binding:"omitempty,max=50"`
	Phone   string `json:"phone" binding:"omitempty,max=20"`
	Email   string `json:"email" binding:"omitempty,email"`
	Address string `json:"address" binding:"omitempty,max=200"`
}

// UpdateSupplierRequest 更新供应商请求
type UpdateSupplierRequest struct {
	Name    string `json:"name" binding:"omitempty,min=2,max=100"`
	Contact string `json:"contact" binding:"omitempty,max=50"`
	Phone   string `json:"phone" binding:"omitempty,max=20"`
	Email   string `json:"email" binding:"omitempty,email"`
	Address string `json:"address" binding:"omitempty,max=200"`
}

// SupplierResponse 供应商响应
type SupplierResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Contact   string    `json:"contact"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SupplierListResponse 供应商列表响应
type SupplierListResponse struct {
	Total   int64             `json:"total"`
	Suppliers []*SupplierResponse `json:"suppliers"`
}