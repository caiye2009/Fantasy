package infra

import (
	"back/internal/user/domain"
	"back/pkg/repo"
	"errors"

	"gorm.io/gorm"
)

// RoleRepo 职位仓储实现
type RoleRepo struct {
	*repo.Repo[domain.Role]
	db *gorm.DB
}

// NewRoleRepo 创建职位仓储实现
func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		Repo: repo.NewRepo[domain.Role](db),
		db:   db,
	}
}

// Create 创建职位
func (r *RoleRepo) Create(role *domain.Role) error {
	return r.Repo.Create(nil, role)
}

// Update 更新职位
func (r *RoleRepo) Update(role *domain.Role) error {
	return r.Repo.Update(nil, role)
}

// FindByID 根据ID查询职位
func (r *RoleRepo) FindByID(id uint) (*domain.Role, error) {
	role, err := r.GetByID(nil, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRoleNotFound
	}
	if err != nil {
		return nil, err
	}
	return role, nil
}

// FindByCode 根据编码查询职位
func (r *RoleRepo) FindByCode(code string) (*domain.Role, error) {
	role, err := r.First(nil, map[string]interface{}{"code": code})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRoleNotFound
	}
	if err != nil {
		return nil, err
	}
	return role, nil
}

// List 查询职位列表
func (r *RoleRepo) List(status *string, page, pageSize int) ([]*domain.Role, int64, error) {
	var roles []*domain.Role
	var total int64

	query := r.db.Model(&domain.Role{})

	// 过滤状态
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("level ASC, created_at DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Delete 软删除职位
func (r *RoleRepo) Delete(id uint) error {
	return r.Repo.Delete(nil, id)
}
