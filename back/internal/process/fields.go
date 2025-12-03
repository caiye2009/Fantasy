package process

import (
	"strconv"
	"back/pkg/fields"
)

type Process struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	fields.DTFields
}

func (Process) TableName() string {
	return "processes"
}

// ========== Indexable 接口实现 ==========

func (p *Process) GetIndexName() string {
	return "processes"
}

func (p *Process) GetDocumentID() string {
	return strconv.Itoa(int(p.ID))
}

func (p *Process) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":          p.ID,
		"name":        p.Name,
		"description": p.Description,
		"created_at":  p.CreatedAt,
		"updated_at":  p.UpdatedAt,
	}
}

// ========== 其他结构体保持不变 ==========

type CreateProcessRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty"`
}

type UpdateProcessRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"omitempty"`
}