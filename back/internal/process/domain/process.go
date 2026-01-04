package domain

import (
	"fmt"
	"time"
	"gorm.io/gorm"
)

// Process 工艺
type Process struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProcessCode string `gorm:"size:50;uniqueIndex;not null" json:"processCode"`
	ProcessName string `gorm:"size:100;not null" json:"processName"`
	ProcessCategory string `gorm:"size:100;not null" json:"processCategory"`
	CreatedBy uint           `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Process) TableName() string {
	return "processes"
}

// ES Indexable 接口实现

// GetIndexName 返回 ES 索引名称
func (p *Process) GetIndexName() string {
	return "process"
}

// GetDocumentID 返回文档 ID
func (p *Process) GetDocumentID() string {
	return fmt.Sprintf("%d", p.ID)
}

// ToDocument 转换为 ES 文档
func (p *Process) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":            p.ID,
		"processCode":   p.ProcessCode,
		"processName":   p.ProcessName,
		"processCategory": p.ProcessCategory,
		"createdBy":     p.CreatedBy,
		"createdAt":     p.CreatedAt,
		"updatedAt":     p.UpdatedAt,
	}
}