package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Supplier 供应商
type Supplier struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	SupplierCode     string         `gorm:"size:50;uniqueIndex;not null" json:"supplierCode"`
	SupplierName     string         `gorm:"size:100;not null" json:"supplierName"`
	SupplierCategory string         `gorm:"size:50" json:"supplierCategory"`
	CreatedBy        string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Supplier) TableName() string {
	return "suppliers"
}

// ES Indexable 接口实现

// GetIndexName 返回 ES 索引名称
func (s *Supplier) GetIndexName() string {
	return "supplier"
}

// GetDocumentID 返回文档 ID
func (s *Supplier) GetDocumentID() string {
	return fmt.Sprintf("%d", s.ID)
}

// ToDocument 转换为 ES 文档
func (s *Supplier) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":               s.ID,
		"supplierCode":     s.SupplierCode,
		"supplierName":     s.SupplierName,
		"supplierCategory": s.SupplierCategory,
		"createdBy":        s.CreatedBy,
		"createdAt":        s.CreatedAt,
		"updatedAt":        s.UpdatedAt,
	}
}
