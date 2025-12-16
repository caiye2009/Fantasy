package domain

import "context"

// SearchRepository 搜索仓储接口
type SearchRepository interface {
	// Search 执行搜索
	Search(ctx context.Context, criteria *SearchCriteria, config *SearchConfig) (*SearchResult, error)
}
