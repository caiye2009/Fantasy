package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Fabric struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	FabricCode string         `gorm:"size:50;uniqueIndex;not null" json:"fabricCode"`
	FabricName string         `gorm:"size:100;not null" json:"fabricName"`
	CreatedBy  string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (Fabric) TableName() string {
	return "fabrics"
}

func (f *Fabric) GetIndexName() string {
	return "fabric"
}

func (f *Fabric) GetDocumentID() string {
	return fmt.Sprintf("%d", f.ID)
}

func (f *Fabric) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":         f.ID,
		"fabricCode": f.FabricCode,
		"fabricName": f.FabricName,
		"createdBy":  f.CreatedBy,
		"createdAt":  f.CreatedAt,
		"updatedAt":  f.UpdatedAt,
	}
}