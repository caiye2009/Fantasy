package domain

// IndexMeta 索引元数据
type IndexMeta struct {
	Name          string   // 索引名称
	Type          string   // 类型名称
	DefaultFields []string // 默认搜索字段（带权重）
}

// IndexConfig 索引配置映射
var IndexConfig = map[string]IndexMeta{
	"vendors": {
		Name: "vendors",
		Type: "vendor",
		DefaultFields: []string{
			"name^5",
			"contact^3",
			"phone^2",
			"email^2",
			"address",
		},
	},
	"clients": {
		Name: "clients",
		Type: "client",
		DefaultFields: []string{
			"name^5",
			"contact^3",
			"phone^2",
			"email^2",
			"address",
		},
	},
	"materials": {
		Name: "materials",
		Type: "material",
		DefaultFields: []string{
			"name^5",
			"spec^3",
			"unit^2",
			"description",
		},
	},
	"processes": {
		Name: "processes",
		Type: "process",
		DefaultFields: []string{
			"name^5",
			"description^3",
		},
	},
	"products": {
		Name: "products",
		Type: "product",
		DefaultFields: []string{
			"name^5",
			"status^2",
		},
	},
	"orders": {
		Name: "orders",
		Type: "order",
		DefaultFields: []string{
			"order_no^5",
			"status^2",
		},
	},
	"plans": {
		Name: "plans",
		Type: "plan",
		DefaultFields: []string{
			"plan_no^5",
			"status^2",
		},
	},
}

// GetAllIndices 获取所有索引名称
func GetAllIndices() []string {
	indices := make([]string, 0, len(IndexConfig))
	for name := range IndexConfig {
		indices = append(indices, name)
	}
	return indices
}

// GetDefaultFields 获取索引的默认搜索字段
func GetDefaultFields(indices []string) []string {
	if len(indices) == 0 {
		return []string{"*"}
	}

	fieldsMap := make(map[string]bool)
	for _, idx := range indices {
		if meta, ok := IndexConfig[idx]; ok {
			for _, field := range meta.DefaultFields {
				fieldsMap[field] = true
			}
		}
	}

	fields := make([]string, 0, len(fieldsMap))
	for field := range fieldsMap {
		fields = append(fields, field)
	}

	if len(fields) == 0 {
		return []string{"*"}
	}

	return fields
}

// GetIndexMeta 获取索引元数据
func GetIndexMeta(index string) (IndexMeta, bool) {
	meta, ok := IndexConfig[index]
	return meta, ok
}