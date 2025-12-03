package client

import (
	"strconv"
	"back/pkg/fields"
)

type Client struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"size:100;not null" json:"name"`
	Contact string `gorm:"size:50" json:"contact"`
	Phone   string `gorm:"size:20" json:"phone"`
	Email   string `gorm:"size:100" json:"email"`
	Address string `gorm:"size:200" json:"address"`
	fields.DTFields
}

func (Client) TableName() string {
	return "clients"
}

// ========== Indexable 接口实现 ==========

func (c *Client) GetIndexName() string {
	return "clients"
}

func (c *Client) GetDocumentID() string {
	return strconv.Itoa(int(c.ID))
}

func (c *Client) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":         c.ID,
		"name":       c.Name,
		"contact":    c.Contact,
		"phone":      c.Phone,
		"email":      c.Email,
		"address":    c.Address,
		"created_at": c.CreatedAt,
		"updated_at": c.UpdatedAt,
	}
}

// ========== 其他结构体保持不变 ==========

type CreateClientRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Contact string `json:"contact" binding:"omitempty,max=50"`
	Phone   string `json:"phone" binding:"omitempty,max=20"`
	Email   string `json:"email" binding:"omitempty,email"`
	Address string `json:"address" binding:"omitempty,max=200"`
}

type UpdateClientRequest struct {
	Name    string `json:"name" binding:"omitempty,min=2,max=100"`
	Contact string `json:"contact" binding:"omitempty,max=50"`
	Phone   string `json:"phone" binding:"omitempty,max=20"`
	Email   string `json:"email" binding:"omitempty,email"`
	Address string `json:"address" binding:"omitempty,max=200"`
}