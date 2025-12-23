package config

import (
	"log"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// InitCasbin 初始化Casbin enforcer
func InitCasbin(cfg *Config, db *gorm.DB) *casbin.Enforcer {
	// 使用GORM adapter（策略存储在数据库）
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Failed to create Casbin adapter: %v", err)
	}

	// Casbin model定义（基于RBAC）
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`)
	if err != nil {
		log.Fatalf("Failed to create Casbin model: %v", err)
	}

	// 创建enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		log.Fatalf("Failed to load Casbin policy: %v", err)
	}

	log.Println("✓ Casbin enforcer initialized")
	return enforcer
}

func InitCasbinPolicies(enforcer *casbin.Enforcer) error {
	policies, err := enforcer.GetPolicy()
	if err != nil {
		return err
	}
	
	if len(policies) > 0 {
		log.Println("✓ Casbin policies already initialized")
		return nil
	}

	log.Println("Initializing Casbin policies...")

	defaultPolicies := [][]string{
		// Admin (全部权限)
		{"admin", "/api/v1/*", "*"},

		// HR (用户管理)
		{"hr", "/api/v1/user/create", "POST"},
		{"hr", "/api/v1/user/list", "GET"},
		{"hr", "/api/v1/user/:id", "GET"},
		{"hr", "/api/v1/user/:id", "PUT"},

		// Sales (业务员)
		{"sales", "/api/v1/vendor/list", "GET"},
		{"sales", "/api/v1/client/list", "GET"},
		{"sales", "/api/v1/client/create", "POST"},
		{"sales", "/api/v1/material/list", "GET"},
		{"sales", "/api/v1/process/list", "GET"},
		{"sales", "/api/v1/product/list", "GET"},
		{"sales", "/api/v1/product/create", "POST"},
		{"sales", "/api/v1/product/cost", "POST"},

		// Follower (跟单员)
		{"follower", "/api/v1/vendor/list", "GET"},
		{"follower", "/api/v1/client/list", "GET"},
		{"follower", "/api/v1/material/list", "GET"},
		{"follower", "/api/v1/process/list", "GET"},
		{"follower", "/api/v1/product/list", "GET"},
		{"follower", "/api/v1/product/cost", "POST"},

		// Assistant (助理,只读)
		{"assistant", "/api/v1/vendor/list", "GET"},
		{"assistant", "/api/v1/client/list", "GET"},
		{"assistant", "/api/v1/material/list", "GET"},
		{"assistant", "/api/v1/process/list", "GET"},
		{"assistant", "/api/v1/product/list", "GET"},

		// ==================== 订单模块权限 ====================
		// SalesManager - 创建订单、查看订单
		{"salesManager", "/api/v1/order", "POST"},
		{"salesManager", "/api/v1/order", "GET"},
		{"salesManager", "/api/v1/order/:id/detail", "GET"},
		{"salesManager", "/api/v1/order/:id/events", "GET"},

		// SalesAssistant - 创建订单、查看订单
		{"salesAssistant", "/api/v1/order", "POST"},
		{"salesAssistant", "/api/v1/order", "GET"},
		{"salesAssistant", "/api/v1/order/:id/detail", "GET"},
		{"salesAssistant", "/api/v1/order/:id/events", "GET"},

		// ProductionDirector - 分配部门
		{"productionDirector", "/api/v1/order", "GET"},
		{"productionDirector", "/api/v1/order/:id/detail", "GET"},
		{"productionDirector", "/api/v1/order/:id/assign-department", "POST"},
		{"productionDirector", "/api/v1/order/:id/events", "GET"},

		// ProductionAssistant - 分配人员
		{"productionAssistant", "/api/v1/order", "GET"},
		{"productionAssistant", "/api/v1/order/:id/detail", "GET"},
		{"productionAssistant", "/api/v1/order/:id/assign-personnel", "POST"},
		{"productionAssistant", "/api/v1/order/:id/events", "GET"},

		// OrderCoordinator - 更新胚布投入、生产进度、回修进度
		{"orderCoordinator", "/api/v1/order", "GET"},
		{"orderCoordinator", "/api/v1/order/:id/detail", "GET"},
		{"orderCoordinator", "/api/v1/order/:id/progress/fabric-input", "POST"},
		{"orderCoordinator", "/api/v1/order/:id/progress/production", "POST"},
		{"orderCoordinator", "/api/v1/order/:id/progress/rework", "POST"},
		{"orderCoordinator", "/api/v1/order/:id/events", "GET"},

		// Warehouse - 更新验货进度、录入次品
		{"warehouse", "/api/v1/order", "GET"},
		{"warehouse", "/api/v1/order/:id/detail", "GET"},
		{"warehouse", "/api/v1/order/:id/progress/warehouse-check", "POST"},
		{"warehouse", "/api/v1/order/:id/defect", "POST"},
		{"warehouse", "/api/v1/order/:id/events", "GET"},
	}

	for _, policy := range defaultPolicies {
		_, err := enforcer.AddPolicy(policy)
		if err != nil {
			return err
		}
	}

	log.Printf("✓ Initialized %d Casbin policies", len(defaultPolicies))
	return nil
}