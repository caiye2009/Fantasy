package config

import (
	"back/pkg/auth"
	casbinPkg "back/pkg/casbin"
	"github.com/casbin/casbin/v2"
)

// InitCasbin 初始化 Casbin enforcer（文件存储）
func InitCasbin(cfg *Config) *casbin.Enforcer {
	// 使用 pkg/casbin 的初始化函数
	enforcer := casbinPkg.InitEnforcer(casbinPkg.DefaultPolicyFile)

	// 初始化默认策略（如果策略文件为空）
	if err := casbinPkg.InitDefaultPolicies(enforcer); err != nil {
		panic(err)
	}

	return enforcer
}

// InitCasbinManager 初始化 Casbin 管理器
func InitCasbinManager(enforcer *casbin.Enforcer, whitelistManager *auth.WhitelistManager) *casbinPkg.Manager {
	return casbinPkg.NewManager(enforcer, whitelistManager)
}
