// back/internal/analytics/infra/return_analysis_repository.go
package infra

import (
	"context"
	"fmt"
	"strings"

	"back/internal/analytics/domain"
	"gorm.io/gorm"
)

// ReturnAnalysisRepositoryImpl 退货分析仓储实现
type ReturnAnalysisRepositoryImpl struct {
	db           *gorm.DB
	exchangeRate float64 // USD 到 RMB 的汇率
}

// NewReturnAnalysisRepository 创建退货分析仓储
func NewReturnAnalysisRepository(db *gorm.DB) domain.ReturnAnalysisRepository {
	return &ReturnAnalysisRepositoryImpl{
		db:           db,
		exchangeRate: 8.0, // 默认汇率
	}
}

// GetCustomerList 获取客户列表
func (r *ReturnAnalysisRepositoryImpl) GetCustomerList(ctx context.Context) ([]domain.CustomerOption, error) {
	var results []struct {
		CustomerNo   string `gorm:"column:CustomNo"`
		CustomerName string `gorm:"column:CustomerName"`
	}

	err := r.db.WithContext(ctx).
		Table("v_md_plan").
		Select(`DISTINCT "CustomNo", MAX("CustomerName") as "CustomerName"`).
		Where(`"Schedule" = ?`, 99).
		Group(`"CustomNo"`).
		Order(`"CustomerName" ASC`).
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("查询客户列表失败: %w", err)
	}

	customers := make([]domain.CustomerOption, len(results))
	for i, r := range results {
		customers[i] = domain.CustomerOption{
			CustomerNo:   r.CustomerNo,
			CustomerName: r.CustomerName,
		}
	}

	return customers, nil
}

// GetOrderSummary 获取订单汇总数据
func (r *ReturnAnalysisRepositoryImpl) GetOrderSummary(ctx context.Context, query domain.ReturnAnalysisQuery) (*domain.OrderSummary, error) {
	var result struct {
		TotalOrders      int64   `gorm:"column:total_orders"`
		TotalMeters      float64 `gorm:"column:total_meters"`
		TotalWeight      float64 `gorm:"column:total_weight"`
		MeterOrderCount  int64   `gorm:"column:meter_order_count"`
		WeightOrderCount int64   `gorm:"column:weight_order_count"`
	}

	// 构建查询
	db := r.db.WithContext(ctx).Table("v_md_plan")

	// 基础过滤条件
	db = db.Where(`"Schedule" = ?`, 99)

	// 客户过滤
	if query.CustomerNo != "" {
		db = db.Where(`"CustomNo" = ?`, query.CustomerNo)
	}

	// 日期范围过滤
	if !query.StartDate.IsZero() {
		db = db.Where(`TO_DATE("AffirmDate", 'YYYY-MM-DD') >= ?`, query.StartDate)
	}
	if !query.EndDate.IsZero() {
		db = db.Where(`TO_DATE("AffirmDate", 'YYYY-MM-DD') <= ?`, query.EndDate)
	}

	// 执行聚合查询
	err := db.Select(`
		COUNT(DISTINCT "Planno") as total_orders,
		COALESCE(SUM(CASE 
			WHEN LOWER("Unit") IN ('米', 'm', 'meter') THEN "TotalQuantity"
			ELSE 0 
		END), 0) as total_meters,
		COALESCE(SUM(CASE 
			WHEN LOWER("Unit") IN ('公斤', 'kg', 'kilogram') THEN "TotalQuantity"
			ELSE 0 
		END), 0) as total_weight,
		COUNT(DISTINCT CASE 
			WHEN LOWER("Unit") IN ('米', 'm', 'meter') THEN "Planno"
		END) as meter_order_count,
		COUNT(DISTINCT CASE 
			WHEN LOWER("Unit") IN ('公斤', 'kg', 'kilogram') THEN "Planno"
		END) as weight_order_count
	`).Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("查询订单汇总失败: %w", err)
	}

	return &domain.OrderSummary{
		TotalOrders:      result.TotalOrders,
		TotalMeters:      result.TotalMeters,
		TotalWeight:      result.TotalWeight,
		MeterOrderCount:  result.MeterOrderCount,
		WeightOrderCount: result.WeightOrderCount,
	}, nil
}

// GetReturnSummary 获取退货汇总数据
func (r *ReturnAnalysisRepositoryImpl) GetReturnSummary(ctx context.Context, query domain.ReturnAnalysisQuery) (*domain.ReturnSummary, error) {
	// 第一步：查询按单位分组的退货数量
	var quantityResults []struct {
		ReturnedMeters float64 `gorm:"column:returned_meters"`
		ReturnedWeight float64 `gorm:"column:returned_weight"`
	}

	quantityDB := r.db.WithContext(ctx).
		Table("\"Ord_ReturnGoodsApply\" r").
		Joins(`JOIN v_md_plan o ON r."PlanNo" = o."Planno"`).
		Where(`o."Schedule" = ?`, 99).
		Where(`r."Currency" IS NOT NULL AND r."Currency" != ''`)

	// 客户过滤
	if query.CustomerNo != "" {
		quantityDB = quantityDB.Where(`o."CustomNo" = ?`, query.CustomerNo)
	}

	// 日期范围过滤
	if !query.StartDate.IsZero() {
		quantityDB = quantityDB.Where(`TO_DATE(o."AffirmDate", 'YYYY-MM-DD') >= ?`, query.StartDate)
	}
	if !query.EndDate.IsZero() {
		quantityDB = quantityDB.Where(`TO_DATE(o."AffirmDate", 'YYYY-MM-DD') <= ?`, query.EndDate)
	}

	err := quantityDB.Select(`
		COALESCE(SUM(CASE 
			WHEN LOWER(o."Unit") IN ('米', 'm', 'meter') THEN r."TotalQuantity"
			ELSE 0 
		END), 0) as returned_meters,
		COALESCE(SUM(CASE 
			WHEN LOWER(o."Unit") IN ('公斤', 'kg', 'kilogram') THEN r."TotalQuantity"
			ELSE 0 
		END), 0) as returned_weight
	`).Scan(&quantityResults).Error

	if err != nil {
		return nil, fmt.Errorf("查询退货数量失败: %w", err)
	}

	var returnedMeters, returnedWeight float64
	if len(quantityResults) > 0 {
		returnedMeters = quantityResults[0].ReturnedMeters
		returnedWeight = quantityResults[0].ReturnedWeight
	}

	// 第二步：查询按币种分组的退款金额
	var amountResults []struct {
		Currency    string  `gorm:"column:currency"`
		TotalAmount float64 `gorm:"column:total_amount"`
	}

	amountDB := r.db.WithContext(ctx).
		Table("\"Ord_ReturnGoodsApply\" r").
		Joins(`JOIN v_md_plan o ON r."PlanNo" = o."Planno"`).
		Where(`o."Schedule" = ?`, 99).
		Where(`r."Currency" IS NOT NULL AND r."Currency" != ''`)

	// 客户过滤
	if query.CustomerNo != "" {
		amountDB = amountDB.Where(`o."CustomNo" = ?`, query.CustomerNo)
	}

	// 日期范围过滤
	if !query.StartDate.IsZero() {
		amountDB = amountDB.Where(`TO_DATE(o."AffirmDate", 'YYYY-MM-DD') >= ?`, query.StartDate)
	}
	if !query.EndDate.IsZero() {
		amountDB = amountDB.Where(`TO_DATE(o."AffirmDate", 'YYYY-MM-DD') <= ?`, query.EndDate)
	}

	err = amountDB.Select(`
		UPPER(TRIM(r."Currency")) as currency,
		COALESCE(SUM(r."TotalAmount"), 0) as total_amount
	`).
		Group("UPPER(TRIM(r.\"Currency\"))").
		Find(&amountResults).Error

	if err != nil {
		return nil, fmt.Errorf("查询退款金额失败: %w", err)
	}

	// 计算统一换算为 RMB 的总金额
	totalAmountRMB := 0.0
	for _, result := range amountResults {
		amount := result.TotalAmount
		currency := strings.ToUpper(strings.TrimSpace(result.Currency))

		// 根据币种换算为 RMB
		switch currency {
		case "RMB", "CNY", "人民币":
			totalAmountRMB += amount
		case "USD", "美元", "美金":
			totalAmountRMB += amount * r.exchangeRate
		default:
			// 其他币种暂不处理，可根据需要扩展
		}
	}

	// 第三步：查询有退货的订单数
	var orderCountResult struct {
		Count int64 `gorm:"column:count"`
	}

	countDB := r.db.WithContext(ctx).
		Table("\"Ord_ReturnGoodsApply\" r").
		Joins(`JOIN v_md_plan o ON r."PlanNo" = o."Planno"`).
		Where(`o."Schedule" = ?`, 99).
		Where(`r."Currency" IS NOT NULL AND r."Currency" != ''`)

	// 客户过滤
	if query.CustomerNo != "" {
		countDB = countDB.Where(`o."CustomNo" = ?`, query.CustomerNo)
	}

	// 日期范围过滤
	if !query.StartDate.IsZero() {
		countDB = countDB.Where(`TO_DATE(o."AffirmDate", 'YYYY-MM-DD') >= ?`, query.StartDate)
	}
	if !query.EndDate.IsZero() {
		countDB = countDB.Where(`TO_DATE(o."AffirmDate", 'YYYY-MM-DD') <= ?`, query.EndDate)
	}

	err = countDB.Select(`COUNT(DISTINCT r."PlanNo") as count`).
		Scan(&orderCountResult).Error

	if err != nil {
		return nil, fmt.Errorf("查询退货订单数失败: %w", err)
	}

	return &domain.ReturnSummary{
		ReturnedMeters:     returnedMeters,
		ReturnedWeight:     returnedWeight,
		ReturnAmountRMB:    totalAmountRMB,
		ReturnedOrderCount: orderCountResult.Count,
	}, nil
}

// GetCustomerInfo 获取客户信息
func (r *ReturnAnalysisRepositoryImpl) GetCustomerInfo(ctx context.Context, customerNo string) (*domain.CustomerOption, error) {
	var result struct {
		CustomerNo   string `gorm:"column:CustomNo"`
		CustomerName string `gorm:"column:CustomerName"`
	}

	err := r.db.WithContext(ctx).
		Table("v_md_plan").
		Select(`"CustomNo", MAX("CustomerName") as "CustomerName"`).
		Where(`"CustomNo" = ?`, customerNo).
		Where(`"Schedule" = ?`, 99).
		Group(`"CustomNo"`).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询客户信息失败: %w", err)
	}

	return &domain.CustomerOption{
		CustomerNo:   result.CustomerNo,
		CustomerName: result.CustomerName,
	}, nil
}