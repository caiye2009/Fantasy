package query

import (
	"fmt"
	"time"
)

// Operator 查询操作符
type Operator string

const (
	OpEqual          Operator = "eq"     // 等于
	OpNotEqual       Operator = "ne"     // 不等于
	OpGreaterThan    Operator = "gt"     // 大于
	OpGreaterOrEqual Operator = "gte"    // 大于等于
	OpLessThan       Operator = "lt"     // 小于
	OpLessOrEqual    Operator = "lte"    // 小于等于
	OpLike           Operator = "like"   // 模糊匹配
	OpIn             Operator = "in"     // 在范围内
	OpNotIn          Operator = "nin"    // 不在范围内
	OpBetween        Operator = "between" // 在两值之间
	OpIsNull         Operator = "null"   // 为空
	OpNotNull        Operator = "notnull" // 不为空
)

// Condition 查询条件
type Condition struct {
	Field    string      `json:"field"`    // 字段名
	Operator Operator    `json:"operator"` // 操作符
	Value    interface{} `json:"value"`    // 值
}

// LogicOperator 逻辑操作符
type LogicOperator string

const (
	LogicAnd LogicOperator = "AND"
	LogicOr  LogicOperator = "OR"
)

// Filter 查询过滤器
type Filter struct {
	Conditions []Condition   `json:"conditions"` // 条件列表
	Logic      LogicOperator `json:"logic"`      // 逻辑操作符，默认AND
}

// DateRange 日期范围
type DateRange struct {
	Start string `json:"start"` // 开始日期 YYYY-MM-DD
	End   string `json:"end"`   // 结束日期 YYYY-MM-DD
}

// ParseDateRange 解析日期范围
func (dr *DateRange) ParseDateRange() (time.Time, time.Time, error) {
	if dr.Start == "" || dr.End == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("日期范围不能为空")
	}

	start, err := time.Parse("2006-01-02", dr.Start)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("开始日期格式错误: %w", err)
	}

	end, err := time.Parse("2006-01-02", dr.End)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("结束日期格式错误: %w", err)
	}

	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("开始日期不能晚于结束日期")
	}

	return start, end, nil
}

// Validate 验证日期范围
func (dr *DateRange) Validate() error {
	_, _, err := dr.ParseDateRange()
	return err
}

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page" form:"page"`         // 页码，从1开始
	PageSize int `json:"pageSize" form:"pageSize"` // 每页数量
}

// GetOffset 获取数据库偏移量
func (p *Pagination) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

// GetLimit 获取限制数量
func (p *Pagination) GetLimit() int {
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// Sort 排序参数
type Sort struct {
	Field string `json:"field"` // 排序字段
	Order string `json:"order"` // asc 或 desc
}

// GetOrder 获取排序方式
func (s *Sort) GetOrder() string {
	if s.Order == "desc" {
		return "DESC"
	}
	return "ASC"
}

// QueryRequest 通用查询请求
type QueryRequest struct {
	Filter     *Filter     `json:"filter"`     // 过滤条件
	Pagination *Pagination `json:"pagination"` // 分页
	Sort       []Sort      `json:"sort"`       // 排序
}
