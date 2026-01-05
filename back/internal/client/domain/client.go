package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Client 客户
type Client struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ClientCode    string         `gorm:"size:50;uniqueIndex;not null" json:"clientCode"`
	ClientName    string         `gorm:"size:100;not null" json:"clientName"`
	ClientCountry string         `gorm:"size:50" json:"clientCountry"`
	CreatedBy     string         `gorm:"size:50;not null" json:"createdBy"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// TableName 表名
func (Client) TableName() string {
	return "clients"
}

// ES Indexable 接口实现

// GetIndexName 返回 ES 索引名称
func (c *Client) GetIndexName() string {
	return "client"
}

// GetDocumentID 返回文档 ID
func (c *Client) GetDocumentID() string {
	return fmt.Sprintf("%d", c.ID)
}

// ToDocument 转换为 ES 文档
func (c *Client) ToDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":            c.ID,
		"clientCode":    c.ClientCode,
		"clientName":    c.ClientName,
		"clientCountry": c.ClientCountry,
		"createdBy":     c.CreatedBy,
		"createdAt":     c.CreatedAt,
		"updatedAt":     c.UpdatedAt,
	}
}
