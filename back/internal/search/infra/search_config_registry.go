package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"back/internal/search/domain"
)

// SearchConfigRegistry 搜索配置注册中心
type SearchConfigRegistry struct {
	configs map[string]*domain.SearchConfig // entityType -> config
}

// NewSearchConfigRegistry 创建配置注册中心
func NewSearchConfigRegistry(configDir string) (*SearchConfigRegistry, error) {
	registry := &SearchConfigRegistry{
		configs: make(map[string]*domain.SearchConfig),
	}

	if err := registry.loadConfigs(configDir); err != nil {
		return nil, err
	}

	return registry, nil
}

// loadConfigs 加载配置目录下的所有 YAML 文件
func (r *SearchConfigRegistry) loadConfigs(configDir string) error {
	// 读取目录下所有文件
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("read config dir: %w", err)
	}

	// 遍历加载每个 YAML 文件
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 只处理 .yaml 和 .yml 文件
		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		filePath := filepath.Join(configDir, entry.Name())
		if err := r.loadConfigFile(filePath); err != nil {
			return fmt.Errorf("load config file %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// loadConfigFile 加载单个配置文件
func (r *SearchConfigRegistry) loadConfigFile(filePath string) error {
	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	// 解析 YAML
	var config domain.SearchConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	// 验证配置
	if err := r.validateConfig(&config); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	// 注册配置
	r.configs[config.EntityType] = &config

	return nil
}

// validateConfig 验证配置合法性
func (r *SearchConfigRegistry) validateConfig(config *domain.SearchConfig) error {
	if config.EntityType == "" {
		return fmt.Errorf("entityType is required")
	}
	if config.IndexName == "" {
		return fmt.Errorf("indexName is required")
	}

	// 验证 Query 字段配置
	for i, qf := range config.QueryFields {
		if qf.Field == "" {
			return fmt.Errorf("queryFields[%d].field is required", i)
		}
	}

	// 验证 Filter 字段配置
	for i, ff := range config.FilterFields {
		if ff.Field == "" {
			return fmt.Errorf("filterFields[%d].field is required", i)
		}
		if ff.Type == "" {
			return fmt.Errorf("filterFields[%d].type is required", i)
		}
		if ff.Operator == "" {
			return fmt.Errorf("filterFields[%d].operator is required", i)
		}

		// 验证 operator 合法性
		validOperators := map[string]bool{
			"term":   true,
			"terms":  true,
			"match":  true,
			"range":  true,
			"exists": true,
		}
		if !validOperators[ff.Operator] {
			return fmt.Errorf("filterFields[%d].operator '%s' is invalid", i, ff.Operator)
		}
	}

	// 验证 Aggregation 字段配置
	for i, af := range config.AggregationFields {
		if af.Field == "" {
			return fmt.Errorf("aggregationFields[%d].field is required", i)
		}
		if af.Type == "" {
			return fmt.Errorf("aggregationFields[%d].type is required", i)
		}
		if af.AggType == "" {
			return fmt.Errorf("aggregationFields[%d].aggType is required", i)
		}

		// 验证 aggType 合法性
		validAggTypes := map[string]bool{
			"terms":           true,
			"composite":       true,
			"stats":           true,
			"date_histogram":  true,
		}
		if !validAggTypes[af.AggType] {
			return fmt.Errorf("aggregationFields[%d].aggType '%s' is invalid", i, af.AggType)
		}
	}

	return nil
}

// GetConfig 获取指定实体的配置
func (r *SearchConfigRegistry) GetConfig(entityType string) (*domain.SearchConfig, bool) {
	config, ok := r.configs[entityType]
	return config, ok
}

// ListEntityTypes 列出所有支持搜索的实体类型
func (r *SearchConfigRegistry) ListEntityTypes() []string {
	entityTypes := make([]string, 0, len(r.configs))
	for entityType := range r.configs {
		entityTypes = append(entityTypes, entityType)
	}
	return entityTypes
}
