package application

import "time"

// CreateMaterialRequest 创建材料请求
type CreateMaterialRequest struct {
	Code         string  `json:"code" binding:"omitempty,max=50"`          // 原料编号（可选，不填自动生成）
	Name         string  `json:"name" binding:"required,min=2,max=100"`
	Spec         string  `json:"spec" binding:"omitempty,max=200"`
	Unit         string  `json:"unit" binding:"omitempty,max=20"`
	Category     string  `json:"category" binding:"omitempty,max=50"`      // 分类
	CurrentPrice float64 `json:"currentPrice" binding:"omitempty,min=0"`   // 当前价格
	Description  string  `json:"description" binding:"omitempty"`
}

// UpdateMaterialRequest 更新材料请求
type UpdateMaterialRequest struct {
	Code         string  `json:"code" binding:"omitempty,max=50"`
	Name         string  `json:"name" binding:"omitempty,min=2,max=100"`
	Spec         string  `json:"spec" binding:"omitempty,max=200"`
	Unit         string  `json:"unit" binding:"omitempty,max=20"`
	Category     string  `json:"category" binding:"omitempty,max=50"`
	CurrentPrice float64 `json:"currentPrice" binding:"omitempty,min=0"`
	Description  string  `json:"description" binding:"omitempty"`
}

// MaterialResponse 材料响应
type MaterialResponse struct {
	ID           uint      `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Spec         string    `json:"spec"`
	Unit         string    `json:"unit"`
	Category     string    `json:"category"`
	CurrentPrice float64   `json:"currentPrice"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// MaterialListResponse 材料列表响应
type MaterialListResponse struct {
	Total     int64               `json:"total"`
	Materials []*MaterialResponse `json:"materials"`
}