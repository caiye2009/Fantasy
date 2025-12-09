// back/internal/analytics/domain/return_analysis_domain.go
package domain

import (
	"context"
	"time"
)

// ReturnAnalysisQuery 退货分析查询参数
type ReturnAnalysisQuery struct {
	CustomerNo string    // 客户编号，空表示所有客户
	StartDate  time.Time // 开始日期
	EndDate    time.Time // 结束日期
}

// MeterDimension 米数维度统计
type MeterDimension struct {
	TotalMeters    float64 `json:"totalMeters"`    // 订单总米数
	ReturnedMeters float64 `json:"returnedMeters"` // 退货总米数
	ReturnRate     string  `json:"returnRate"`     // 退货率（可能是 "N/A"）
	OrderCount     int64   `json:"orderCount"`     // 涉及订单数
}

// WeightDimension 重量维度统计
type WeightDimension struct {
	TotalWeight    float64 `json:"totalWeight"`    // 订单总重量（公斤）
	ReturnedWeight float64 `json:"returnedWeight"` // 退货总重量
	ReturnRate     string  `json:"returnRate"`     // 退货率（可能是 "N/A"）
	OrderCount     int64   `json:"orderCount"`     // 涉及订单数
}

// AmountDimension 金额维度统计
type AmountDimension struct {
	TotalAmountRMB    float64 `json:"totalAmountRMB"`    // 退款总金额（统一换算为RMB）
	ReturnedOrderCount int64   `json:"returnedOrderCount"` // 有退货的订单数
}

// ReturnAnalysis 退货分析结果
type ReturnAnalysis struct {
	CustomerNo   string          `json:"customerNo"`   // 客户编号
	CustomerName string          `json:"customerName"` // 客户名称
	MeterStats   MeterDimension  `json:"meterStats"`   // 米数统计
	WeightStats  WeightDimension `json:"weightStats"`  // 重量统计
	AmountStats  AmountDimension `json:"amountStats"`  // 金额统计
	TotalOrders  int64           `json:"totalOrders"`  // 总订单数
}

// CustomerOption 客户选项（用于下拉列表）
type CustomerOption struct {
	CustomerNo   string `json:"customerNo"`   // 客户编号
	CustomerName string `json:"customerName"` // 客户名称
}

// OrderSummary 订单汇总数据（内部使用）
type OrderSummary struct {
	TotalOrders       int64
	TotalMeters       float64
	TotalWeight       float64
	MeterOrderCount   int64
	WeightOrderCount  int64
}

// ReturnSummary 退货汇总数据（内部使用）
type ReturnSummary struct {
	ReturnedMeters      float64
	ReturnedWeight      float64
	ReturnAmountRMB     float64 // 统一换算为RMB
	ReturnedOrderCount  int64
}

// ReturnAnalysisRepository 退货分析仓储接口
type ReturnAnalysisRepository interface {
	// GetCustomerList 获取客户列表（用于下拉选择）
	GetCustomerList(ctx context.Context) ([]CustomerOption, error)

	// GetOrderSummary 获取订单汇总数据
	GetOrderSummary(ctx context.Context, query ReturnAnalysisQuery) (*OrderSummary, error)

	// GetReturnSummary 获取退货汇总数据
	GetReturnSummary(ctx context.Context, query ReturnAnalysisQuery) (*ReturnSummary, error)

	// GetCustomerInfo 获取客户信息
	GetCustomerInfo(ctx context.Context, customerNo string) (*CustomerOption, error)
}