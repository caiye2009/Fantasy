package interfaces

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/fabric/application"
	"back/internal/fabric/domain"
	"back/pkg/endpoint"
	"back/pkg/handler"
)

// FabricHandler Handler
type FabricHandler struct {
	service *application.FabricService
}

// NewFabricHandler 创建 Handler
func NewFabricHandler(service *application.FabricService) *FabricHandler {
	return &FabricHandler{service: service}
}

// Create 创建面料
func (h *FabricHandler) Create(c *gin.Context) {
	handler.HandleCreate(c, h.service.Create, func(resp *application.FabricResponse) interface{} {
		return resp.ID
	})
}

// Get 获取面料
func (h *FabricHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleGet(c, uint(id), h.service.Get, domain.ErrFabricNotFound)
}

// List 列表
func (h *FabricHandler) List(c *gin.Context) {
	handler.HandleList(c, h.service.List)
}

// Update 更新面料
func (h *FabricHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleUpdate(c, uint(id), h.service.Get, h.service.Update, domain.ErrFabricNotFound)
}

// Delete 删除面料
func (h *FabricHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleDelete(c, uint(id), h.service.Get, h.service.Delete, domain.ErrFabricNotFound)
}

// GetRoutes 获取路由定义
func (h *FabricHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/fabrics", Handler: h.Create, Domain: "fabric", Action: "create"},
		{Method: "GET", Path: "/fabrics/:id", Handler: h.Get, Domain: "fabric", Action: "get"},
		{Method: "GET", Path: "/fabrics", Handler: h.List, Domain: "fabric", Action: "list"},
		{Method: "PUT", Path: "/fabrics/:id", Handler: h.Update, Domain: "fabric", Action: "update"},
		{Method: "DELETE", Path: "/fabrics/:id", Handler: h.Delete, Domain: "fabric", Action: "delete"},
	}
}