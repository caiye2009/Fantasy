package domain

import (
	"time"

	"gorm.io/gorm"
)

// MaterialAllocation 胚布调配
type MaterialAllocation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	AllocationCode string `gorm:"size:50;uniqueIndex;not null" json:"allocationCode"`
	OrderCode string `gorm:"size:50;not null;index" json:"orderCode"`
	MaterialCode string `gorm:"size:50;not null;index" json:"materialCode"`
	AllocationQty float64 `gorm:"type:decimal(10,2);not null" json:"allocationQty"`
	ActualQty float64 `gorm:"type:decimal(10,2)" json:"actualQty"`
	AllocationStatus string `gorm:"size:50;not null;default:'pending'" json:"allocationStatus"`
	AllocatedBy uint `gorm:"not null" json:"allocatedBy"`
	AllocatedAt time.Time `json:"allocatedAt"`
	CreatedBy string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (MaterialAllocation) TableName() string {
	return "material_allocations"
}
