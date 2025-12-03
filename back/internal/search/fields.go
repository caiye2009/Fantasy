package search

// SearchRequest 通用搜索请求
type SearchRequest struct {
	Query   string                 `json:"query"`             // 搜索关键词
	Indices []string               `json:"indices,omitempty"` // 要搜索的索引列表，空=全部索引
	Fields  []string               `json:"fields,omitempty"`  // 要搜索的字段列表，空=全部字段
	Filters map[string]interface{} `json:"filters,omitempty"` // 筛选条件 {"status": "已复核"}
	Sort    []SortField            `json:"sort,omitempty"`    // 排序
	From    int                    `json:"from"`              // 分页起始，默认0
	Size    int                    `json:"size"`              // 每页大小，默认20
}

// SortField 排序字段
type SortField struct {
	Field string `json:"field"` // 字段名
	Order string `json:"order"` // asc/desc
}

// SearchResponse 通用搜索响应
type SearchResponse struct {
	Total      int            `json:"total"`                // 总数
	Took       int            `json:"took"`                 // 搜索耗时(ms)
	MaxScore   float64        `json:"maxScore"`             // 最高分
	Results    []SearchResult `json:"results"`              // 搜索结果
	Aggregations interface{}  `json:"aggregations,omitempty"` // 聚合结果（可选）
}

// SearchResult 单条搜索结果
type SearchResult struct {
	Index     string                 `json:"index"`               // 索引名: materials/samples/products
	Type      string                 `json:"type"`                // 类型: material/sample/product
	ID        string                 `json:"id"`                  // 文档ID
	Score     float64                `json:"score"`               // 相关度评分
	Source    map[string]interface{} `json:"source"`              // 原始数据
	Highlight map[string][]string    `json:"highlight,omitempty"` // 高亮片段
}

// SuggestRequest 搜索建议请求
type SuggestRequest struct {
	Query  string   `json:"query"`             // 输入的部分关键词
	Indices []string `json:"indices,omitempty"` // 建议的索引，空=全部
	Size   int      `json:"size"`              // 建议数量
}

// SuggestResponse 搜索建议响应
type SuggestResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
}

// Suggestion 单个建议
type Suggestion struct {
	Text  string  `json:"text"`  // 建议文本
	Score float64 `json:"score"` // 评分
}