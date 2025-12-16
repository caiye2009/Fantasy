package application

import (
	"context"
	"fmt"

	"back/internal/search/domain"
	"back/internal/search/infra"
)

// SearchService 搜索应用服务（全新实现）
type SearchService struct {
	registry *infra.SearchConfigRegistry
	repo     domain.SearchRepository
}

// NewSearchService 创建搜索服务
func NewSearchService(registry *infra.SearchConfigRegistry, repo domain.SearchRepository) *SearchService {
	return &SearchService{
		registry: registry,
		repo:     repo,
	}
}

// Search 执行搜索
func (s *SearchService) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	// 1. 加载配置
	config, ok := s.registry.GetConfig(req.EntityType)
	if !ok {
		return nil, fmt.Errorf("unsupported entity type: %s", req.EntityType)
	}

	// 2. 验证请求参数
	if err := s.validateRequest(config, req); err != nil {
		return nil, err
	}

	// 3. 构建搜索条件（DTO → Domain）
	criteria := s.buildCriteria(config, req)

	// 4. 执行搜索
	result, err := s.repo.Search(ctx, criteria, config)
	if err != nil {
		return nil, err
	}

	// 5. 格式化响应（Domain → DTO）
	return s.formatResponse(result), nil
}

// validateRequest 验证请求参数
func (s *SearchService) validateRequest(config *domain.SearchConfig, req *SearchRequest) error {
	// 验证 filter 字段是否在白名单
	for field := range req.Filters {
		if !config.IsFilterableField(field) {
			return fmt.Errorf("field '%s' is not filterable", field)
		}
	}

	// 验证 aggregation 字段是否在白名单
	for field := range req.AggRequests {
		if !config.IsAggregableField(field) {
			return fmt.Errorf("field '%s' is not aggregable", field)
		}
	}

	// 验证分页参数
	if req.Pagination.Size < 0 {
		return fmt.Errorf("pagination size must be >= 0")
	}
	if req.Pagination.Size > 100 {
		return fmt.Errorf("pagination size must be <= 100")
	}
	if req.Pagination.Offset < 0 {
		return fmt.Errorf("pagination offset must be >= 0")
	}

	return nil
}

// buildCriteria 构建搜索条件
func (s *SearchService) buildCriteria(config *domain.SearchConfig, req *SearchRequest) *domain.SearchCriteria {
	// 转换 Sort
	sortFields := make([]domain.SortField, len(req.Sort))
	for i, s := range req.Sort {
		sortFields[i] = domain.SortField{
			Field: s.Field,
			Order: s.Order,
		}
	}

	// 转换 AggRequests
	aggRequests := make(map[string]domain.AggRequest)
	for field, aggReq := range req.AggRequests {
		aggRequests[field] = domain.AggRequest{
			Search: aggReq.Search,
			After:  aggReq.After,
			Size:   aggReq.Size,
		}
	}

	// 初始化 Filters
	filters := req.Filters
	if filters == nil {
		filters = make(map[string]interface{})
	}

	// 默认分页
	pagination := req.Pagination
	if pagination.Size == 0 {
		pagination.Size = 20 // 默认 20 条
	}

	return &domain.SearchCriteria{
		EntityType:  req.EntityType,
		IndexName:   config.IndexName,
		Query:       req.Query,
		Filters:     filters,
		AggRequests: aggRequests,
		Pagination: domain.Pagination{
			Offset: pagination.Offset,
			Size:   pagination.Size,
		},
		Sort: sortFields,
	}
}

// formatResponse 格式化响应
func (s *SearchService) formatResponse(result *domain.SearchResult) *SearchResponse {
	// 转换 Aggregations
	aggregations := make(map[string]AggResult)
	for field, aggResult := range result.Aggregations {
		// 转换 Buckets
		buckets := make([]Bucket, len(aggResult.Buckets))
		for i, b := range aggResult.Buckets {
			buckets[i] = Bucket{
				Key:      b.Key,
				DocCount: b.DocCount,
			}
		}

		aggregations[field] = AggResult{
			Buckets: buckets,
			After:   aggResult.After,
			HasMore: aggResult.HasMore,
			Min:     aggResult.Min,
			Max:     aggResult.Max,
			Avg:     aggResult.Avg,
		}
	}

	return &SearchResponse{
		Items:        result.Items,
		Total:        result.Total,
		Took:         result.Took,
		Aggregations: aggregations,
	}
}

// GetEntityTypes 获取所有支持搜索的实体类型
func (s *SearchService) GetEntityTypes() *IndexListResponse {
	entityTypes := s.registry.ListEntityTypes()
	return &IndexListResponse{
		EntityTypes: entityTypes,
	}
}