package domain

import (
	"time"

	"gorm.io/gorm"
)

// 库存单位常量
const (
	UnitMeter     = "米"  // 米
	UnitKilogram  = "kg"  // 公斤
)

// 库存类别常量
const (
	CategoryRawMaterial  = "raw_material"  // 原材料
	CategorySemiFinished = "semi_finished" // 半成品
	CategoryFinished     = "finished"      // 成品
)

// Inventory 库存聚合根
type Inventory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProductID   uint           `gorm:"not null;index" json:"productId"`        // 产品ID
	Category    string         `gorm:"size:50;not null;index" json:"category"` // 类别
	BatchID     string         `gorm:"size:100;not null;index" json:"batchId"` // 批次ID
	Quantity    float64        `gorm:"not null" json:"quantity"`               // 数量
	Unit        string         `gorm:"size:20;not null" json:"unit"`           // 单位（米或kg）
	UnitCost    float64        `gorm:"not null" json:"unitCost"`               // 单价
	TotalCost   float64        `gorm:"not null" json:"totalCost"`              // 总成本
	Remark      string         `gorm:"size:500" json:"remark"`                 // 备注
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Inventory) TableName() string {
	return "inventories"
}

// Validate 验证库存数据
func (i *Inventory) Validate() error {
	if i.ProductID == 0 {
		return ErrProductIDRequired
	}

	if i.Category == "" {
		return ErrCategoryRequired
	}

	if i.BatchID == "" {
		return ErrBatchIDRequired
	}

	if i.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	if i.Unit == "" {
		return ErrUnitRequired
	}

	if i.UnitCost < 0 {
		return ErrInvalidUnitCost
	}

	if i.TotalCost < 0 {
		return ErrInvalidTotalCost
	}

	return nil
}

// CalculateTotalCost 计算总成本
func (i *Inventory) CalculateTotalCost() {
	i.TotalCost = i.Quantity * i.UnitCost
}

// UpdateQuantity 更新数量
func (i *Inventory) UpdateQuantity(quantity float64) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	i.Quantity = quantity
	i.CalculateTotalCost()
	return nil
}

// UpdateUnitCost 更新单价
func (i *Inventory) UpdateUnitCost(unitCost float64) error {
	if unitCost < 0 {
		return ErrInvalidUnitCost
	}

	i.UnitCost = unitCost
	i.CalculateTotalCost()
	return nil
}

// Deduct 扣减库存
func (i *Inventory) Deduct(quantity float64) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if i.Quantity < quantity {
		return ErrInsufficientInventory
	}

	i.Quantity -= quantity
	i.CalculateTotalCost()
	return nil
}

// Add 增加库存
func (i *Inventory) Add(quantity float64) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	i.Quantity += quantity
	i.CalculateTotalCost()
	return nil
}
