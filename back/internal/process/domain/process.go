package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Process 工序聚合根
type Process struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Process) TableName() string {
	return "processes"
}

// Validate 验证工序数据
func (p *Process) Validate() error {
	if p.Name == "" {
		return ErrProcessNameEmpty
	}
	
	if len(p.Name) < 2 || len(p.Name) > 100 {
		return ErrProcessNameInvalid
	}
	
	return nil
}

// UpdateName 更新名称
func (p *Process) UpdateName(newName string) error {
	if newName == "" {
		return ErrProcessNameEmpty
	}
	if len(newName) < 2 || len(newName) > 100 {
		return ErrProcessNameInvalid
	}
	p.Name = newName
	return nil
}

// ToDocument 转换为 ES 文档
func (p *Process) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":          p.ID,
		"name":        p.Name,
		"description": p.Description,
		"createdAt":   p.CreatedAt,
		"updatedAt":   p.UpdatedAt,
	}
}

// GetIndexName ES 索引名称
func (p *Process) GetIndexName() string {
	return "process"
}

// GetDocumentID ES 文档 ID
func (p *Process) GetDocumentID() string {
	return fmt.Sprintf("%d", p.ID)
}