package application

import (
	"context"
	
	pricingApp "back/internal/pricing/application"
)

// SupplierServiceAdapter Supplier 服务适配器（供 Pricing 模块使用）
type SupplierServiceAdapter struct {
	service *SupplierService
}

// NewSupplierServiceAdapter 创建适配器
func NewSupplierServiceAdapter(service *SupplierService) pricingApp.SupplierServiceInterface {
	return &SupplierServiceAdapter{service: service}
}

// GetSupplierInfo 获取供应商信息（适配为 Pricing 模块需要的类型）
func (a *SupplierServiceAdapter) GetSupplierInfo(ctx context.Context, id uint) (*pricingApp.SupplierInfo, error) {
	info, err := a.service.GetSupplierInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &pricingApp.SupplierInfo{
		ID:   info.ID,
		Name: info.Name,
	}, nil
}