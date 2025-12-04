package domain

import (
	"time"
)

// Material 材料聚合根
type Material struct {
	ID          uint
	Name        string
	Spec        string
	Unit        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate 验证材料数据
func (m *Material) Validate() error {
	if m.Name == "" {
		return ErrMaterialNameEmpty
	}
	
	if len(m.Name) < 2 || len(m.Name) > 100 {
		return ErrMaterialNameInvalid
	}
	
	if len(m.Spec) > 200 {
		return ErrMaterialSpecTooLong
	}
	
	if len(m.Unit) > 20 {
		return ErrMaterialUnitTooLong
	}
	
	return nil
}

// UpdateSpec 更新规格
func (m *Material) UpdateSpec(newSpec string) error {
	if len(newSpec) > 200 {
		return ErrMaterialSpecTooLong
	}
	m.Spec = newSpec
	return nil
}

// UpdateName 更新名称
func (m *Material) UpdateName(newName string) error {
	if newName == "" {
		return ErrMaterialNameEmpty
	}
	if len(newName) < 2 || len(newName) > 100 {
		return ErrMaterialNameInvalid
	}
	m.Name = newName
	return nil
}

// ToDocument 转换为 ES 文档
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

// GetIndexName ES 索引名称
func (m *Material) GetIndexName() string {
	return "materials"
}

// GetDocumentID ES 文档 ID
func (m *Material) GetDocumentID() string {
	return string(rune(m.ID))
}