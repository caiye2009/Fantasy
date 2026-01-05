package application

import "time"

// CreateClientRequest 创建客户请求
type CreateClientRequest struct {
	ClientCode    string `json:"clientCode" binding:"required,max=50"`
	ClientName    string `json:"clientName" binding:"required,max=100"`
	ClientCountry string `json:"clientCountry" binding:"max=50"`
}

// ClientResponse 客户响应
type ClientResponse struct {
	ID            uint      `json:"id"`
	ClientCode    string    `json:"clientCode"`
	ClientName    string    `json:"clientName"`
	ClientCountry string    `json:"clientCountry"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
