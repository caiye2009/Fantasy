package config

import (
	"path/filepath"

	searchInfra "back/internal/search/infra"
)

// InitSearchRegistry 初始化搜索配置注册中心
func InitSearchRegistry() (*searchInfra.SearchConfigRegistry, error) {
	// 配置文件目录（相对于项目根目录）
	configDir := filepath.Join("config", "search")

	// 加载所有 YAML 配置
	registry, err := searchInfra.NewSearchConfigRegistry(configDir)
	if err != nil {
		return nil, err
	}

	return registry, nil
}
