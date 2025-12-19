package infra

import (
	"fmt"
	"strings"

	"back/internal/search/domain"
)

// QueryBuilder Query DSL 构建器
type QueryBuilder struct{}

// NewQueryBuilder 创建 Query 构建器
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// Build 构建 ES Query
func (b *QueryBuilder) Build(criteria *domain.SearchCriteria, config *domain.SearchConfig) map[string]interface{} {
	boolQuery := map[string]interface{}{
		"must":   []interface{}{},
		"filter": []interface{}{},
	}

	// 1. 构建全文搜索（Query 字段）
	if criteria.Query != "" {
		mustClause := b.buildQueryClause(criteria.Query, config.QueryFields)
		if mustClause != nil {
			boolQuery["must"] = append(boolQuery["must"].([]interface{}), mustClause)
		}
	}

	// 2. 构建结构化筛选（Filter 字段）
	filterClauses := b.buildFilterClauses(criteria.Filters, config)
	if len(filterClauses) > 0 {
		boolQuery["filter"] = filterClauses
	}

	// 3. 如果没有任何查询条件，使用 match_all
	if len(boolQuery["must"].([]interface{})) == 0 && len(boolQuery["filter"].([]interface{})) == 0 {
		return map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	return map[string]interface{}{
		"bool": boolQuery,
	}
}

// buildQueryClause 构建全文搜索子句
func (b *QueryBuilder) buildQueryClause(queryText string, queryFields []domain.QueryFieldConfig) map[string]interface{} {
	if len(queryFields) == 0 {
		return nil
	}

	// 构建带权重的字段列表
	fields := make([]string, 0, len(queryFields))
	for _, qf := range queryFields {
		if qf.Boost > 0 && qf.Boost != 1.0 {
			fields = append(fields, fmt.Sprintf("%s^%.1f", qf.Field, qf.Boost))
		} else {
			fields = append(fields, qf.Field)
		}
	}

	return map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query":  queryText,
			"fields": fields,
			"type":   "best_fields",
		},
	}
}

// buildFilterClauses 构建 Filter 子句列表
func (b *QueryBuilder) buildFilterClauses(filters map[string]interface{}, config *domain.SearchConfig) []interface{} {
	if len(filters) == 0 {
		return nil
	}

	clauses := []interface{}{}
	processedFields := make(map[string]bool) // 标记已处理的字段（避免重复处理范围查询）

	for filterField, filterValue := range filters {
		// 跳过已处理的字段
		if processedFields[filterField] {
			continue
		}

		// 获取字段配置
		fieldConfig := config.GetFilterField(filterField)
		if fieldConfig == nil {
			// 检查是否是范围查询的 Min/Max 后缀
			baseField := b.extractBaseField(filterField)
			if baseField != "" {
				fieldConfig = config.GetFilterField(baseField)
			}
		}

		if fieldConfig == nil {
			continue // 字段不在白名单，忽略
		}

		// 根据 operator 类型构建子句
		switch fieldConfig.Operator {
		case "term":
			// 对于 text/keyword 类型（在 ES 中是 text + keyword multi-field），term 查询需要使用 .keyword 子字段
			fieldName := filterField
			if fieldConfig.Type == "text" || fieldConfig.Type == "keyword" {
				fieldName = filterField + ".keyword"
			}
			clauses = append(clauses, map[string]interface{}{
				"term": map[string]interface{}{
					fieldName: filterValue,
				},
			})

		case "terms":
			// 对于 text/keyword 类型（在 ES 中是 text + keyword multi-field），terms 查询需要使用 .keyword 子字段
			fieldName := filterField
			if fieldConfig.Type == "text" || fieldConfig.Type == "keyword" {
				fieldName = filterField + ".keyword"
			}
			clauses = append(clauses, map[string]interface{}{
				"terms": map[string]interface{}{
					fieldName: filterValue,
				},
			})

		case "match":
			clauses = append(clauses, map[string]interface{}{
				"match": map[string]interface{}{
					filterField: filterValue,
				},
			})

		case "range":
			// 构建范围查询（处理 Min/Max 后缀）
			rangeClause := b.buildRangeClause(fieldConfig.Field, filters)
			if rangeClause != nil {
				clauses = append(clauses, rangeClause)
				// 标记相关字段已处理
				processedFields[fieldConfig.Field] = true
				processedFields[fieldConfig.Field+"Min"] = true
				processedFields[fieldConfig.Field+"Max"] = true
			}

		case "exists":
			clauses = append(clauses, map[string]interface{}{
				"exists": map[string]interface{}{
					"field": filterField,
				},
			})
		}
	}

	return clauses
}

// buildRangeClause 构建范围查询子句
func (b *QueryBuilder) buildRangeClause(baseField string, filters map[string]interface{}) map[string]interface{} {
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

// extractBaseField 提取基础字段名（去除 Min/Max 后缀）
func (b *QueryBuilder) extractBaseField(field string) string {
	if strings.HasSuffix(field, "Min") {
		return strings.TrimSuffix(field, "Min")
	}
	if strings.HasSuffix(field, "Max") {
		return strings.TrimSuffix(field, "Max")
	}
	return ""
}
