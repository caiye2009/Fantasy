package application

// AddUserPermissionRequest 添加用户权限请求
type AddUserPermissionRequest struct {
	LoginID    string `json:"loginId" binding:"required"`
	Permission string `json:"permission" binding:"required"`
}

// RemoveUserPermissionRequest 删除用户权限请求
type RemoveUserPermissionRequest struct {
	LoginID    string `json:"loginId" binding:"required"`
	Permission string `json:"permission" binding:"required"`
}

// GetUserPermissionsRequest 获取用户权限请求
type GetUserPermissionsRequest struct {
	LoginID string `json:"loginId" binding:"required"`
}

// UserPermissionsResponse 用户权限响应
type UserPermissionsResponse struct {
	LoginID     string   `json:"loginId"`
	Permissions []string `json:"permissions"`
}

// PermissionItem 权限项
type PermissionItem struct {
	Name   string `json:"name"`   // user.create
	Domain string `json:"domain"` // user
	Action string `json:"action"` // create
	Desc   string `json:"desc"`   // 创建用户
}

// PermissionsListResponse 权限列表响应
type PermissionsListResponse struct {
	Total       int              `json:"total"`
	Permissions []PermissionItem `json:"permissions"`
}

// PermissionsByDomainResponse 按域分组的权限列表
type PermissionsByDomainResponse struct {
	Domains map[string][]PermissionItem `json:"domains"`
}
