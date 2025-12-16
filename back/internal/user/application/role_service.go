package application

import (
	"time"

	"back/internal/user/domain"
	"back/internal/user/infra"
)

// RoleService 职位应用服务
type RoleService struct {
	repo *infra.RoleRepo
}

// NewRoleService 创建职位服务
func NewRoleService(repo *infra.RoleRepo) *RoleService {
	return &RoleService{repo: repo}
}

// Create 创建职位
func (s *RoleService) Create(req *CreateRoleRequest) (*RoleResponse, error) {
	// 1. DTO → Domain Model
	role := &domain.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      domain.RoleStatusActive,
		Level:       req.Level,
	}

	// 2. 领域验证
	if err := role.Validate(); err != nil {
		return nil, err
	}

	// 3. 检查编码是否重复
	existing, err := s.repo.FindByCode(role.Code)
	if err == nil && existing != nil {
		return nil, domain.ErrRoleCodeDuplicate
	}
	if err != nil && err != domain.ErrRoleNotFound {
		return nil, err
	}

	// 4. 保存到数据库
	if err := s.repo.Create(role); err != nil {
		return nil, err
	}

	// 5. 返回响应
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
	}, nil
}

// Update 更新职位
func (s *RoleService) Update(id uint, req *UpdateRoleRequest) (*RoleResponse, error) {
	// 1. 查询现有职位
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 2. 更新字段
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Code != "" {
		// 检查编码是否重复
		existing, err := s.repo.FindByCode(req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, domain.ErrRoleCodeDuplicate
		}
		if err != nil && err != domain.ErrRoleNotFound {
			return nil, err
		}
		role.Code = req.Code
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Level > 0 {
		role.Level = req.Level
	}

	// 3. 领域验证
	if err := role.Validate(); err != nil {
		return nil, err
	}

	// 4. 更新到数据库
	if err := s.repo.Update(role); err != nil {
		return nil, err
	}

	// 5. 返回响应
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
	}, nil
}

// Get 获取职位
func (s *RoleService) Get(id uint) (*RoleResponse, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

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
	}, nil
}

// List 职位列表
func (s *RoleService) List(status *string, page, pageSize int) (*RoleListResponse, error) {
	roles, total, err := s.repo.List(status, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]*RoleResponse, len(roles))
	for i, role := range roles {
		var deletedAt *time.Time
		if role.DeletedAt.Valid {
			deletedAt = &role.DeletedAt.Time
		}
		responses[i] = &RoleResponse{
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

	return &RoleListResponse{
		Total: total,
		Roles: responses,
	}, nil
}

// Deactivate 停用职位（软删除）
func (s *RoleService) Deactivate(id uint) error {
	// 1. 查询职位
	role, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 2. 停用职位
	role.Deactivate()

	// 3. 更新到数据库
	if err := s.repo.Update(role); err != nil {
		return err
	}

	// 4. 软删除
	return s.repo.Delete(id)
}

// Activate 激活职位
func (s *RoleService) Activate(id uint) error {
	// 1. 查询职位
	role, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 2. 激活职位
	role.Activate()

	// 3. 更新到数据库
	return s.repo.Update(role)
}
