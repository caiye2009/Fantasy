package application

import (
	"time"
	"back/internal/user/domain"
)

// ToRoleDomain 将创建请求转换为领域对象
func ToRoleDomain(req *CreateRoleRequest) *domain.Role {
	return &domain.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      domain.RoleStatusActive,
		Level:       req.Level,
	}
}

// ToRoleResponse 将领域对象转换为响应
func ToRoleResponse(role *domain.Role) *RoleResponse {
	var deletedAt *time.Time
	if role.DeletedAt.Valid {
		deletedAt = &role.DeletedAt.Time
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Status:      role.Status,
		Level:       role.Level,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

// ToRoleListResponse 将领域对象列表转换为响应
func ToRoleListResponse(roles []*domain.Role, total int64) *RoleListResponse {
	responses := make([]*RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = ToRoleResponse(role)
	}
	return &RoleListResponse{
		Total: total,
		Roles: responses,
	}
}
