package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	
	"back/internal/search/domain"
)

// ESSearchRepo Elasticsearch 搜索仓储实现
type ESSearchRepo struct {
	client *elasticsearch.Client
}

// NewESSearchRepo 创建 ES 搜索仓储
func NewESSearchRepo(client *elasticsearch.Client) *ESSearchRepo {
	return &ESSearchRepo{client: client}
}

// Search 执行搜索
func (r *ESSearchRepo) Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResponse, error) {
	// 1. 构建 ES 查询
	esQuery := r.buildESQuery(query)
	
	// 2. 执行搜索
	return r.executeSearch(ctx, query.GetIndices(), esQuery)
}

// buildESQuery 构建 ES 查询
func (r *ESSearchRepo) buildESQuery(query *domain.SearchQuery) map[string]interface{} {
	esQuery := map[string]interface{}{
		"from": query.From,
		"size": query.Size,
	}
	
	// 构建 bool 查询
	boolQuery := map[string]interface{}{
		"must":   []interface{}{},
		"filter": []interface{}{},
	}
	
	// 添加关键词搜索
	if query.Keyword != "" {
		boolQuery["must"] = append(
			boolQuery["must"].([]interface{}),
			map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":     query.Keyword,
					"fields":    query.GetFields(),
					"type":      "best_fields",
					"fuzziness": "AUTO",
				},
			},
		)
	}
	
	// 添加筛选条件
	for field, value := range query.Filters {
		boolQuery["filter"] = append(
			boolQuery["filter"].([]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{
					field: value,
				},
			},
		)
	}
	
	// 添加查询
	if query.Keyword != "" || len(query.Filters) > 0 {
		esQuery["query"] = map[string]interface{}{
			"bool": boolQuery,
		}
	} else {
		esQuery["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}
	
	// 添加排序
	if len(query.Sort) > 0 {
		sortArray := make([]map[string]interface{}, 0, len(query.Sort))
		for _, s := range query.Sort {
			sortArray = append(sortArray, map[string]interface{}{
				s.Field: map[string]interface{}{
					"order": s.Order,
				},
			})
		}
		esQuery["sort"] = sortArray
	}
	
	// 添加高亮
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

// executeSearch 执行搜索
func (r *ESSearchRepo) executeSearch(ctx context.Context, indices []string, query map[string]interface{}) (*domain.SearchResponse, error) {
	// 序列化查询
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("encode query: %w", err)
	}
	
	// 执行搜索
	indexNames := strings.Join(indices, ",")
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(indexNames),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search request: %w", err)
	}
	defer res.Body.Close()
	
	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}
	
	// 解析响应
	return r.parseSearchResponse(res)
}

// parseSearchResponse 解析搜索响应
func (r *ESSearchRepo) parseSearchResponse(res *esapi.Response) (*domain.SearchResponse, error) {
	var result struct {
		Took int `json:"took"`
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			MaxScore float64 `json:"max_score"`
			Hits     []struct {
				Index     string                 `json:"_index"`
				ID        string                 `json:"_id"`
				Score     float64                `json:"_score"`
				Source    map[string]interface{} `json:"_source"`
				Highlight map[string][]string    `json:"highlight"`
			} `json:"hits"`
		} `json:"hits"`
	}
	
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	
	response := &domain.SearchResponse{
		Total:    result.Hits.Total.Value,
		Took:     result.Took,
		MaxScore: result.Hits.MaxScore,
		Results:  make([]*domain.SearchResult, 0, len(result.Hits.Hits)),
	}
	
	for _, hit := range result.Hits.Hits {
		// 根据索引名确定类型
		dataType := hit.Index
		if meta, ok := domain.GetIndexMeta(hit.Index); ok {
			dataType = meta.Type
		}
		
		response.Results = append(response.Results, &domain.SearchResult{
			Index:     hit.Index,
			Type:      dataType,
			ID:        hit.ID,
			Score:     hit.Score,
			Source:    hit.Source,
			Highlight: hit.Highlight,
		})
	}
	
	return response, nil
}

// IndexDocument 索引文档
func (r *ESSearchRepo) IndexDocument(ctx context.Context, index, id string, doc interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal document: %w", err)
	}
	
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}
	
	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("index document: %w", err)
	}
	defer res.Body.Close()
	
	if res.IsError() {
		return fmt.Errorf("index error: %s", res.String())
	}
	
	return nil
}

// UpdateDocument 更新文档
func (r *ESSearchRepo) UpdateDocument(ctx context.Context, index, id string, doc interface{}) error {
	data, err := json.Marshal(map[string]interface{}{
		"doc": doc,
	})
	if err != nil {
		return fmt.Errorf("marshal document: %w", err)
	}
	
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}
	
	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("update document: %w", err)
	}
	defer res.Body.Close()
	
	if res.IsError() {
		return fmt.Errorf("update error: %s", res.String())
	}
	
	return nil
}

// DeleteDocument 删除文档
func (r *ESSearchRepo) DeleteDocument(ctx context.Context, index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}
	
	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	defer res.Body.Close()
	
	if res.IsError() {
		return fmt.Errorf("delete error: %s", res.String())
	}
	
	return nil
}