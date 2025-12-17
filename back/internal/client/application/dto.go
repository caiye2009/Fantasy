package application

import "time"

// CreateClientRequest 创建客户请求
type CreateClientRequest struct {
	CustomNo     string     `json:"customNo" binding:"required,max=50"`
	CustomerCode string     `json:"customerCode" binding:"omitempty,max=50"`
	InputDate    *time.Time `json:"inputDate"`
	Sales        string     `json:"sales" binding:"omitempty,max=50"`
	CustomName   string     `json:"customName" binding:"required,min=2,max=200"`
	StateChNm    string     `json:"stateChNm" binding:"omitempty,max=100"`
	Country      string     `json:"country" binding:"omitempty,max=50"`
	Address      string     `json:"address" binding:"omitempty,max=500"`
	AddressEn    string     `json:"addressEn" binding:"omitempty,max=500"`
	CustomNameEn string     `json:"customNameEn" binding:"omitempty,max=200"`
	Contactor    string     `json:"contactor" binding:"omitempty,max=100"`
	UnitPhone    string     `json:"unitPhone" binding:"omitempty,max=50"`
	Mobile       string     `json:"mobile" binding:"omitempty,max=50"`
	FaxNum       string     `json:"faxNum" binding:"omitempty,max=50"`
	Email        string     `json:"email" binding:"omitempty,email"`
	PyCustomName string     `json:"pyCustomName" binding:"omitempty,max=200"`
	CheckRequest string     `json:"checkRequest" binding:"omitempty"`
	CustomStatus string     `json:"customStatus" binding:"omitempty,max=50"`
	DocMan       string     `json:"docMan" binding:"omitempty,max=50"`
}

// UpdateClientRequest 更新客户请求
type UpdateClientRequest struct {
	CustomNo     string     `json:"customNo" binding:"omitempty,max=50"`
	CustomerCode string     `json:"customerCode" binding:"omitempty,max=50"`
	InputDate    *time.Time `json:"inputDate"`
	Sales        string     `json:"sales" binding:"omitempty,max=50"`
	CustomName   string     `json:"customName" binding:"omitempty,min=2,max=200"`
	StateChNm    string     `json:"stateChNm" binding:"omitempty,max=100"`
	Country      string     `json:"country" binding:"omitempty,max=50"`
	Address      string     `json:"address" binding:"omitempty,max=500"`
	AddressEn    string     `json:"addressEn" binding:"omitempty,max=500"`
	CustomNameEn string     `json:"customNameEn" binding:"omitempty,max=200"`
	Contactor    string     `json:"contactor" binding:"omitempty,max=100"`
	UnitPhone    string     `json:"unitPhone" binding:"omitempty,max=50"`
	Mobile       string     `json:"mobile" binding:"omitempty,max=50"`
	FaxNum       string     `json:"faxNum" binding:"omitempty,max=50"`
	Email        string     `json:"email" binding:"omitempty,email"`
	PyCustomName string     `json:"pyCustomName" binding:"omitempty,max=200"`
	CheckRequest string     `json:"checkRequest" binding:"omitempty"`
	CustomStatus string     `json:"customStatus" binding:"omitempty,max=50"`
	DocMan       string     `json:"docMan" binding:"omitempty,max=50"`
}

// ClientResponse 客户响应
type ClientResponse struct {
	ID           uint       `json:"id"`
	CustomNo     string     `json:"customNo"`
	CustomerCode string     `json:"customerCode"`
	InputDate    *time.Time `json:"inputDate"`
	Sales        string     `json:"sales"`
	CustomName   string     `json:"customName"`
	StateChNm    string     `json:"stateChNm"`
	Country      string     `json:"country"`
	Address      string     `json:"address"`
	AddressEn    string     `json:"addressEn"`
	CustomNameEn string     `json:"customNameEn"`
	Contactor    string     `json:"contactor"`
	UnitPhone    string     `json:"unitPhone"`
	Mobile       string     `json:"mobile"`
	FaxNum       string     `json:"faxNum"`
	Email        string     `json:"email"`
	PyCustomName string     `json:"pyCustomName"`
	CheckRequest string     `json:"checkRequest"`
	CustomStatus string     `json:"customStatus"`
	DocMan       string     `json:"docMan"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// ClientListResponse 客户列表响应
type ClientListResponse struct {
	Total   int64             `json:"total"`
	Clients []*ClientResponse `json:"clients"`
}