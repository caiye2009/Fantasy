package application

import (
	"back/internal/user/domain"
	"back/internal/user/infra"
)

// DepartmentService 部门应用服务
type DepartmentService struct {
	repo *infra.DepartmentRepo
}

// NewDepartmentService 创建部门服务
func NewDepartmentService(repo *infra.DepartmentRepo) *DepartmentService {
	return &DepartmentService{repo: repo}
}

// Create 创建部门
func (s *DepartmentService) Create(req *CreateDepartmentRequest) (*DepartmentResponse, error) {
	// 1. DTO → Domain Model
	department := ToDepartmentDomain(req)

	// 2. 领域验证
	if err := department.Validate(); err != nil {
		return nil, err
	}

	// 3. 检查编码是否重复（如果提供了编码）
	if department.Code != "" {
		existing, err := s.repo.FindByCode(department.Code)
		if err == nil && existing != nil {
			return nil, domain.ErrDepartmentCodeDuplicate
		}
		if err != nil && err != domain.ErrDepartmentNotFound {
			return nil, err
		}
	}

	// 4. 保存到数据库
	if err := s.repo.Create(department); err != nil {
		return nil, err
	}

	// 5. 返回响应
	return ToDepartmentResponse(department), nil
}

// Update 更新部门
func (s *DepartmentService) Update(id uint, req *UpdateDepartmentRequest) (*DepartmentResponse, error) {
	// 1. 查询现有部门
	department, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 2. 更新字段
	if req.Name != "" {
		department.Name = req.Name
	}
	if req.Code != "" {
		// 检查编码是否重复
		existing, err := s.repo.FindByCode(req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, domain.ErrDepartmentCodeDuplicate
		}
		if err != nil && err != domain.ErrDepartmentNotFound {
			return nil, err
		}
		department.Code = req.Code
	}
	if req.Description != "" {
		department.Description = req.Description
	}
	if req.ParentID != nil {
		department.ParentID = req.ParentID
	}

	// 3. 领域验证
	if err := department.Validate(); err != nil {
		return nil, err
	}

	// 4. 更新到数据库
	if err := s.repo.Update(department); err != nil {
		return nil, err
	}

	// 5. 返回响应
	return ToDepartmentResponse(department), nil
}

// Get 获取部门
func (s *DepartmentService) Get(id uint) (*DepartmentResponse, error) {
	department, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return ToDepartmentResponse(department), nil
}

// List 部门列表
func (s *DepartmentService) List(status *string, page, pageSize int) (*DepartmentListResponse, error) {
	departments, total, err := s.repo.List(status, page, pageSize)
	if err != nil {
		return nil, err
	}

	return ToDepartmentListResponse(departments, total), nil
}

// Deactivate 停用部门（软删除）
func (s *DepartmentService) Deactivate(id uint) error {
	// 1. 查询部门
	department, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 2. 停用部门
	department.Deactivate()

	// 3. 更新到数据库
	if err := s.repo.Update(department); err != nil {
		return err
	}

	// 4. 软删除
	return s.repo.Delete(id)
}

// Activate 激活部门
func (s *DepartmentService) Activate(id uint) error {
	// 1. 查询部门
	department, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 2. 激活部门
	department.Activate()

	// 3. 更新到数据库
	return s.repo.Update(department)
}
