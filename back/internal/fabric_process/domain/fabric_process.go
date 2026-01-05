package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type FabricProcess struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	FabricCode  string         `gorm:"size:50;not null;index" json:"fabricCode"`
	ProcessCode string         `gorm:"size:50;not null" json:"processCode"`
	StepOrder   int            `gorm:"not null" json:"stepOrder"`
	CreatedBy   string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (FabricProcess) TableName() string {
	return "fabric_processes"
}

func (f *FabricProcess) GetIndexName() string {
	return "fabric_process"
}

func (f *FabricProcess) GetDocumentID() string {
	return fmt.Sprintf("%d", f.ID)
}

func (f *FabricProcess) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":          f.ID,
		"fabricCode":  f.FabricCode,
		"processCode": f.ProcessCode,
		"stepOrder":   f.StepOrder,
		"createdBy":   f.CreatedBy,
		"createdAt":   f.CreatedAt,
		"updatedAt":   f.UpdatedAt,
	}
}