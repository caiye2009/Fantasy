package casbin

import (
	"errors"
	"log"

	"github.com/casbin/casbin/v2"
	"back/pkg/auth"
)

// Manager Casbin 权限管理器
type Manager struct {
	enforcer         *casbin.Enforcer
	whitelistManager *auth.WhitelistManager
}

// NewManager 创建权限管理器
func NewManager(enforcer *casbin.Enforcer, whitelistManager *auth.WhitelistManager) *Manager {
	return &Manager{
		enforcer:         enforcer,
		whitelistManager: whitelistManager,
	}
}

// ==================== 用户权限管理 ====================

// AddPermissionForUser 给用户添加权限
func (m *Manager) AddPermissionForUser(loginID, permission string) error {
	added, err := m.enforcer.AddPolicy(loginID, permission, "*")
	if err != nil {
		return err
	}

	if !added {
		return errors.New("权限已存在")
	}

	// 清除白名单，强制用户重新登录
	if m.whitelistManager != nil {
		if err := m.whitelistManager.RemoveAllForUser(loginID); err != nil {
			log.Printf("Warning: Failed to remove whitelist for user %s: %v", loginID, err)
		}
	}

	log.Printf("✓ Added permission for user %s: %s", loginID, permission)
	return nil
}

// RemovePermissionForUser 删除用户权限
func (m *Manager) RemovePermissionForUser(loginID, permission string) error {
	removed, err := m.enforcer.RemovePolicy(loginID, permission, "*")
	if err != nil {
		return err
	}

	if !removed {
		return errors.New("权限不存在")
	}

	// 清除白名单
	if m.whitelistManager != nil {
		if err := m.whitelistManager.RemoveAllForUser(loginID); err != nil {
			log.Printf("Warning: Failed to remove whitelist for user %s: %v", loginID, err)
		}
	}

	log.Printf("✓ Removed permission for user %s: %s", loginID, permission)
	return nil
}

// GetPermissionsForUser 获取用户的所有权限
func (m *Manager) GetPermissionsForUser(loginID string) ([]string, error) {
	permissions, err := m.enforcer.GetPermissionsForUser(loginID)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(permissions))
	for _, p := range permissions {
		if len(p) > 1 {
			result = append(result, p[1]) // p[1] 是 permission
		}
	}
	return result, nil
}

// HasPermissionForUser 检查用户是否有某个权限
func (m *Manager) HasPermissionForUser(loginID, permission string) (bool, error) {
	has, err := m.enforcer.HasPolicy(loginID, permission, "*")
	return has, err
}

// ==================== 角色权限管理 ====================

// AddPermissionForRole 给角色添加权限
func (m *Manager) AddPermissionForRole(role, permission string) error {
	added, err := m.enforcer.AddPolicy(role, permission, "*")
	if err != nil {
		return err
	}

	if !added {
		return errors.New("权限已存在")
	}

	log.Printf("✓ Added permission for role %s: %s", role, permission)
	return nil
}

// RemovePermissionForRole 删除角色权限
func (m *Manager) RemovePermissionForRole(role, permission string) error {
	removed, err := m.enforcer.RemovePolicy(role, permission, "*")
	if err != nil {
		return err
	}

	if !removed {
		return errors.New("权限不存在")
	}

	log.Printf("✓ Removed permission for role %s: %s", role, permission)
	return nil
}

// GetPermissionsForRole 获取角色的所有权限
func (m *Manager) GetPermissionsForRole(role string) ([]string, error) {
	permissions, err := m.enforcer.GetPermissionsForUser(role)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(permissions))
	for _, p := range permissions {
		if len(p) > 1 {
			result = append(result, p[1])
		}
	}
	return result, nil
}

// ==================== 权限检查 ====================

// Enforce 检查权限（两层检查：用户权限 + 角色权限）
func (m *Manager) Enforce(loginID, role, permission string) bool {
	// 1. 先检查用户个性化权限
	allowed, _ := m.enforcer.Enforce(loginID, permission, "*")
	if allowed {
		return true
	}

	// 2. 检查角色权限
	allowed, _ = m.enforcer.Enforce(role, permission, "*")
	return allowed
}

// ==================== 批量操作 ====================

// RemoveAllPermissionsForUser 删除用户的所有个性化权限
func (m *Manager) RemoveAllPermissionsForUser(loginID string) error {
	_, err := m.enforcer.RemoveFilteredPolicy(0, loginID)
	if err != nil {
		return err
	}

	// 清除白名单
	if m.whitelistManager != nil {
		if err := m.whitelistManager.RemoveAllForUser(loginID); err != nil {
			log.Printf("Warning: Failed to remove whitelist for user %s: %v", loginID, err)
		}
	}

	log.Printf("✓ Removed all permissions for user %s", loginID)
	return nil
}

// ==================== 策略管理 ====================

// GetAllPolicies 获取所有策略
func (m *Manager) GetAllPolicies() ([][]string, error) {
	policies, err := m.enforcer.GetPolicy()
	if err != nil {
		return nil, err
	}
	return policies, nil
}

// ReloadPolicies 重新加载策略（从文件）
func (m *Manager) ReloadPolicies() error {
	return m.enforcer.LoadPolicy()
}

// SavePolicies 手动保存策略到文件
func (m *Manager) SavePolicies() error {
	return m.enforcer.SavePolicy()
}
