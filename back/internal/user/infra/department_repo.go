package infra

import (
	"back/internal/user/domain"
	"back/pkg/repo"
	"errors"

	"gorm.io/gorm"
)

// DepartmentRepo 部门仓储实现
type DepartmentRepo struct {
	*repo.Repo[domain.Department]
	db *gorm.DB
}

// NewDepartmentRepo 创建部门仓储实现
func NewDepartmentRepo(db *gorm.DB) *DepartmentRepo {
	return &DepartmentRepo{
		Repo: repo.NewRepo[domain.Department](db),
		db:   db,
	}
}

// Create 创建部门
func (r *DepartmentRepo) Create(department *domain.Department) error {
	return r.Repo.Create(nil, department)
}

// Update 更新部门
func (r *DepartmentRepo) Update(department *domain.Department) error {
	return r.Repo.Update(nil, department)
}

// FindByID 根据ID查询部门
func (r *DepartmentRepo) FindByID(id uint) (*domain.Department, error) {
	department, err := r.GetByID(nil, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrDepartmentNotFound
	}
	if err != nil {
		return nil, err
	}
	return department, nil
}

// FindByCode 根据编码查询部门
func (r *DepartmentRepo) FindByCode(code string) (*domain.Department, error) {
	department, err := r.First(nil, map[string]interface{}{"code": code})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrDepartmentNotFound
	}
	if err != nil {
		return nil, err
	}
	return department, nil
}

// List 查询部门列表
func (r *DepartmentRepo) List(status *string, page, pageSize int) ([]*domain.Department, int64, error) {
	var departments []*domain.Department
	var total int64

	query := r.db.Model(&domain.Department{})

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
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&departments).Error; err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}

// Delete 软删除部门
func (r *DepartmentRepo) Delete(id uint) error {
	return r.Repo.Delete(nil, id)
}
