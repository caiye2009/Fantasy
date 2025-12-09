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

	// 验证过滤字段
	if err := q.ValidateFilterFields(); err != nil {
		return err
	}

	// 验证排序字段
	if err := q.ValidateSortFields(); err != nil {
		return err
	}

	return nil
}

// ValidateFilterFields 验证过滤字段是否在白名单内
func (q *SearchQuery) ValidateFilterFields() error {
	if len(q.Filters) == 0 {
		return nil
	}

	// 获取所有索引
	indices := q.GetIndices()

	// 收集所有允许的过滤字段
	allowedFields := make(map[string]bool)
	for _, index := range indices {
		meta, ok := IndexConfig[index]
		if !ok {
			continue
		}
		for _, field := range meta.FilterableFields {
			allowedFields[field] = true
		}
	}

	// 验证每个过滤字段
	for field := range q.Filters {
		if !allowedFields[field] {
			// 返回第一个索引名称用于错误信息
			indexName := "unknown"
			if len(indices) > 0 {
				indexName = indices[0]
			}
			return ErrInvalidFilterField(field, indexName)
		}
	}

	return nil
}

// ValidateSortFields 验证排序字段是否在白名单内
func (q *SearchQuery) ValidateSortFields() error {
	if len(q.Sort) == 0 {
		return nil
	}

	// 获取所有索引
	indices := q.GetIndices()

	// 收集所有允许的排序字段
	allowedFields := make(map[string]bool)
	for _, index := range indices {
		meta, ok := IndexConfig[index]
		if !ok {
			continue
		}
		for _, field := range meta.SortableFields {
			allowedFields[field] = true
		}
	}

	// 验证每个排序字段
	for _, sortField := range q.Sort {
		if !allowedFields[sortField.Field] {
			// 返回第一个索引名称用于错误信息
			indexName := "unknown"
			if len(indices) > 0 {
				indexName = indices[0]
			}
			return ErrInvalidSortField(sortField.Field, indexName)
		}
	}

	return nil
}

// NormalizeSize 标准化分页大小
func (q *SearchQuery) NormalizeSize() {
	if q.Size <= 0 {
		q.Size = 10 // 默认每页 10 条
	}
	if q.Size > 100 {
		q.Size = 100 // 最大每页 100 条
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