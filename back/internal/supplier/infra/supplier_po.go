package infra

import (
	"time"

	"gorm.io/gorm"

	"back/internal/supplier/domain"
)

// SupplierPO 供应商持久化对象
type SupplierPO struct {
	ID        uint           `gorm:"primaryKey"`
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
func (SupplierPO) TableName() string {
	return "suppliers"
}

// ToDomain 转换为领域模型
func (po *SupplierPO) ToDomain() *domain.Supplier {
	return &domain.Supplier{
		ID:        po.ID,
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
func FromDomain(s *domain.Supplier) *SupplierPO {
	return &SupplierPO{
		ID:        s.ID,
		Name:      s.Name,
		Contact:   s.Contact,
		Phone:     s.Phone,
		Email:     s.Email,
		Address:   s.Address,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
