package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Material 原料
type Material struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MaterialCode string `gorm:"size:50;uniqueIndex;not null" json:"materialCode"`
	MaterialName string `gorm:"size:100;not null" json:"materialName"`
	MaterialCategory string `gorm:"size:50" json:"materialCategory"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Material) TableName() string {
	return "materials"
}

// ES Indexable 接口实现

// GetIndexName 返回 ES 索引名称
func (m *Material) GetIndexName() string {
	return "material"
}

// GetDocumentID 返回文档 ID
func (m *Material) GetDocumentID() string {
	return fmt.Sprintf("%d", m.ID)
}

// ToDocument 转换为 ES 文档
func (m *Material) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":              m.ID,
		"materialCode":    m.MaterialCode,
		"materialName":    m.MaterialName,
		"materialCategory": m.MaterialCategory,
		"createdBy":       m.CreatedBy,
		"createdAt":       m.CreatedAt,
		"updatedAt":       m.UpdatedAt,
	}
}
