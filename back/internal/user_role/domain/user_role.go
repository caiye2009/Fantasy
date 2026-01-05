package domain

import (
	"time"
	"gorm.io/gorm"
)

// UserRole 用户角色
type UserRole struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	RoleCode    string         `gorm:"size:50;uniqueIndex;not null" json:"roleCode"`
	RoleName    string         `gorm:"size:100;not null" json:"roleName"`
	Description string         `gorm:"size:500" json:"description"`
	CreatedBy   string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (UserRole) TableName() string {
	return "user_roles"
}
