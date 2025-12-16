package domain

// SearchResult 搜索结果（领域对象）
type SearchResult struct {
	Items        []map[string]interface{} // 搜索结果项
	Total        int                      // 总条数
	Took         int                      // 耗时 (ms)
	Aggregations map[string]AggResult     // 聚合结果
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