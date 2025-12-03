package vendor

import (
	"strconv"
	"back/pkg/fields"
)

type Vendor struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"size:100;not null" json:"name"`
	Contact string `gorm:"size:50" json:"contact"`
	Phone   string `gorm:"size:20" json:"phone"`
	Email   string `gorm:"size:100" json:"email"`
	Address string `gorm:"size:200" json:"address"`
	fields.DTFields
}

func (Vendor) TableName() string {
	return "vendors"
}

// ========== Indexable 接口实现 ==========

func (v *Vendor) GetIndexName() string {
	return "vendors"
}

func (v *Vendor) GetDocumentID() string {
	return strconv.Itoa(int(v.ID))
}

func (v *Vendor) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":         v.ID,
		"name":       v.Name,
		"contact":    v.Contact,
		"phone":      v.Phone,
		"email":      v.Email,
		"address":    v.Address,
		"created_at": v.CreatedAt,
		"updated_at": v.UpdatedAt,
	}
}

// ========== 其他结构体保持不变 ==========

type CreateVendorRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Contact string `json:"contact" binding:"omitempty,max=50"`
	Phone   string `json:"phone" binding:"omitempty,max=20"`
	Email   string `json:"email" binding:"omitempty,email"`
	Address string `json:"address" binding:"omitempty,max=200"`
}

type UpdateVendorRequest struct {
	Name    string `json:"name" binding:"omitempty,min=2,max=100"`
	Contact string `json:"contact" binding:"omitempty,max=50"`
	Phone   string `json:"phone" binding:"omitempty,max=20"`
	Email   string `json:"email" binding:"omitempty,email"`
	Address string `json:"address" binding:"omitempty,max=200"`
}