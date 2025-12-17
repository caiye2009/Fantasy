package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	clientDomain "back/internal/client/domain"
	materialDomain "back/internal/material/domain"
	orderDomain "back/internal/order/domain"
	planDomain "back/internal/plan/domain"
	processDomain "back/internal/process/domain"
	productDomain "back/internal/product/domain"
	supplierDomain "back/internal/supplier/domain"

	searchDomain "back/internal/search/domain"
	"back/internal/search/infra"
)

// InitSearchRegistry 初始化搜索注册中心（Domain 感知）
func InitSearchRegistry() (*infra.DomainAwareRegistry, error) {
	log.Println("=== Initializing Search Registry (Domain-Aware) ===")

	registry := infra.NewDomainAwareRegistry()

	// 1. 注册所有 Domain 模型（字段定义的唯一来源）
	if err := registerAllDomains(registry); err != nil {
		return nil, fmt.Errorf("register domains: %w", err)
	}

	// 2. 加载配置文件（自动验证字段、补全类型）
	configDir := "./config/search"
	if err := loadAllConfigs(registry, configDir); err != nil {
		return nil, fmt.Errorf("load configs: %w", err)
	}

	log.Printf("✓ Search registry initialized successfully")
	return registry, nil
}

// registerAllDomains 注册所有 Domain 模型
func registerAllDomains(registry *infra.DomainAwareRegistry) error {
	domains := []struct {
		index string
		model interface{}
	}{
		{"client", &clientDomain.Client{}},
		{"supplier", &supplierDomain.Supplier{}},
		{"material", &materialDomain.Material{}},
		{"product", &productDomain.Product{}},
		{"process", &processDomain.Process{}},
		{"order", &orderDomain.Order{}},
		{"plan", &planDomain.Plan{}},
	}

	for _, d := range domains {
		if err := registry.RegisterDomain(d.index, d.model); err != nil {
			return fmt.Errorf("register %s: %w", d.index, err)
		}
	}

	return nil
}

// loadAllConfigs 加载所有配置文件
func loadAllConfigs(registry *infra.DomainAwareRegistry, configDir string) error {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("read config dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 只处理 .yaml 和 .yml 文件
		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		filePath := filepath.Join(configDir, entry.Name())
		if err := loadConfigFile(registry, filePath); err != nil {
			log.Printf("Warning: failed to load config file %s: %v", entry.Name(), err)
			continue
		}
	}

	return nil
}

// loadConfigFile 加载单个配置文件
func loadConfigFile(registry *infra.DomainAwareRegistry, filePath string) error {
	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	// 解析 YAML
	var config searchDomain.SearchConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	// 加载到注册中心（会自动验证和补全）
	if err := registry.LoadConfig(&config); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}
