package infra

import (
	"time"
	
	"gorm.io/gorm"
	
	"back/internal/user/domain"
)

// UserPO 用户持久化对象
type UserPO struct {
	ID           uint           `gorm:"primaryKey"`
	LoginID      string         `gorm:"size:50;not null;uniqueIndex"`
	Username     string         `gorm:"size:100;not null"`
	PasswordHash string         `gorm:"size:255;not null"`
	Email        string         `gorm:"size:100"`
	Role         string         `gorm:"size:50;not null;index"`
	Status       string         `gorm:"size:20;default:'active';index"`
	HasInitPass  bool           `gorm:"default:true"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName 表名
func (UserPO) TableName() string {
	return "users"
}

// ToDomain 转换为领域模型
func (po *UserPO) ToDomain() *domain.User {
	return &domain.User{
		ID:           po.ID,
		LoginID:      po.LoginID,
		Username:     po.Username,
		PasswordHash: po.PasswordHash,
		Email:        po.Email,
		Role:         domain.UserRole(po.Role),
		Status:       domain.UserStatus(po.Status),
		HasInitPass:  po.HasInitPass,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
	}
}

// FromDomain 从领域模型转换
func FromDomain(u *domain.User) *UserPO {
	return &UserPO{
		ID:           u.ID,
		LoginID:      u.LoginID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Email:        u.Email,
		Role:         string(u.Role),
		Status:       string(u.Status),
		HasInitPass:  u.HasInitPass,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}