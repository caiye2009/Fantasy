package domain

// SearchCriteria 搜索条件（领域对象）
type SearchCriteria struct {
	EntityType  string
	IndexName   string
	Query       string                 // 全文搜索关键词
	Filters     map[string]interface{} // 筛选条件
	AggRequests map[string]AggRequest  // 聚合请求
	Pagination  Pagination
	Sort        []SortField
}

// AggRequest 聚合请求
type AggRequest struct {
	Search string                 // 下拉框搜索词
	After  map[string]interface{} // 分页游标
	Size   int                    // 每页条数（可选，覆盖配置默认值）
}

// Pagination 分页
type Pagination struct {
	Offset int
	Size   int
}

// SortField 排序字段
type SortField struct {
	Field string
	Order string // asc, desc
}
