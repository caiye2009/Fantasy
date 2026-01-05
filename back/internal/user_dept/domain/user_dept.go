package domain

import (
	"time"
	"gorm.io/gorm"
)

// UserDept 用户部门
type UserDept struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DeptCode    string         `gorm:"size:50;uniqueIndex;not null" json:"deptCode"`
	DeptName    string         `gorm:"size:100;not null" json:"deptName"`
	Description string         `gorm:"size:500" json:"description"`
	CreatedBy   string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (UserDept) TableName() string {
	return "user_depts"
}
