package domain

import (
	"time"
)

// 进度类型
const (
	ProgressTypeFabricInput    = "fabric_input"     // 胚布投入进度（跟单更新）
	ProgressTypeProduction     = "production"       // 加工进度（跟单更新）
	ProgressTypeWarehouseCheck = "warehouse_check"  // 验货进度（仓管更新）
	ProgressTypeRework         = "rework"           // 回修进度（跟单更新，有次品时才存在）
)

// OrderProgress 订单进度实体
type OrderProgress struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	OrderID           uint      `gorm:"not null;index" json:"order_id"`                             // 订单ID
	Type              string    `gorm:"size:30;not null;index" json:"type"`                         // 进度类型
	TargetQuantity    float64   `gorm:"type:decimal(10,2);not null" json:"target_quantity"`         // 目标数量
	CompletedQuantity float64   `gorm:"type:decimal(10,2);default:0" json:"completed_quantity"`     // 已完成数量
	Progress          int       `gorm:"default:0" json:"progress"`                                  // 百分比（0-100）
	Exists            bool      `gorm:"default:true" json:"exists"`                                 // 是否存在
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 表名
func (OrderProgress) TableName() string {
	return "order_progresses"
}

// CalculateProgress 计算进度百分比
func (p *OrderProgress) CalculateProgress() {
	if p.TargetQuantity > 0 {
		p.Progress = int((p.CompletedQuantity / p.TargetQuantity) * 100)
		if p.Progress > 100 {
			p.Progress = 100
		}
	} else {
		p.Progress = 0
	}
}

// UpdateCompleted 更新完成数量
func (p *OrderProgress) UpdateCompleted(quantity float64) error {
	if quantity < 0 {
		return ErrInvalidQuantity
	}
	p.CompletedQuantity += quantity
	p.CalculateProgress()
	return nil
}

// SetCompleted 设置完成数量（直接设置，不是累加）
func (p *OrderProgress) SetCompleted(quantity float64) error {
	if quantity < 0 {
		return ErrInvalidQuantity
	}
	p.CompletedQuantity = quantity
	p.CalculateProgress()
	return nil
}

// SetTarget 设置目标数量
func (p *OrderProgress) SetTarget(quantity float64) error {
	if quantity < 0 {
		return ErrInvalidQuantity
	}
	p.TargetQuantity = quantity
	p.CalculateProgress()
	return nil
}

// MarkAsNonExistent 标记为不存在
func (p *OrderProgress) MarkAsNonExistent() {
	p.Exists = false
	p.TargetQuantity = 0
	p.CompletedQuantity = 0
	p.Progress = 0
}

// MarkAsExistent 标记为存在
func (p *OrderProgress) MarkAsExistent() {
	p.Exists = true
}
