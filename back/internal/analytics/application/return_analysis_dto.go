// back/internal/analytics/application/return_analysis_dto.go
package application

import (
	"fmt"
	"time"
)

// DateRange 日期范围
type DateRange struct {
	Start string `json:"start" example:"2024-01-01"` // 开始日期 YYYY-MM-DD，可选
	End   string `json:"end" example:"2024-12-31"`   // 结束日期 YYYY-MM-DD，可选
}

// Validate 验证日期范围
func (d *DateRange) Validate() error {
	// 检查是否只输入了一个日期
	hasStart := d.Start != ""
	hasEnd := d.End != ""
	
	if hasStart && !hasEnd {
		return fmt.Errorf("开始日期和结束日期必须同时提供，或者都不提供")
	}
	
	if !hasStart && hasEnd {
		return fmt.Errorf("开始日期和结束日期必须同时提供，或者都不提供")
	}
	
	// 如果都没有输入，表示查询全部数据，验证通过
	if !hasStart && !hasEnd {
		return nil
	}
	
	// 验证日期格式
	startDate, err := time.Parse("2006-01-02", d.Start)
	if err != nil {
		return fmt.Errorf("开始日期格式错误: %w", err)
	}
	
	endDate, err := time.Parse("2006-01-02", d.End)
	if err != nil {
		return fmt.Errorf("结束日期格式错误: %w", err)
	}
	
	if endDate.Before(startDate) {
		return fmt.Errorf("结束日期不能早于开始日期")
	}
	
	return nil
}

// ParseDateRange 解析日期范围
func (d *DateRange) ParseDateRange() (time.Time, time.Time, error) {
	// 如果都没有输入，返回零值（表示不过滤日期）
	if d.Start == "" && d.End == "" {
		return time.Time{}, time.Time{}, nil
	}
	
	startDate, err := time.Parse("2006-01-02", d.Start)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("解析开始日期失败: %w", err)
	}
	
	endDate, err := time.Parse("2006-01-02", d.End)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("解析结束日期失败: %w", err)
	}
	
	return startDate, endDate, nil
}

// ReturnAnalysisRequest 退货分析请求
type ReturnAnalysisRequest struct {
	CustomerNo string    `json:"customerNo" form:"customerNo" example:"CS0678"` // 客户编号，空表示所有客户
	DateRange  DateRange `json:"dateRange"`                                      // 日期范围，start和end都不传表示查询全部时间
}

// MeterDimensionResponse 米数维度响应
type MeterDimensionResponse struct {
	TotalMeters    float64 `json:"totalMeters" example:"50000.0"`    // 订单总米数
	ReturnedMeters float64 `json:"returnedMeters" example:"5000.0"`  // 退货总米数
	ReturnRate     string  `json:"returnRate" example:"10.00%"`      // 退货率（百分比字符串或"N/A"）
	OrderCount     int64   `json:"orderCount" example:"80"`          // 涉及订单数
}

// WeightDimensionResponse 重量维度响应
type WeightDimensionResponse struct {
	TotalWeight    float64 `json:"totalWeight" example:"20000.0"`    // 订单总重量（公斤）
	ReturnedWeight float64 `json:"returnedWeight" example:"2000.0"`  // 退货总重量
	ReturnRate     string  `json:"returnRate" example:"10.00%"`      // 退货率（百分比字符串或"N/A"）
	OrderCount     int64   `json:"orderCount" example:"30"`          // 涉及订单数
}

// AmountDimensionResponse 金额维度响应
type AmountDimensionResponse struct {
	TotalAmountRMB     float64 `json:"totalAmountRMB" example:"120000.0"`     // 退款总金额（RMB，USD已按汇率8换算）
	ReturnedOrderCount int64   `json:"returnedOrderCount" example:"15"`       // 有退货的订单数
}

// ReturnAnalysisResponse 退货分析响应
type ReturnAnalysisResponse struct {
	QueryConditions QueryConditionsResponse  `json:"queryConditions"` // 查询条件回显
	MeterStats      MeterDimensionResponse   `json:"meterStats"`      // 米数统计
	WeightStats     WeightDimensionResponse  `json:"weightStats"`     // 重量统计
	AmountStats     AmountDimensionResponse  `json:"amountStats"`     // 金额统计
	TotalOrders     int64                    `json:"totalOrders" example:"100"` // 总订单数
}

// QueryConditionsResponse 查询条件回显
type QueryConditionsResponse struct {
	DateRange    DateRange `json:"dateRange"`              // 日期范围
	CustomerNo   string    `json:"customerNo,omitempty" example:"CS0678"`   // 客户编号
	CustomerName string    `json:"customerName,omitempty" example:"客户A"` // 客户名称
}

// CustomerOptionResponse 客户选项响应
type CustomerOptionResponse struct {
	CustomerNo   string `json:"customerNo" example:"CS0678"`   // 客户编号
	CustomerName string `json:"customerName" example:"客户A"` // 客户名称
}