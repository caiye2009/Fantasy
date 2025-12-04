package application

import (
	"context"
	
	"back/internal/search/domain"
)

// SearchRepository 搜索仓储接口
type SearchRepository interface {
	Search(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResponse, error)
	IndexDocument(ctx context.Context, index, id string, doc interface{}) error
	UpdateDocument(ctx context.Context, index, id string, doc interface{}) error
	DeleteDocument(ctx context.Context, index, id string) error
}

// SearchService 搜索应用服务
type SearchService struct {
	repo SearchRepository
}

// NewSearchService 创建搜索服务
func NewSearchService(repo SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

// Search 执行搜索
func (s *SearchService) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	// 1. DTO → Domain Model
	query := ToSearchQuery(req)
	
	// 2. 领域验证
	if err := query.Validate(); err != nil {
		return nil, err
	}
	
	// 3. 标准化参数
	query.NormalizeSize()
	
	// 4. 执行搜索
	resp, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	
	// 5. Domain Model → DTO
	return ToSearchResponse(resp), nil
}

// GetIndices 获取所有索引
func (s *SearchService) GetIndices() *IndexListResponse {
	return ToIndexListResponse()
}

// IndexDocument 索引文档
func (s *SearchService) IndexDocument(ctx context.Context, index, id string, doc interface{}) error {
	return s.repo.IndexDocument(ctx, index, id, doc)
}

// UpdateDocument 更新文档
func (s *SearchService) UpdateDocument(ctx context.Context, index, id string, doc interface{}) error {
	return s.repo.UpdateDocument(ctx, index, id, doc)
}

// DeleteDocument 删除文档
func (s *SearchService) DeleteDocument(ctx context.Context, index, id string) error {
	return s.repo.DeleteDocument(ctx, index, id)
}