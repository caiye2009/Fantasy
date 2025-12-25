package casbin

import (
	"log"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

const (
	// DefaultPolicyFile 默认策略文件路径
	DefaultPolicyFile = "./config/policies.csv"
)

// InitEnforcer 初始化 Casbin enforcer（文件存储）
func InitEnforcer(policyFile string) *casbin.Enforcer {
	if policyFile == "" {
		policyFile = DefaultPolicyFile
	}

	// 确保策略文件存在
	if _, err := os.Stat(policyFile); os.IsNotExist(err) {
		log.Printf("Policy file not found: %s, will create on first save", policyFile)
		// 创建空文件
		f, err := os.Create(policyFile)
		if err != nil {
			log.Fatalf("Failed to create policy file: %v", err)
		}
		f.Close()
	}

	// 使用 File adapter
	adapter := fileadapter.NewAdapter(policyFile)

	// Casbin Model 定义
	// 支持用户个性化权限 + 角色权限
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (r.obj == p.obj || keyMatch2(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`)
	if err != nil {
		log.Fatalf("Failed to create Casbin model: %v", err)
	}

	// 创建 enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		log.Printf("Warning: Failed to load policy, starting with empty policy: %v", err)
	}

	// 启用自动保存（AddPolicy/RemovePolicy 时自动保存到文件）
	enforcer.EnableAutoSave(true)

	log.Printf("✓ Casbin enforcer initialized with file adapter: %s", policyFile)
	return enforcer
}

// InitDefaultPolicies 已废弃，默认权限现在直接配置在 policies.csv 文件中
// 启动时会自动从文件加载
func InitDefaultPolicies(enforcer *casbin.Enforcer) error {
	policies, err := enforcer.GetPolicy()
	if err != nil {
		log.Printf("Warning: Failed to get policies: %v", err)
		return err
	}
	log.Printf("✓ Loaded %d policies from file", len(policies))
	return nil
}
