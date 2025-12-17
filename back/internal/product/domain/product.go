package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

// 产品状态常量
const (
	ProductStatusDraft     = "draft"     // 草稿
	ProductStatusSubmitted = "submitted" // 已提交
	ProductStatusApproved  = "approved"  // 已审批
	ProductStatusRejected  = "rejected"  // 已拒绝
)

// MaterialConfig 原料配置
type MaterialConfig struct {
	MaterialID uint    `json:"material_id"`
	Ratio      float64 `json:"ratio"` // 占比（总和必须为1）
}

// ProcessConfig 工艺配置
type ProcessConfig struct {
	ProcessID uint    `json:"process_id"`
	Quantity  float64 `json:"quantity"` // 数量（可选）
}

// MaterialConfigJSON 原料配置 JSON（用于 GORM JSONB）
type MaterialConfigJSON []MaterialConfig

// Scan 实现 sql.Scanner 接口
func (m *MaterialConfigJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

// Value 实现 driver.Valuer 接口
func (m MaterialConfigJSON) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}

// ProcessConfigJSON 工艺配置 JSON（用于 GORM JSONB）
type ProcessConfigJSON []ProcessConfig

// Scan 实现 sql.Scanner 接口
func (p *ProcessConfigJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// Value 实现 driver.Valuer 接口
func (p ProcessConfigJSON) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}
	return json.Marshal(p)
}

// Product 产品聚合根
type Product struct {
	ID        uint               `gorm:"primaryKey" json:"id"`
	Name      string             `gorm:"size:100;not null;index" json:"name"`
	Status    string             `gorm:"size:20;default:draft;index" json:"status"`
	Materials MaterialConfigJSON `gorm:"type:jsonb" json:"materials"`
	Processes ProcessConfigJSON  `gorm:"type:jsonb" json:"processes"`
	CreatedAt time.Time          `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time          `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
}

// TableName 表名
func (Product) TableName() string {
	return "products"
}

// Validate 验证产品数据
func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrProductNameEmpty
	}
	
	if len(p.Name) < 2 || len(p.Name) > 100 {
		return ErrProductNameInvalid
	}
	
	if len(p.Materials) == 0 {
		return ErrMaterialsRequired
	}
	
	// 验证原料占比总和为1
	sum := 0.0
	for _, m := range p.Materials {
		if m.Ratio <= 0 {
			return ErrInvalidMaterialRatio
		}
		sum += m.Ratio
	}
	if math.Abs(sum-1.0) > 0.0001 {
		return ErrMaterialRatioSumNotOne
	}
	
	if len(p.Processes) == 0 {
		return ErrProcessesRequired
	}
	
	return nil
}

// Submit 提交产品（从草稿变为已提交）
func (p *Product) Submit() error {
	if p.Status != ProductStatusDraft {
		return ErrCannotSubmit
	}

	// 提交前验证
	if err := p.Validate(); err != nil {
		return err
	}

	p.Status = ProductStatusSubmitted
	return nil
}

// Approve 审批通过
func (p *Product) Approve() error {
	if p.Status != ProductStatusSubmitted {
		return ErrCannotApprove
	}

	p.Status = ProductStatusApproved
	return nil
}

// Reject 拒绝
func (p *Product) Reject() error {
	if p.Status != ProductStatusSubmitted {
		return ErrCannotReject
	}

	p.Status = ProductStatusRejected
	return nil
}

// UpdateName 更新名称
func (p *Product) UpdateName(newName string) error {
	if newName == "" {
		return ErrProductNameEmpty
	}
	if len(newName) < 2 || len(newName) > 100 {
		return ErrProductNameInvalid
	}
	
	// 已审批的产品不能修改
	if p.Status == ProductStatusApproved {
		return ErrCannotUpdateApprovedProduct
	}
	
	p.Name = newName
	return nil
}

// UpdateMaterials 更新原料配置
func (p *Product) UpdateMaterials(materials []MaterialConfig) error {
	if len(materials) == 0 {
		return ErrMaterialsRequired
	}
	
	// 验证占比总和
	sum := 0.0
	for _, m := range materials {
		if m.Ratio <= 0 {
			return ErrInvalidMaterialRatio
		}
		sum += m.Ratio
	}
	if math.Abs(sum-1.0) > 0.0001 {
		return ErrMaterialRatioSumNotOne
	}
	
	// 已审批的产品不能修改
	if p.Status == ProductStatusApproved {
		return ErrCannotUpdateApprovedProduct
	}
	
	p.Materials = materials
	return nil
}

// UpdateProcesses 更新工艺配置
func (p *Product) UpdateProcesses(processes []ProcessConfig) error {
	if len(processes) == 0 {
		return ErrProcessesRequired
	}
	
	// 已审批的产品不能修改
	if p.Status == ProductStatusApproved {
		return ErrCannotUpdateApprovedProduct
	}
	
	p.Processes = processes
	return nil
}

// CanDelete 是否可以删除
func (p *Product) CanDelete() bool {
	// 已审批的产品不能删除
	return p.Status != ProductStatusApproved
}

// IsApproved 是否已审批
func (p *Product) IsApproved() bool {
	return p.Status == ProductStatusApproved
}

// ToDocument 转换为 ES 文档
func (p *Product) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":         p.ID,
		"name":       p.Name,
		"status":     p.Status,
		"materials":  p.Materials,
		"processes":  p.Processes,
		"createdAt":  p.CreatedAt,
		"updatedAt":  p.UpdatedAt,
	}
}

// GetIndexName ES 索引名称
func (p *Product) GetIndexName() string {
	return "product"
}

// GetDocumentID ES 文档 ID
func (p *Product) GetDocumentID() string {
	return fmt.Sprintf("%d", p.ID)
}