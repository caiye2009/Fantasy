package fields

import (
	"strings"
	"unicode"
)

// ToSnakeCase 将 camelCase 转换为 snake_case
// 例如: supplierId -> supplier_id, createdAt -> created_at
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	result.Grow(len(s) + 5) // 预分配，通常会多几个下划线

	for i, r := range s {
		// 如果是大写字母
		if unicode.IsUpper(r) {
			// 不是第一个字符，前面加下划线
			if i > 0 {
				result.WriteRune('_')
			}
			// 转小写
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ToCamelCase 将 snake_case 转换为 camelCase（辅助函数，可选）
// 例如: supplier_id -> supplierId
func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}

	var result strings.Builder
	result.WriteString(parts[0]) // 第一个部分保持小写

	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			// 首字母大写
			result.WriteString(strings.ToUpper(parts[i][:1]))
			if len(parts[i]) > 1 {
				result.WriteString(parts[i][1:])
			}
		}
	}

	return result.String()
}

// ConvertMapKeysToSnakeCase 转换 map 的 key 从 camelCase 到 snake_case
func ConvertMapKeysToSnakeCase(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}

	result := make(map[string]interface{}, len(m))
	for key, value := range m {
		result[ToSnakeCase(key)] = value
	}
	return result
}
