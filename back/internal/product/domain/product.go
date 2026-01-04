package domain

import (
	"fmt"
	"time"
	"gorm.io/gorm"
)

// Product 产品
type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductCode string `gorm:"size:50;uniqueIndex;not null" json:"productCode"`
	ProductName string `gorm:"size:100;not null" json:"productName"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Product) TableName() string {
	return "products"
}

// ES Indexable 接口实现

// GetIndexName 返回 ES 索引名称
func (p *Product) GetIndexName() string {
	return "product"
}

// GetDocumentID 返回文档 ID
func (p *Product) GetDocumentID() string {
	return fmt.Sprintf("%d", p.ID)
}

// ToDocument 转换为 ES 文档
func (p *Product) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":          p.ID,
		"productCode": p.ProductCode,
		"productName": p.ProductName,
		"createdBy":   p.CreatedBy,
		"createdAt":   p.CreatedAt,
		"updatedAt":   p.UpdatedAt,
	}
}
