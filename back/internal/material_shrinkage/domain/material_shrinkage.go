package domain

import (
	"time"

	"gorm.io/gorm"
)

// MaterialShrinkage 原料缩率表
type MaterialShrinkage struct {
	ShrinkageID string `gorm:"size:50;primaryKey" json:"shrinkageId"`
	MaterialCode string `gorm:"size:50;not null;index" json:"materialCode"`
	SupplierCode string `gorm:"size:50;not null;index" json:"supplierCode"`
	ShrinkageRate float64 `gorm:"type:decimal(5,2);not null" json:"shrinkageRate"`
	CreatedBy string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (MaterialShrinkage) TableName() string {
	return "material_shrinkages"
}