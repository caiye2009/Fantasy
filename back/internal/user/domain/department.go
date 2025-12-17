package domain

import (
	"time"
	"gorm.io/gorm"
)

// 部门状态常量
const (
	DepartmentStatusActive   = "active"   // 激活
	DepartmentStatusInactive = "inactive" // 停用
)

// Department 部门聚合根
type Department struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"size:50;uniqueIndex" json:"code"`
	Description string         `gorm:"size:500" json:"description"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	ParentID    *uint          `gorm:"index" json:"parentId"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Department) TableName() string {
	return "departments"
}

// Validate 验证部门数据
func (d *Department) Validate() error {
	if d.Name == "" {
		return ErrDepartmentNameEmpty
	}

	if len(d.Name) < 2 || len(d.Name) > 100 {
		return ErrDepartmentNameInvalid
	}

	if d.Code != "" && (len(d.Code) < 2 || len(d.Code) > 50) {
		return ErrDepartmentCodeInvalid
	}

	return nil
}

// Deactivate 停用部门
func (d *Department) Deactivate() {
	d.Status = DepartmentStatusInactive
}

// Activate 激活部门
func (d *Department) Activate() {
	d.Status = DepartmentStatusActive
}

// DepartmentRepository 部门仓储接口
type DepartmentRepository interface {
	Create(department *Department) error
	Update(department *Department) error
	FindByID(id uint) (*Department, error)
	FindByCode(code string) (*Department, error)
	List(status *string, page, pageSize int) ([]*Department, int64, error)
	Delete(id uint) error
}
