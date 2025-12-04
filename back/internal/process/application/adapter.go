package application

import (
	"context"
	
	productApp "back/internal/product/application"
)

// ProcessServiceAdapter Process 服务适配器（供 Product 模块使用）
type ProcessServiceAdapter struct {
	service *ProcessService
}

// NewProcessServiceAdapter 创建适配器
func NewProcessServiceAdapter(service *ProcessService) productApp.ProcessServiceInterface {
	return &ProcessServiceAdapter{service: service}
}

// Get 获取工序信息（适配为 Product 模块需要的类型）
func (a *ProcessServiceAdapter) Get(ctx context.Context, id uint) (*productApp.ProcessInfo, error) {
	resp, err := a.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &productApp.ProcessInfo{
		ID:   resp.ID,
		Name: resp.Name,
	}, nil
}

// Exists 检查工序是否存在
func (a *ProcessServiceAdapter) Exists(ctx context.Context, id uint) (bool, error) {
	return a.service.Exists(ctx, id)
}