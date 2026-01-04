package interfaces

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/process/application"
	"back/internal/process/domain"
	"back/pkg/endpoint"
	"back/pkg/handler"
)

// ProcessHandler 工序 Handler
type ProcessHandler struct {
	service *application.ProcessService
}

// NewProcessHandler 创建 Handler
func NewProcessHandler(service *application.ProcessService) *ProcessHandler {
	return &ProcessHandler{service: service}
}

// Create 创建工序
func (h *ProcessHandler) Create(c *gin.Context) {
	handler.HandleCreate(c, h.service.Create, func(resp *application.ProcessResponse) interface{} {
		return resp.ID
	})
}

// Get 获取工序
func (h *ProcessHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleGet(c, uint(id), h.service.Get, domain.ErrProcessNotFound)
}

// List 列表
func (h *ProcessHandler) List(c *gin.Context) {
	handler.HandleList(c, h.service.List)
}

// Update 更新工序
func (h *ProcessHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleUpdate(c, uint(id), h.service.Get, h.service.Update, domain.ErrProcessNotFound)
}

// Delete 删除工序
func (h *ProcessHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleDelete(c, uint(id), h.service.Get, h.service.Delete, domain.ErrProcessNotFound)
}

// GetRoutes 获取路由定义
func (h *ProcessHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/processs", Handler: h.Create, Domain: "process", Action: "create"},
		{Method: "GET", Path: "/processs/:id", Handler: h.Get, Domain: "process", Action: "get"},
		{Method: "GET", Path: "/processs", Handler: h.List, Domain: "process", Action: "list"},
		{Method: "PUT", Path: "/processs/:id", Handler: h.Update, Domain: "process", Action: "update"},
		{Method: "DELETE", Path: "/processs/:id", Handler: h.Delete, Domain: "process", Action: "delete"},
	}
}
