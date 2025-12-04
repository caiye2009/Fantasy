package application

import (
	"time"
	
	"back/internal/product/domain"
)

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
	Name      string                   `json:"name" binding:"required,min=2,max=100"`
	Materials []domain.MaterialConfig  `json:"materials" binding:"required,dive"`
	Processes []domain.ProcessConfig   `json:"processes" binding:"required,dive"`
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
	Name      string                   `json:"name" binding:"omitempty,min=2,max=100"`
	Status    string                   `json:"status" binding:"omitempty,oneof=draft submitted approved rejected"`
	Materials []domain.MaterialConfig  `json:"materials" binding:"omitempty,dive"`
	Processes []domain.ProcessConfig   `json:"processes" binding:"omitempty,dive"`
}

// ProductResponse 产品响应
type ProductResponse struct {
	ID        uint                     `json:"id"`
	Name      string                   `json:"name"`
	Status    string                   `json:"status"`
	Materials []domain.MaterialConfig  `json:"materials"`
	Processes []domain.ProcessConfig   `json:"processes"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

// ProductListResponse 产品列表响应
type ProductListResponse struct {
	Total    int64              `json:"total"`
	Products []*ProductResponse `json:"products"`
}

// CalculateCostRequest 计算成本请求
type CalculateCostRequest struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	UseMinPrice bool    `json:"use_min_price"`
}