package interfaces

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/material/application"
	"back/internal/material/domain"
	"back/pkg/endpoint"
	"back/pkg/handler"
)

// MaterialHandler 客户 Handler
type MaterialHandler struct {
	service *application.MaterialService
}

// NewMaterialHandler 创建 Handler
func NewMaterialHandler(service *application.MaterialService) *MaterialHandler {
	return &MaterialHandler{service: service}
}

// Create 创建客户
func (h *MaterialHandler) Create(c *gin.Context) {
	handler.HandleCreate(c, h.service.Create, func(resp *application.MaterialResponse) interface{} {
		return resp.ID
	})
}

// Get 获取客户
func (h *MaterialHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleGet(c, uint(id), h.service.Get, domain.ErrMaterialNotFound)
}

// List 列表
func (h *MaterialHandler) List(c *gin.Context) {
	handler.HandleList(c, h.service.List)
}

// Update 更新客户
func (h *MaterialHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleUpdate(c, uint(id), h.service.Get, h.service.Update, domain.ErrMaterialNotFound)
}

// Delete 删除客户
func (h *MaterialHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleDelete(c, uint(id), h.service.Get, h.service.Delete, domain.ErrMaterialNotFound)
}

// GetRoutes 获取路由定义
func (h *MaterialHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/materials", Handler: h.Create, Domain: "material", Action: "create"},
		{Method: "GET", Path: "/materials/:id", Handler: h.Get, Domain: "material", Action: "get"},
		{Method: "GET", Path: "/materials", Handler: h.List, Domain: "material", Action: "list"},
		{Method: "PUT", Path: "/materials/:id", Handler: h.Update, Domain: "material", Action: "update"},
		{Method: "DELETE", Path: "/materials/:id", Handler: h.Delete, Domain: "material", Action: "delete"},
	}
}
