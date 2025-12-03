package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type SearchService struct {
	client *elasticsearch.Client
}

func NewSearchService(client *elasticsearch.Client) *SearchService {
	return &SearchService{client: client}
}

// Search 通用搜索方法
func (s *SearchService) Search(req SearchRequest) (*SearchResponse, error) {
	// 默认值处理
	if req.Size == 0 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100 // 限制最大返回数量
	}

	// 确定搜索的索引
	indices := req.Indices
	if len(indices) == 0 {
		indices = GetAllIndices() // 搜索所有索引
	}

	// 确定搜索的字段
	fields := req.Fields
	if len(fields) == 0 {
		fields = GetDefaultFields(indices) // 使用默认字段
	}

	// 构建查询
	query := s.buildQuery(req, fields)

	// 执行搜索
	return s.executeSearch(indices, query)
}

// buildQuery 构建 ES 查询
func (s *SearchService) buildQuery(req SearchRequest, fields []string) map[string]interface{} {
	query := map[string]interface{}{
		"from": req.From,
		"size": req.Size,
	}

	// 构建 bool 查询
	boolQuery := map[string]interface{}{
		"must":   []interface{}{},
		"filter": []interface{}{},
	}

	// 添加关键词搜索
	if req.Query != "" {
		boolQuery["must"] = append(
			boolQuery["must"].([]interface{}),
			map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  req.Query,
					"fields": fields,
					"type":   "best_fields",
					"fuzziness": "AUTO", // 模糊匹配
				},
			},
		)
	}

	// 添加筛选条件
	for field, value := range req.Filters {
		boolQuery["filter"] = append(
			boolQuery["filter"].([]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{
					field: value,
				},
			},
		)
	}

	// 如果有查询条件，添加 bool 查询
	if req.Query != "" || len(req.Filters) > 0 {
		query["query"] = map[string]interface{}{
			"bool": boolQuery,
		}
	} else {
		// 没有查询条件，返回所有
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	// 添加排序
	if len(req.Sort) > 0 {
		sortArray := make([]map[string]interface{}, 0, len(req.Sort))
		for _, s := range req.Sort {
			sortArray = append(sortArray, map[string]interface{}{
				s.Field: map[string]interface{}{
					"order": s.Order,
				},
			})
		}
		query["sort"] = sortArray
	}

	// 添加高亮
	query["highlight"] = map[string]interface{}{
		"fields": map[string]interface{}{
			"*": map[string]interface{}{
				"pre_tags":  []string{"<mark>"},
				"post_tags": []string{"</mark>"},
			},
		},
	}

	return query
}

// executeSearch 执行搜索
func (s *SearchService) executeSearch(indices []string, query map[string]interface{}) (*SearchResponse, error) {
	// 序列化查询
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("encode query: %w", err)
	}

	// 执行搜索
	indexNames := strings.Join(indices, ",")
	res, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(indexNames),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	// 解析响应
	return s.parseSearchResponse(res)
}

// parseSearchResponse 解析搜索响应
func (s *SearchService) parseSearchResponse(res *esapi.Response) (*SearchResponse, error) {
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

	response := &SearchResponse{
		Total:    result.Hits.Total.Value,
		Took:     result.Took,
		MaxScore: result.Hits.MaxScore,
		Results:  make([]SearchResult, 0, len(result.Hits.Hits)),
	}

	for _, hit := range result.Hits.Hits {
		// 根据索引名确定类型
		dataType := hit.Index
		if meta, ok := IndexConfig[hit.Index]; ok {
			dataType = meta.Type
		}

		response.Results = append(response.Results, SearchResult{
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

// IndexDocument 索引单个文档
func (s *SearchService) IndexDocument(index, id string, doc interface{}) error {
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

	res, err := req.Do(context.Background(), s.client)
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
func (s *SearchService) UpdateDocument(index, id string, doc interface{}) error {
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

	res, err := req.Do(context.Background(), s.client)
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
func (s *SearchService) DeleteDocument(index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), s.client)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete error: %s", res.String())
	}

	return nil
}