package application

// SearchRequest 搜索请求（全新结构）
type SearchRequest struct {
	EntityType  string                 `json:"entityType" binding:"required"` // material, order, client, etc.
	Query       string                 `json:"query"`                         // 全文搜索关键词
	Filters     map[string]interface{} `json:"filters"`                       // 筛选条件
	AggRequests map[string]AggRequest  `json:"aggRequests"`                   // 聚合请求
	Pagination  PaginationRequest      `json:"pagination"`
	Sort        []SortRequest          `json:"sort"`
}

// AggRequest 聚合请求
type AggRequest struct {
	Search string                 `json:"search"` // 下拉框搜索词
	After  map[string]interface{} `json:"after"`  // 分页游标
	Size   int                    `json:"size"`   // 每页条数（可选，覆盖配置默认值）
}

// PaginationRequest 分页请求
type PaginationRequest struct {
	Offset int `json:"offset"`
	Size   int `json:"size" binding:"max=100"` // 最大 100 条
}

// SortRequest 排序请求
type SortRequest struct {
	Field string `json:"field"`
	Order string `json:"order"` // asc, desc
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Items        []map[string]interface{} `json:"items"`        // 搜索结果
	Total        int                      `json:"total"`        // 总条数
	Took         int                      `json:"took"`         // 耗时 (ms)
	Aggregations map[string]AggResult     `json:"aggregations"` // 聚合结果
}

// AggResult 聚合结果
type AggResult struct {
	Buckets []Bucket               `json:"buckets,omitempty"` // terms/composite 聚合
	After   map[string]interface{} `json:"after,omitempty"`   // 下一页游标
	HasMore bool                   `json:"hasMore"`           // 是否有更多

	// stats 聚合
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
	Avg *float64 `json:"avg,omitempty"`
}

// Bucket 聚合桶
type Bucket struct {
	Key      interface{} `json:"key"`
	DocCount int         `json:"docCount"`
}

// IndexListResponse 索引列表响应
type IndexListResponse struct {
	EntityTypes []string `json:"entityTypes"`
}