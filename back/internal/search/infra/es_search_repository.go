package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"

	"back/internal/search/domain"
)

// ESSearchRepository ES 搜索仓储实现
type ESSearchRepository struct {
	client   *elasticsearch.Client
	queryBuilder *QueryBuilder
	aggBuilder   *AggregationBuilder
}

// NewESSearchRepository 创建 ES 搜索仓储
func NewESSearchRepository(client *elasticsearch.Client) *ESSearchRepository {
	return &ESSearchRepository{
		client:   client,
		queryBuilder: NewQueryBuilder(),
		aggBuilder:   NewAggregationBuilder(),
	}
}

// Search 执行搜索
func (r *ESSearchRepository) Search(
	ctx context.Context,
	criteria *domain.SearchCriteria,
	config *domain.SearchConfig,
) (*domain.SearchResult, error) {
	// 1. 构建 ES 查询 DSL
	esQuery := r.buildESQuery(criteria, config)

	// 2. 执行查询
	esResponse, err := r.executeSearch(ctx, config.Index, esQuery)
	if err != nil {
		return nil, err
	}

	// 3. 解析响应
	result, err := r.parseResponse(esResponse, criteria)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// buildESQuery 构建完整的 ES 查询
func (r *ESSearchRepository) buildESQuery(
	criteria *domain.SearchCriteria,
	config *domain.SearchConfig,
) map[string]interface{} {
	esQuery := map[string]interface{}{}

	// 1. Query 部分
	esQuery["query"] = r.queryBuilder.Build(criteria, config)

	// 2. 分页
	esQuery["from"] = criteria.Pagination.Offset
	esQuery["size"] = criteria.Pagination.Size

	// 3. 排序
	if len(criteria.Sort) > 0 {
		sortArray := make([]map[string]interface{}, 0, len(criteria.Sort))
		for _, s := range criteria.Sort {
			sortConfig := map[string]interface{}{
				"order": s.Order,
			}
			// 处理缺失值
			if s.Missing != "" {
				sortConfig["missing"] = s.Missing
			}
			sortArray = append(sortArray, map[string]interface{}{
				s.Field: sortConfig,
			})
		}
		esQuery["sort"] = sortArray
	}

	// 4. 聚合
	if len(criteria.AggRequests) > 0 {
		aggs := r.aggBuilder.Build(criteria.AggRequests, config, criteria.Filters)
		if aggs != nil {
			esQuery["aggs"] = aggs
		}
	}

	// 5. 高亮
	esQuery["highlight"] = map[string]interface{}{
		"fields": map[string]interface{}{
			"*": map[string]interface{}{
				"pre_tags":  []string{"<mark>"},
				"post_tags": []string{"</mark>"},
			},
		},
	}

	return esQuery
}

// executeSearch 执行 ES 查询
func (r *ESSearchRepository) executeSearch(
	ctx context.Context,
	index string,
	esQuery map[string]interface{},
) (map[string]interface{}, error) {
	// 序列化查询
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, fmt.Errorf("encode query: %w", err)
	}

	// 执行搜索
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(index),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		errorBody := res.String()
		return nil, fmt.Errorf("search error: %s", errorBody)
	}

	// 解析响应
	var esResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return esResponse, nil
}

// parseResponse 解析 ES 响应
func (r *ESSearchRepository) parseResponse(
	esResponse map[string]interface{},
	criteria *domain.SearchCriteria,
) (*domain.SearchResult, error) {
	result := &domain.SearchResult{
		Items:        []map[string]interface{}{},
		Aggregations: map[string]domain.AggResult{},
	}

	// 1. 解析命中结果
	if hits, ok := esResponse["hits"].(map[string]interface{}); ok {
		// 总数
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				result.Total = int(value)
			}
		}

		// 文档列表
		if hitsList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					if source, ok := hitMap["_source"].(map[string]interface{}); ok {
						result.Items = append(result.Items, source)
					}
				}
			}
		}
	}

	// 2. 耗时
	if took, ok := esResponse["took"].(float64); ok {
		result.Took = int(took)
	}

	// 3. 解析聚合结果
	if aggs, ok := esResponse["aggregations"].(map[string]interface{}); ok {
		result.Aggregations = r.parseAggregations(aggs, criteria)
	}

	return result, nil
}

// parseAggregations 解析聚合结果
func (r *ESSearchRepository) parseAggregations(
	aggs map[string]interface{},
	criteria *domain.SearchCriteria,
) map[string]domain.AggResult {
	results := map[string]domain.AggResult{}

	for aggField, aggData := range aggs {
		aggDataMap, ok := aggData.(map[string]interface{})
		if !ok {
			continue
		}

		// 检查是否是 filter aggregation 包装
		if _, hasFilter := aggDataMap["doc_count"]; hasFilter {
			// 这是 filter aggregation，提取内层聚合
			innerAggName := aggField + "_inner"
			if innerAgg, ok := aggDataMap[innerAggName].(map[string]interface{}); ok {
				aggDataMap = innerAgg
			}
		}

		// 根据聚合类型解析
		var aggResult domain.AggResult

		// Terms 聚合
		if buckets, ok := aggDataMap["buckets"].([]interface{}); ok {
			aggResult.Buckets = r.parseBuckets(buckets)
			aggResult.HasMore = false // terms 聚合没有分页
		}

		// Composite 聚合
		if buckets, ok := aggDataMap["buckets"].([]interface{}); ok {
			aggResult.Buckets = r.parseBuckets(buckets)
			// 检查是否有 after_key
			if afterKey, ok := aggDataMap["after_key"].(map[string]interface{}); ok {
				aggResult.After = afterKey
				// 如果有 after_key，说明可能还有更多数据
				aggReq, exists := criteria.AggRequests[aggField]
				if exists {
					aggResult.HasMore = len(aggResult.Buckets) >= aggReq.Size
				} else {
					aggResult.HasMore = true
				}
			} else {
				aggResult.HasMore = false
			}
		}

		// Stats 聚合
		if min, ok := aggDataMap["min"].(float64); ok {
			aggResult.Min = &min
		}
		if max, ok := aggDataMap["max"].(float64); ok {
			aggResult.Max = &max
		}
		if avg, ok := aggDataMap["avg"].(float64); ok {
			aggResult.Avg = &avg
		}

		results[aggField] = aggResult
	}

	return results
}

// parseBuckets 解析聚合桶
func (r *ESSearchRepository) parseBuckets(buckets []interface{}) []domain.Bucket {
	result := make([]domain.Bucket, 0, len(buckets))

	for _, bucket := range buckets {
		bucketMap, ok := bucket.(map[string]interface{})
		if !ok {
			continue
		}

		b := domain.Bucket{}

		// Key
		if key, ok := bucketMap["key"]; ok {
			b.Key = key
		}

		// DocCount
		if docCount, ok := bucketMap["doc_count"].(float64); ok {
			b.DocCount = int(docCount)
		}

		result = append(result, b)
	}

	return result
}
