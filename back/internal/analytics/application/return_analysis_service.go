// back/internal/analytics/application/return_analysis_service.go
package application

import (
	"context"
	"fmt"

	"back/internal/analytics/domain"
)

// ReturnAnalysisService 退货分析服务
type ReturnAnalysisService struct {
	repo domain.ReturnAnalysisRepository
}

// NewReturnAnalysisService 创建退货分析服务
func NewReturnAnalysisService(repo domain.ReturnAnalysisRepository) *ReturnAnalysisService {
	return &ReturnAnalysisService{
		repo: repo,
	}
}

// GetCustomerList 获取客户列表
func (s *ReturnAnalysisService) GetCustomerList(ctx context.Context) ([]CustomerOptionResponse, error) {
	customers, err := s.repo.GetCustomerList(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取客户列表失败: %w", err)
	}

	resp := make([]CustomerOptionResponse, len(customers))
	for i, c := range customers {
		resp[i] = CustomerOptionResponse{
			CustomerNo:   c.CustomerNo,
			CustomerName: c.CustomerName,
		}
	}

	return resp, nil
}

// GetReturnAnalysis 获取退货分析
func (s *ReturnAnalysisService) GetReturnAnalysis(ctx context.Context, req *ReturnAnalysisRequest) (*ReturnAnalysisResponse, error) {
	// 验证日期范围
	if err := req.DateRange.Validate(); err != nil {
		return nil, fmt.Errorf("日期范围验证失败: %w", err)
	}

	// 解析日期
	startDate, endDate, err := req.DateRange.ParseDateRange()
	if err != nil {
		return nil, err
	}

	// 构建查询
	query := domain.ReturnAnalysisQuery{
		CustomerNo: req.CustomerNo,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	// 获取订单汇总
	orderSummary, err := s.repo.GetOrderSummary(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查询订单汇总失败: %w", err)
	}

	// 获取退货汇总
	returnSummary, err := s.repo.GetReturnSummary(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查询退货汇总失败: %w", err)
	}

	// 获取客户信息（如果指定了客户）
	var customerName string
	if req.CustomerNo != "" {
		customerInfo, err := s.repo.GetCustomerInfo(ctx, req.CustomerNo)
		if err == nil && customerInfo != nil {
			customerName = customerInfo.CustomerName
		}
	}

	// 计算米数维度
	meterStats := s.calculateMeterDimension(orderSummary, returnSummary)

	// 计算重量维度
	weightStats := s.calculateWeightDimension(orderSummary, returnSummary)

	// 金额维度
	amountStats := AmountDimensionResponse{
		TotalAmountRMB:     returnSummary.ReturnAmountRMB,
		ReturnedOrderCount: returnSummary.ReturnedOrderCount,
	}

	// 构建响应
	return &ReturnAnalysisResponse{
		QueryConditions: QueryConditionsResponse{
			DateRange:    req.DateRange,
			CustomerNo:   req.CustomerNo,
			CustomerName: customerName,
		},
		MeterStats:  meterStats,
		WeightStats: weightStats,
		AmountStats: amountStats,
		TotalOrders: orderSummary.TotalOrders,
	}, nil
}

// calculateMeterDimension 计算米数维度
func (s *ReturnAnalysisService) calculateMeterDimension(order *domain.OrderSummary, ret *domain.ReturnSummary) MeterDimensionResponse {
	resp := MeterDimensionResponse{
		TotalMeters:    order.TotalMeters,
		ReturnedMeters: ret.ReturnedMeters,
		OrderCount:     order.MeterOrderCount,
	}

	// 计算退货率
	if order.TotalMeters > 0 {
		rate := (ret.ReturnedMeters / order.TotalMeters) * 100
		resp.ReturnRate = fmt.Sprintf("%.2f%%", rate)
	} else {
		resp.ReturnRate = "N/A"
	}

	return resp
}

// calculateWeightDimension 计算重量维度
func (s *ReturnAnalysisService) calculateWeightDimension(order *domain.OrderSummary, ret *domain.ReturnSummary) WeightDimensionResponse {
	resp := WeightDimensionResponse{
		TotalWeight:    order.TotalWeight,
		ReturnedWeight: ret.ReturnedWeight,
		OrderCount:     order.WeightOrderCount,
	}

	// 计算退货率
	if order.TotalWeight > 0 {
		rate := (ret.ReturnedWeight / order.TotalWeight) * 100
		resp.ReturnRate = fmt.Sprintf("%.2f%%", rate)
	} else {
		resp.ReturnRate = "N/A"
	}

	return resp
}