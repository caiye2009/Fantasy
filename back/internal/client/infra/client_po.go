package infra

import (
	"time"

	"gorm.io/gorm"

	"back/internal/client/domain"
)

// ClientPO 客户持久化对象
type ClientPO struct {
	ID        uint           `gorm:"primaryKey"`
	Code      string         `gorm:"size:50;index"` // 新增 Code 字段
	Name      string         `gorm:"size:100;not null;index"`
	Contact   string         `gorm:"size:50"`
	Phone     string         `gorm:"size:20;index"`
	Email     string         `gorm:"size:100;index"`
	Address   string         `gorm:"size:200"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 表名
func (ClientPO) TableName() string {
	return "clients"
}

// ToDomain 转换为领域模型
func (po *ClientPO) ToDomain() *domain.Client {
	return &domain.Client{
		ID:        po.ID,
		Code:      po.Code,    // 转换 Code 字段
		Name:      po.Name,
		Contact:   po.Contact,
		Phone:     po.Phone,
		Email:     po.Email,
		Address:   po.Address,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

// FromDomain 从领域模型转换
func FromDomain(c *domain.Client) *ClientPO {
	return &ClientPO{
		ID:        c.ID,
		Code:      c.Code,     // 转换 Code 字段
		Name:      c.Name,
		Contact:   c.Contact,
		Phone:     c.Phone,
		Email:     c.Email,
		Address:   c.Address,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}