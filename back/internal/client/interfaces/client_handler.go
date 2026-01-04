package interfaces

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/client/application"
	"back/internal/client/domain"
	"back/pkg/endpoint"
	"back/pkg/handler"
)

// ClientHandler 客户 Handler
type ClientHandler struct {
	service *application.ClientService
}

// NewClientHandler 创建 Handler
func NewClientHandler(service *application.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

// Create 创建客户
func (h *ClientHandler) Create(c *gin.Context) {
	handler.HandleCreate(c, h.service.Create, func(resp *application.ClientResponse) interface{} {
		return resp.ID
	})
}

// Get 获取客户
func (h *ClientHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleGet(c, uint(id), h.service.Get, domain.ErrClientNotFound)
}

// List 列表
func (h *ClientHandler) List(c *gin.Context) {
	handler.HandleList(c, h.service.List)
}

// Update 更新客户
func (h *ClientHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleUpdate(c, uint(id), h.service.Get, h.service.Update, domain.ErrClientNotFound)
}

// Delete 删除客户
func (h *ClientHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleDelete(c, uint(id), h.service.Get, h.service.Delete, domain.ErrClientNotFound)
}

// GetRoutes 获取路由定义
func (h *ClientHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/clients", Handler: h.Create, Domain: "client", Action: "create"},
		{Method: "GET", Path: "/clients/:id", Handler: h.Get, Domain: "client", Action: "get"},
		{Method: "GET", Path: "/clients", Handler: h.List, Domain: "client", Action: "list"},
		{Method: "PUT", Path: "/clients/:id", Handler: h.Update, Domain: "client", Action: "update"},
		{Method: "DELETE", Path: "/clients/:id", Handler: h.Delete, Domain: "client", Action: "delete"},
	}
}
