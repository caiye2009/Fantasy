package domain

import (
	"time"
	"gorm.io/gorm"
)

// User 用户
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	LoginID string `gorm:"size:50;uniqueIndex;not null" json:"loginId"`
	Username string `gorm:"size:100;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	RealName string `gorm:"size:100;not null" json:"realName"`
	Email string `gorm:"size:100" json:"email"`
	Role string `gorm:"size:50;not null;default:'user'" json:"role"`
	Department string `gorm:"size:100" json:"department"`
	IsActive bool `gorm:"default:true;not null" json:"isActive"`
	HasInitPass bool `gorm:"default:false;not null" json:"hasInitPass"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}
