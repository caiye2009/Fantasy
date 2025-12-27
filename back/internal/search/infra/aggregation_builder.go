package infra

import (
	"strings"

	"back/internal/search/domain"
)

// AggregationBuilder Aggregation DSL 构建器
type AggregationBuilder struct{}

// NewAggregationBuilder 创建 Aggregation 构建器
func NewAggregationBuilder() *AggregationBuilder {
	return &AggregationBuilder{}
}

// Build 构建 ES Aggregations
func (b *AggregationBuilder) Build(
	aggRequests map[string]domain.AggRequest,
	config *domain.SearchConfig,
	currentFilters map[string]interface{},
) map[string]interface{} {
	if len(aggRequests) == 0 {
		return nil
	}

	aggregations := map[string]interface{}{}

	for aggField, aggReq := range aggRequests {
		aggConfig := config.GetAggField(aggField)
		if aggConfig == nil {
			continue // 字段不在白名单，忽略
		}

		// 构建聚合
		agg := b.buildAggregation(aggField, aggReq, aggConfig, currentFilters, config)
		if agg != nil {
			aggregations[aggField] = agg
		}
	}

	return aggregations
}

// buildAggregation 构建单个聚合
func (b *AggregationBuilder) buildAggregation(
	aggField string,
	aggReq domain.AggRequest,
	aggConfig *domain.AggregationFieldConfig,
	currentFilters map[string]interface{},
	searchConfig *domain.SearchConfig,
) map[string]interface{} {
	// 1. 构建内层聚合
	innerAgg := b.buildInnerAggregation(aggField, aggReq, aggConfig)
	if innerAgg == nil {
		return nil
	}

	// 2. 如果需要联动（excludeSelf），包装 filter aggregation
	if aggConfig.ExcludeSelf && len(currentFilters) > 0 {
		filterConditions := b.buildFilterConditions(aggField, currentFilters, searchConfig)
		if len(filterConditions) > 0 {
			return map[string]interface{}{
				"filter": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": filterConditions,
					},
				},
				"aggs": map[string]interface{}{
					aggField + "_inner": innerAgg,
				},
			}
		}
	}

	// 不需要联动或没有 filter，直接返回内层聚合
	return innerAgg
}

// buildInnerAggregation 构建内层聚合
func (b *AggregationBuilder) buildInnerAggregation(
	aggField string,
	aggReq domain.AggRequest,
	aggConfig *domain.AggregationFieldConfig,
) map[string]interface{} {
	size := aggConfig.Size
	if aggReq.Size > 0 {
		size = aggReq.Size // 请求覆盖配置
	}

	switch aggConfig.AggType {
	case "terms":
		return b.buildTermsAggregation(aggField, aggReq, aggConfig, size)

	case "composite":
		return b.buildCompositeAggregation(aggField, aggReq, aggConfig, size)

	case "stats":
		return b.buildStatsAggregation(aggField)

	case "date_histogram":
		return b.buildDateHistogramAggregation(aggField, size)

	default:
		return nil
	}
}

// buildTermsAggregation 构建 terms 聚合
func (b *AggregationBuilder) buildTermsAggregation(
	aggField string,
	aggReq domain.AggRequest,
	aggConfig *domain.AggregationFieldConfig,
	size int,
) map[string]interface{} {
	// 对于 text/keyword 类型（在 ES 中是 text + keyword multi-field），聚合需要使用 .keyword 子字段
	fieldName := aggField
	if aggConfig.Type == "text" || aggConfig.Type == "keyword" {
		fieldName = aggField + ".keyword"
	}

	termsAgg := map[string]interface{}{
		"field": fieldName,
		"size":  size,
	}

	// 下拉框搜索（正则匹配）
	if aggConfig.SupportSearch && aggReq.Search != "" {
		termsAgg["include"] = ".*" + aggReq.Search + ".*"
	}

	return map[string]interface{}{
		"terms": termsAgg,
	}
}

// buildCompositeAggregation 构建 composite 聚合（支持深度分页）
func (b *AggregationBuilder) buildCompositeAggregation(
	aggField string,
	aggReq domain.AggRequest,
	aggConfig *domain.AggregationFieldConfig,
	size int,
) map[string]interface{} {
	// 对于 text/keyword 类型（在 ES 中是 text + keyword multi-field），聚合需要使用 .keyword 子字段
	fieldName := aggField
	if aggConfig.Type == "text" || aggConfig.Type == "keyword" {
		fieldName = aggField + ".keyword"
	}

	termsSource := map[string]interface{}{
		"field": fieldName,
	}

	// 下拉框搜索（正则匹配）
	if aggConfig.SupportSearch && aggReq.Search != "" {
		termsSource["include"] = ".*" + aggReq.Search + ".*"
	}

	compositeAgg := map[string]interface{}{
		"sources": []map[string]interface{}{
			{
				aggField: map[string]interface{}{
					"terms": termsSource,
				},
			},
		},
		"size": size,
	}

	// 分页游标
	if aggReq.After != nil && len(aggReq.After) > 0 {
		compositeAgg["after"] = aggReq.After
	}

	return map[string]interface{}{
		"composite": compositeAgg,
	}
}

// buildStatsAggregation 构建 stats 聚合（数值统计）
func (b *AggregationBuilder) buildStatsAggregation(aggField string) map[string]interface{} {
	return map[string]interface{}{
		"stats": map[string]interface{}{
			"field": aggField,
		},
	}
}

// buildDateHistogramAggregation 构建日期直方图聚合
func (b *AggregationBuilder) buildDateHistogramAggregation(aggField string, size int) map[string]interface{} {
	return map[string]interface{}{
		"date_histogram": map[string]interface{}{
			"field":    aggField,
			"interval": "day", // 可以根据需要配置
		},
	}
}

// buildFilterConditions 构建联动 filter 条件（排除当前聚合字段自身）
func (b *AggregationBuilder) buildFilterConditions(
	aggField string,
	currentFilters map[string]interface{},
	searchConfig *domain.SearchConfig,
) []interface{} {
	conditions := []interface{}{}

	for filterField, filterValue := range currentFilters {
		// 跳过当前聚合字段自身及其范围查询后缀
		if b.isRelatedField(filterField, aggField) {
			continue
		}

		// 获取字段配置
		fieldConfig := searchConfig.GetFilterField(filterField)
		if fieldConfig == nil {
			// 检查是否是范围查询的 Min/Max 后缀
			baseField := b.extractBaseField(filterField)
			if baseField != "" {
				fieldConfig = searchConfig.GetFilterField(baseField)
				// 如果基础字段是当前聚合字段，跳过
				if baseField == aggField {
					continue
				}
			}
		}

		if fieldConfig == nil {
			continue
		}

		// 根据 operator 构建 filter 条件
		var condition map[string]interface{}
		switch fieldConfig.Operator {
		case "term":
			condition = map[string]interface{}{
				"term": map[string]interface{}{
					filterField: filterValue,
				},
			}
		case "terms":
			condition = map[string]interface{}{
				"terms": map[string]interface{}{
					filterField: filterValue,
				},
			}
		case "match":
			condition = map[string]interface{}{
				"match": map[string]interface{}{
					filterField: filterValue,
				},
			}
		case "range":
			// 范围查询需要特殊处理（在 QueryBuilder 中已处理）
			rangeClause := b.buildRangeClause(fieldConfig.Field, currentFilters)
			if rangeClause != nil && fieldConfig.Field != aggField {
				condition = rangeClause
			}
		}

		if condition != nil {
			conditions = append(conditions, condition)
		}
	}

	return conditions
}

// isRelatedField 判断是否为关联字段（如 price, priceMin, priceMax）
func (b *AggregationBuilder) isRelatedField(filterField, aggField string) bool {
	return filterField == aggField ||
		filterField == aggField+"Min" ||
		filterField == aggField+"Max"
}

// extractBaseField 提取基础字段名（去除 Min/Max 后缀）
func (b *AggregationBuilder) extractBaseField(field string) string {
	if strings.HasSuffix(field, "Min") {
		return strings.TrimSuffix(field, "Min")
	}
	if strings.HasSuffix(field, "Max") {
		return strings.TrimSuffix(field, "Max")
	}
	return ""
}

// buildRangeClause 构建范围查询子句
func (b *AggregationBuilder) buildRangeClause(baseField string, filters map[string]interface{}) map[string]interface{} {
	minKey := baseField + "Min"
	maxKey := baseField + "Max"

	rangeCondition := map[string]interface{}{}
	if minVal, ok := filters[minKey]; ok {
		rangeCondition["gte"] = minVal
	}
	if maxVal, ok := filters[maxKey]; ok {
		rangeCondition["lte"] = maxVal
	}

	if len(rangeCondition) == 0 {
		return nil
	}

	return map[string]interface{}{
		"range": map[string]interface{}{
			baseField: rangeCondition,
		},
	}
}
