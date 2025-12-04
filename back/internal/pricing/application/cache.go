package application

import "context"

// PriceCache 价格缓存接口
type PriceCache interface {
	GetMin(ctx context.Context, targetType string, targetID uint) (*PriceData, error)
	SetMin(ctx context.Context, targetType string, targetID uint, data *PriceData) error
	UpdateMin(ctx context.Context, targetType string, targetID uint, newPrice *PriceData) error
	
	GetMax(ctx context.Context, targetType string, targetID uint) (*PriceData, error)
	SetMax(ctx context.Context, targetType string, targetID uint, data *PriceData) error
	UpdateMax(ctx context.Context, targetType string, targetID uint, newPrice *PriceData) error
}