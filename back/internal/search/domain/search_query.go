package domain

// SearchQuery 搜索查询值对象
type SearchQuery struct {
	Keyword string
	Indices []string
	Fields  []string
	Filters map[string]interface{}
	Sort    []SortField
	From    int
	Size    int
}

// SortField 排序字段
type SortField struct {
	Field string
	Order string // asc/desc
}

// Validate 验证搜索查询
func (q *SearchQuery) Validate() error {
	if q.Size < 0 {
		return ErrInvalidSize
	}
	
	if q.Size > 100 {
		return ErrSizeTooLarge
	}
	
	if q.From < 0 {
		return ErrInvalidFrom
	}
	
	return nil
}

// NormalizeSize 标准化分页大小
func (q *SearchQuery) NormalizeSize() {
	if q.Size == 0 {
		q.Size = 20
	}
	if q.Size > 100 {
		q.Size = 100
	}
}

// GetIndices 获取搜索索引
func (q *SearchQuery) GetIndices() []string {
	if len(q.Indices) == 0 {
		return GetAllIndices()
	}
	return q.Indices
}

// GetFields 获取搜索字段
func (q *SearchQuery) GetFields() []string {
	if len(q.Fields) == 0 {
		return GetDefaultFields(q.GetIndices())
	}
	return q.Fields
}