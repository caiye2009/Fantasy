package es

import (
	"fmt"
	"reflect"
	"strings"
)

// FieldInfo 字段信息（从 Domain 提取）
type FieldInfo struct {
	Name       string // JSON 字段名（camelCase）
	GoName     string // Go 字段名（PascalCase）
	GoType     string // Go 类型
	ESType     string // ES 类型
	IsKeyword  bool   // 是否是 keyword 类型
	IsText     bool   // 是否是 text 类型
	IsNumeric  bool   // 是否是数值类型
	IsDate     bool   // 是否是日期类型
}

// DomainSchema Domain 模型的字段 schema
type DomainSchema struct {
	EntityType string
	IndexName  string
	Fields     map[string]*FieldInfo // fieldName -> FieldInfo
}

// ExtractSchema 从 Domain 模型提取字段 schema
// 示例: ExtractSchema("client", "clients", &domain.Client{})
func ExtractSchema(entityType, indexName string, domainModel interface{}) (*DomainSchema, error) {
	schema := &DomainSchema{
		EntityType: entityType,
		IndexName:  indexName,
		Fields:     make(map[string]*FieldInfo),
	}

	t := reflect.TypeOf(domainModel)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("domain model must be a struct")
	}

	// 遍历所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 跳过未导出的字段
		if !field.IsExported() {
			continue
		}

		// 跳过没有 json tag 的字段（如 DeletedAt 带 json:"-"）
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// 解析 json tag 获取字段名
		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		// 推断 ES 类型
		esType, isKeyword, isText, isNumeric, isDate := inferESType(field.Type)

		fieldInfo := &FieldInfo{
			Name:      fieldName,
			GoName:    field.Name,
			GoType:    field.Type.String(),
			ESType:    esType,
			IsKeyword: isKeyword,
			IsText:    isText,
			IsNumeric: isNumeric,
			IsDate:    isDate,
		}

		schema.Fields[fieldName] = fieldInfo
	}

	return schema, nil
}

// inferESType 根据 Go 类型推断 ES 类型
func inferESType(t reflect.Type) (esType string, isKeyword, isText, isNumeric, isDate bool) {
	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	typeName := t.String()

	switch t.Kind() {
	case reflect.String:
		// 字符串字段
		// 规则：大部分字符串用 text（支持全文搜索）+ keyword 子字段（支持精确匹配、排序、聚合）
		return "text", false, true, false, false

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "long", false, false, true, false

	case reflect.Float32, reflect.Float64:
		return "double", false, false, true, false

	case reflect.Bool:
		return "boolean", false, false, false, false

	case reflect.Struct:
		// 处理 time.Time
		if typeName == "time.Time" {
			return "date", false, false, false, true
		}
		// 其他结构体（如 gorm.DeletedAt）忽略
		return "", false, false, false, false

	default:
		return "text", false, true, false, false
	}
}

// GetField 获取字段信息
func (s *DomainSchema) GetField(fieldName string) (*FieldInfo, bool) {
	field, ok := s.Fields[fieldName]
	return field, ok
}

// ValidateField 验证字段是否存在于 Domain
func (s *DomainSchema) ValidateField(fieldName string) error {
	if _, ok := s.Fields[fieldName]; !ok {
		return fmt.Errorf("field '%s' not found in domain model %s", fieldName, s.EntityType)
	}
	return nil
}

// GetFieldType 获取字段的 ES 类型
func (s *DomainSchema) GetFieldType(fieldName string) (string, error) {
	field, ok := s.Fields[fieldName]
	if !ok {
		return "", fmt.Errorf("field '%s' not found", fieldName)
	}
	return field.ESType, nil
}

// ListFields 列出所有字段名
func (s *DomainSchema) ListFields() []string {
	fields := make([]string, 0, len(s.Fields))
	for name := range s.Fields {
		fields = append(fields, name)
	}
	return fields
}

// ListKeywordFields 列出所有 keyword 类型字段（用于精确匹配、聚合）
func (s *DomainSchema) ListKeywordFields() []string {
	fields := make([]string, 0)
	for name, info := range s.Fields {
		if info.IsKeyword {
			fields = append(fields, name)
		}
	}
	return fields
}

// ListTextFields 列出所有 text 类型字段（用于全文搜索）
func (s *DomainSchema) ListTextFields() []string {
	fields := make([]string, 0)
	for name, info := range s.Fields {
		if info.IsText {
			fields = append(fields, name)
		}
	}
	return fields
}

// ListNumericFields 列出所有数值类型字段
func (s *DomainSchema) ListNumericFields() []string {
	fields := make([]string, 0)
	for name, info := range s.Fields {
		if info.IsNumeric {
			fields = append(fields, name)
		}
	}
	return fields
}

// ListDateFields 列出所有日期类型字段
func (s *DomainSchema) ListDateFields() []string {
	fields := make([]string, 0)
	for name, info := range s.Fields {
		if info.IsDate {
			fields = append(fields, name)
		}
	}
	return fields
}
