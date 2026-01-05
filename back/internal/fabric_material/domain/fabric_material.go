package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type FabricMaterial struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	FabricCode   string         `gorm:"size:50;not null;index" json:"fabricCode"`
	MaterialCode string         `gorm:"size:50;not null;index" json:"materialCode"`
	Ratio        float64        `gorm:"type:decimal(5,2);not null" json:"ratio"`
	CreatedBy    string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (FabricMaterial) TableName() string {
	return "fabric_materials"
}

func (f *FabricMaterial) GetIndexName() string {
	return "fabric_material"
}

func (f *FabricMaterial) GetDocumentID() string {
	return fmt.Sprintf("%d", f.ID)
}

func (f *FabricMaterial) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":           f.ID,
		"fabricCode":   f.FabricCode,
		"materialCode": f.MaterialCode,
		"ratio":        f.Ratio,
		"createdBy":    f.CreatedBy,
		"createdAt":    f.CreatedAt,
		"updatedAt":    f.UpdatedAt,
	}
}