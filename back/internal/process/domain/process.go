package domain

import (
	"time"

	"gorm.io/gorm"
)

// Process 工序聚合根
type Process struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
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
		"created_at":  p.CreatedAt,
		"updated_at":  p.UpdatedAt,
	}
}

// GetIndexName ES 索引名称
func (p *Process) GetIndexName() string {
	return "processes"
}

// GetDocumentID ES 文档 ID
func (p *Process) GetDocumentID() string {
	return string(rune(p.ID))
}