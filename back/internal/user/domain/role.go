package domain

import (
	"time"
	"gorm.io/gorm"
)

// 职位状态常量
const (
	RoleStatusActive   = "active"   // 激活
	RoleStatusInactive = "inactive" // 停用
)

// Role 职位聚合根
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"size:50;not null;uniqueIndex" json:"code"`
	Description string         `gorm:"size:500" json:"description"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	Level       int            `gorm:"default:0" json:"level"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Role) TableName() string {
	return "roles"
}

// Validate 验证职位数据
func (r *Role) Validate() error {
	if r.Name == "" {
		return ErrRoleNameEmpty
	}

	if len(r.Name) < 2 || len(r.Name) > 100 {
		return ErrRoleNameInvalid
	}

	if r.Code == "" {
		return ErrRoleCodeEmpty
	}

	if len(r.Code) < 2 || len(r.Code) > 50 {
		return ErrRoleCodeInvalid
	}

	return nil
}

// Deactivate 停用职位
func (r *Role) Deactivate() {
	r.Status = RoleStatusInactive
}

// Activate 激活职位
func (r *Role) Activate() {
	r.Status = RoleStatusActive
}

// RoleRepository 职位仓储接口
type RoleRepository interface {
	Create(role *Role) error
	Update(role *Role) error
	FindByID(id uint) (*Role, error)
	FindByCode(code string) (*Role, error)
	List(status *string, page, pageSize int) ([]*Role, int64, error)
	Delete(id uint) error
}
