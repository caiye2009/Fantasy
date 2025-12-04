package domain

// SearchResult 搜索结果值对象
type SearchResult struct {
	Index     string
	Type      string
	ID        string
	Score     float64
	Source    map[string]interface{}
	Highlight map[string][]string
}

// SearchResponse 搜索响应聚合根
type SearchResponse struct {
	Total    int
	Took     int
	MaxScore float64
	Results  []*SearchResult
}

// IsEmpty 是否为空结果
func (r *SearchResponse) IsEmpty() bool {
	return r.Total == 0 || len(r.Results) == 0
}

// GetResultCount 获取结果数量
func (r *SearchResponse) GetResultCount() int {
	return len(r.Results)
}