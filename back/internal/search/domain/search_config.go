package domain

// SearchConfig 搜索配置
type SearchConfig struct {
	Index             string                   `yaml:"index"`
	QueryFields       []QueryFieldConfig       `yaml:"queryFields"`
	FilterFields      []FilterFieldConfig      `yaml:"filterFields"`
	AggregationFields []AggregationFieldConfig `yaml:"aggregationFields"`
	DefaultSort       []SortFieldConfig        `yaml:"defaultSort"` // 默认排序（后端自动添加）
}

// QueryFieldConfig Query字段配置
type QueryFieldConfig struct {
	Field string  `yaml:"field"`
	Boost float64 `yaml:"boost"`
}

// FilterFieldConfig Filter字段配置
type FilterFieldConfig struct {
	Field    string `yaml:"field"`
	Type     string `yaml:"type"`     // keyword, text, numeric, date
	Operator string `yaml:"operator"` // term, terms, match, range, exists
}

// AggregationFieldConfig Aggregation字段配置
type AggregationFieldConfig struct {
	Field         string `yaml:"field"`
	Type          string `yaml:"type"`          // keyword, numeric, date
	AggType       string `yaml:"aggType"`       // terms, composite, stats, date_histogram
	Size          int    `yaml:"size"`          // 每页条数
	SupportSearch bool   `yaml:"supportSearch"` // 是否支持下拉框搜索
	ExcludeSelf   bool   `yaml:"excludeSelf"`   // 联动时是否排除自身条件
}

// SortFieldConfig 排序字段配置
type SortFieldConfig struct {
	Field      string `yaml:"field"`                // 字段名
	Order      string `yaml:"order"`                // asc, desc
	Type       string `yaml:"type,omitempty"`       // computed（计算字段）, 普通字段不填
	Missing    string `yaml:"missing,omitempty"`    // _last, _first（缺失值处理）
}

// GetFilterField 获取 Filter 字段配置
func (c *SearchConfig) GetFilterField(field string) *FilterFieldConfig {
	for i := range c.FilterFields {
		if c.FilterFields[i].Field == field {
			return &c.FilterFields[i]
		}
	}
	return nil
}

// GetAggField 获取 Aggregation 字段配置
func (c *SearchConfig) GetAggField(field string) *AggregationFieldConfig {
	for i := range c.AggregationFields {
		if c.AggregationFields[i].Field == field {
			return &c.AggregationFields[i]
		}
	}
	return nil
}

// IsFilterableField 判断字段是否可过滤
func (c *SearchConfig) IsFilterableField(field string) bool {
	// 检查是否是基础字段
	for _, f := range c.FilterFields {
		if f.Field == field {
			return true
		}
		// 检查范围查询后缀（如 priceMin, priceMax）
		if f.Operator == "range" && (field == f.Field+"Min" || field == f.Field+"Max") {
			return true
		}
	}
	return false
}

// IsAggregableField 判断字段是否可聚合
func (c *SearchConfig) IsAggregableField(field string) bool {
	for _, f := range c.AggregationFields {
		if f.Field == field {
			return true
		}
	}
	return false
}

// IsQueryableField 判断字段是否可搜索
func (c *SearchConfig) IsQueryableField(field string) bool {
	for _, f := range c.QueryFields {
		if f.Field == field {
			return true
		}
	}
	return false
}

// GetQueryField 获取查询字段配置
func (c *SearchConfig) GetQueryField(field string) *QueryFieldConfig {
	for i := range c.QueryFields {
		if c.QueryFields[i].Field == field {
			return &c.QueryFields[i]
		}
	}
	return nil
}
