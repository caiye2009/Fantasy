package application

import (
	"context"
	
	productApp "back/internal/product/application"
)

// MaterialServiceAdapter Material 服务适配器（供 Product 模块使用）
type MaterialServiceAdapter struct {
	service *MaterialService
}

// NewMaterialServiceAdapter 创建适配器
func NewMaterialServiceAdapter(service *MaterialService) productApp.MaterialServiceInterface {
	return &MaterialServiceAdapter{service: service}
}

// Get 获取材料信息（适配为 Product 模块需要的类型）
func (a *MaterialServiceAdapter) Get(ctx context.Context, id uint) (*productApp.MaterialInfo, error) {
	resp, err := a.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &productApp.MaterialInfo{
		ID:   resp.ID,
		Name: resp.Name,
	}, nil
}

// Exists 检查材料是否存在
func (a *MaterialServiceAdapter) Exists(ctx context.Context, id uint) (bool, error) {
	return a.service.Exists(ctx, id)
}