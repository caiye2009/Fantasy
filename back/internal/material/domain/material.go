package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Material 材料聚合根
type Material struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;index" json:"name"`
	Spec        string         `gorm:"size:200" json:"spec"`
	Unit        string         `gorm:"size:20" json:"unit"`
	Category    string         `gorm:"size:50;index" json:"category"`
	SupplierID  uint           `gorm:"index" json:"supplierId"`
	Price       float64        `gorm:"type:decimal(10,2)" json:"price"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Material) TableName() string {
	return "materials"
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

// ToDocument 转换为 ES 文档（小驼峰字段名）
func (m *Material) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":          m.ID,
		"name":        m.Name,
		"spec":        m.Spec,
		"unit":        m.Unit,
		"category":    m.Category,
		"supplierId":  m.SupplierID,
		"price":       m.Price,
		"status":      m.Status,
		"description": m.Description,
		"createdAt":   m.CreatedAt,
		"updatedAt":   m.UpdatedAt,
	}
}

// GetIndexName ES 索引名称
func (m *Material) GetIndexName() string {
	return "material"
}

// GetDocumentID ES 文档 ID
func (m *Material) GetDocumentID() string {
	return fmt.Sprintf("%d", m.ID)
}