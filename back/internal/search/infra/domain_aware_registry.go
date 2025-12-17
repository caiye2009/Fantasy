package infra

import (
	"fmt"
	"log"

	"back/internal/search/domain"
	"back/pkg/es"
)

// DomainAwareRegistry Domain 感知的搜索配置注册中心
// 以 Domain 为唯一字段来源，自动验证和补全配置
type DomainAwareRegistry struct {
	schemas map[string]*es.DomainSchema     // indexName -> schema
	configs map[string]*domain.SearchConfig // indexName -> config
}

// NewDomainAwareRegistry 创建 Domain 感知的注册中心
func NewDomainAwareRegistry() *DomainAwareRegistry {
	return &DomainAwareRegistry{
		schemas: make(map[string]*es.DomainSchema),
		configs: make(map[string]*domain.SearchConfig),
	}
}

// RegisterDomain 注册 Domain 模型
// 示例: RegisterDomain("client", "clients", &clientDomain.Client{})
func (r *DomainAwareRegistry) RegisterDomain(entityType, indexName string, domainModel interface{}) error {
	// 提取 Domain schema
	schema, err := es.ExtractSchema(entityType, indexName, domainModel)
	if err != nil {
		return fmt.Errorf("extract schema for %s: %w", entityType, err)
	}

	r.schemas[indexName] = schema
	log.Printf("✓ Registered domain schema for index '%s' (%s) with %d fields", indexName, entityType, len(schema.Fields))

	return nil
}

// LoadConfig 加载配置并自动验证、补全
func (r *DomainAwareRegistry) LoadConfig(config *domain.SearchConfig) error {
	indexName := config.IndexName

	// 1. 检查是否已注册 Domain
	schema, ok := r.schemas[indexName]
	if !ok {
		return fmt.Errorf("domain not registered for index: %s", indexName)
	}

	// 2. 验证所有配置的字段都存在于 Domain
	if err := r.validateQueryFields(schema, config); err != nil {
		return fmt.Errorf("validate query fields: %w", err)
	}
	if err := r.validateFilterFields(schema, config); err != nil {
		return fmt.Errorf("validate filter fields: %w", err)
	}
	if err := r.validateAggregationFields(schema, config); err != nil {
		return fmt.Errorf("validate aggregation fields: %w", err)
	}
	if err := r.validateDefaultSort(schema, config); err != nil {
		return fmt.Errorf("validate default sort: %w", err)
	}

	// 3. 自动补全字段类型（从 Domain schema）
	if err := r.enrichFilterFields(schema, config); err != nil {
		return fmt.Errorf("enrich filter fields: %w", err)
	}
	if err := r.enrichAggregationFields(schema, config); err != nil {
		return fmt.Errorf("enrich aggregation fields: %w", err)
	}

	// 4. 保存配置
	r.configs[indexName] = config
	log.Printf("✓ Loaded search config for index '%s'", indexName)

	return nil
}

// validateQueryFields 验证 query 字段
func (r *DomainAwareRegistry) validateQueryFields(schema *es.DomainSchema, config *domain.SearchConfig) error {
	for i, qf := range config.QueryFields {
		if err := schema.ValidateField(qf.Field); err != nil {
			return fmt.Errorf("queryFields[%d]: %w", i, err)
		}
	}
	return nil
}

// validateFilterFields 验证 filter 字段
func (r *DomainAwareRegistry) validateFilterFields(schema *es.DomainSchema, config *domain.SearchConfig) error {
	for i, ff := range config.FilterFields {
		if err := schema.ValidateField(ff.Field); err != nil {
			return fmt.Errorf("filterFields[%d]: %w", i, err)
		}
	}
	return nil
}

// validateAggregationFields 验证 aggregation 字段
func (r *DomainAwareRegistry) validateAggregationFields(schema *es.DomainSchema, config *domain.SearchConfig) error {
	for i, af := range config.AggregationFields {
		if err := schema.ValidateField(af.Field); err != nil {
			return fmt.Errorf("aggregationFields[%d]: %w", i, err)
		}
	}
	return nil
}

// validateDefaultSort 验证默认排序字段
func (r *DomainAwareRegistry) validateDefaultSort(schema *es.DomainSchema, config *domain.SearchConfig) error {
	for i, sf := range config.DefaultSort {
		// 如果是计算字段（type=computed），跳过验证
		if sf.Type == "computed" {
			log.Printf("  - defaultSort[%d]: %s (computed field, skip validation)", i, sf.Field)
			continue
		}
		// 普通字段必须存在于 Domain
		if err := schema.ValidateField(sf.Field); err != nil {
			return fmt.Errorf("defaultSort[%d]: %w", i, err)
		}
	}
	return nil
}

// enrichFilterFields 从 Domain 自动推断并补全 filter 字段类型
func (r *DomainAwareRegistry) enrichFilterFields(schema *es.DomainSchema, config *domain.SearchConfig) error {
	for i := range config.FilterFields {
		ff := &config.FilterFields[i]

		// 如果配置中没有指定 type，从 Domain 自动推断
		if ff.Type == "" {
			fieldInfo, _ := schema.GetField(ff.Field)
			ff.Type = mapESTypeToFilterType(fieldInfo.ESType)
		}
	}
	return nil
}

// enrichAggregationFields 从 Domain 自动推断并补全 aggregation 字段类型
func (r *DomainAwareRegistry) enrichAggregationFields(schema *es.DomainSchema, config *domain.SearchConfig) error {
	for i := range config.AggregationFields {
		af := &config.AggregationFields[i]

		// 如果配置中没有指定 type，从 Domain 自动推断
		if af.Type == "" {
			fieldInfo, _ := schema.GetField(af.Field)
			af.Type = mapESTypeToAggregationType(fieldInfo.ESType)
		}
	}
	return nil
}

// mapESTypeToFilterType 映射 ES 类型到 filter 类型
func mapESTypeToFilterType(esType string) string {
	switch esType {
	case "text":
		return "text"
	case "keyword":
		return "keyword"
	case "long", "integer", "short", "byte", "double", "float":
		return "numeric"
	case "date":
		return "date"
	case "boolean":
		return "keyword"
	default:
		return "keyword"
	}
}

// mapESTypeToAggregationType 映射 ES 类型到 aggregation 类型
func mapESTypeToAggregationType(esType string) string {
	switch esType {
	case "text", "keyword":
		return "keyword"
	case "long", "integer", "short", "byte", "double", "float":
		return "numeric"
	case "date":
		return "date"
	case "boolean":
		return "keyword"
	default:
		return "keyword"
	}
}

// GetConfigByIndex 通过索引名获取配置
func (r *DomainAwareRegistry) GetConfigByIndex(indexName string) (*domain.SearchConfig, bool) {
	config, ok := r.configs[indexName]
	return config, ok
}

// GetSchema 获取 Domain schema
func (r *DomainAwareRegistry) GetSchema(indexName string) (*es.DomainSchema, bool) {
	schema, ok := r.schemas[indexName]
	return schema, ok
}

// ListIndices 列出所有索引
func (r *DomainAwareRegistry) ListIndices() []string {
	indices := make([]string, 0, len(r.configs))
	for idx := range r.configs {
		indices = append(indices, idx)
	}
	return indices
}
