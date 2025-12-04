package application

// SearchRequest 搜索请求
type SearchRequest struct {
	Query   string                 `json:"query"`
	Indices []string               `json:"indices,omitempty"`
	Fields  []string               `json:"fields,omitempty"`
	Filters map[string]interface{} `json:"filters,omitempty"`
	Sort    []SortFieldDTO         `json:"sort,omitempty"`
	From    int                    `json:"from"`
	Size    int                    `json:"size"`
}

// SortFieldDTO 排序字段 DTO
type SortFieldDTO struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Total    int                    `json:"total"`
	Took     int                    `json:"took"`
	MaxScore float64                `json:"max_score"`
	Results  []*SearchResultDTO     `json:"results"`
}

// SearchResultDTO 搜索结果 DTO
type SearchResultDTO struct {
	Index     string                 `json:"index"`
	Type      string                 `json:"type"`
	ID        string                 `json:"id"`
	Score     float64                `json:"score"`
	Source    map[string]interface{} `json:"source"`
	Highlight map[string][]string    `json:"highlight,omitempty"`
}

// IndexListResponse 索引列表响应
type IndexListResponse struct {
	Indices []*IndexInfoDTO `json:"indices"`
}

// IndexInfoDTO 索引信息 DTO
type IndexInfoDTO struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}