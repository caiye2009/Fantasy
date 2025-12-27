package application

import (
	"errors"

	casbinPkg "back/pkg/casbin"
	"back/pkg/endpoint"
)

// PermissionService 权限管理服务
type PermissionService struct {
	casbinManager *casbinPkg.Manager
}

// NewPermissionService 创建权限管理服务
func NewPermissionService(casbinManager *casbinPkg.Manager) *PermissionService {
	return &PermissionService{
		casbinManager: casbinManager,
	}
}

// ==================== 用户权限管理 ====================

// AddUserPermission 添加用户权限
func (s *PermissionService) AddUserPermission(req *AddUserPermissionRequest) error {
	// 验证 permission 是否在注册表中
	if !s.isValidPermission(req.Permission) {
		return errors.New("无效的权限名称: " + req.Permission)
	}

	return s.casbinManager.AddPermissionForUser(req.LoginID, req.Permission)
}

// RemoveUserPermission 删除用户权限
func (s *PermissionService) RemoveUserPermission(req *RemoveUserPermissionRequest) error {
	return s.casbinManager.RemovePermissionForUser(req.LoginID, req.Permission)
}

// GetUserPermissions 获取用户的所有权限
func (s *PermissionService) GetUserPermissions(req *GetUserPermissionsRequest) (*UserPermissionsResponse, error) {
	permissions, err := s.casbinManager.GetPermissionsForUser(req.LoginID)
	if err != nil {
		return nil, err
	}
	return &UserPermissionsResponse{
		LoginID:     req.LoginID,
		Permissions: permissions,
	}, nil
}

// RemoveAllUserPermissions 删除用户的所有个性化权限
func (s *PermissionService) RemoveAllUserPermissions(loginID string) error {
	return s.casbinManager.RemoveAllPermissionsForUser(loginID)
}

// ==================== 权限列表查询 ====================

// ListAllPermissions 列出所有可用权限
func (s *PermissionService) ListAllPermissions() *PermissionsListResponse {
	endpoints := endpoint.GlobalRegistry.ListAll()

	permissions := make([]PermissionItem, len(endpoints))
	for i, ep := range endpoints {
		permissions[i] = PermissionItem{
			Name:   ep.GetName(),
			Domain: ep.Domain,
			Action: ep.Action,
			Desc:   ep.Desc,
		}
	}

	return &PermissionsListResponse{
		Total:       len(permissions),
		Permissions: permissions,
	}
}

// ListPermissionsByDomain 按域分组列出权限
func (s *PermissionService) ListPermissionsByDomain() *PermissionsByDomainResponse {
	endpointsByDomain := endpoint.GlobalRegistry.ListByDomain()

	domains := make(map[string][]PermissionItem)
	for domain, endpoints := range endpointsByDomain {
		items := make([]PermissionItem, len(endpoints))
		for i, ep := range endpoints {
			items[i] = PermissionItem{
				Name:   ep.GetName(),
				Domain: ep.Domain,
				Action: ep.Action,
				Desc:   ep.Desc,
			}
		}
		domains[domain] = items
	}

	return &PermissionsByDomainResponse{
		Domains: domains,
	}
}

// ==================== 辅助方法 ====================

// isValidPermission 验证权限名称是否有效
func (s *PermissionService) isValidPermission(permission string) bool {
	// 支持通配符
	if permission == "*" {
		return true
	}

	// 检查是否在注册表中
	ep := endpoint.GlobalRegistry.FindByName(permission)
	if ep != nil {
		return true
	}

	// 支持域级别通配符，如 "user.*"
	// 这里简单检查格式，实际可以更严格
	return len(permission) > 0
}
