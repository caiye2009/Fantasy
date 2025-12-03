package material

import (
	"strconv"
	"back/pkg/fields"
)

type Material struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Spec        string `gorm:"size:200" json:"spec"`
	Unit        string `gorm:"size:20" json:"unit"`
	Description string `gorm:"type:text" json:"description"`
	fields.DTFields
}

func (Material) TableName() string {
	return "materials"
}

// ========== Indexable 接口实现 ==========

func (m *Material) GetIndexName() string {
	return "materials"
}

func (m *Material) GetDocumentID() string {
	return strconv.Itoa(int(m.ID))
}

func (m *Material) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":          m.ID,
		"name":        m.Name,
		"spec":        m.Spec,
		"unit":        m.Unit,
		"description": m.Description,
		"created_at":  m.CreatedAt,
		"updated_at":  m.UpdatedAt,
	}
}

// ========== 其他结构体保持不变 ==========

type CreateMaterialRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Spec        string `json:"spec" binding:"omitempty,max=200"`
	Unit        string `json:"unit" binding:"omitempty,max=20"`
	Description string `json:"description" binding:"omitempty"`
}

type UpdateMaterialRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Spec        string `json:"spec" binding:"omitempty,max=200"`
	Unit        string `json:"unit" binding:"omitempty,max=20"`
	Description string `json:"description" binding:"omitempty"`
}