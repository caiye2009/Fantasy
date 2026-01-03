package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/process/application"
	"back/internal/process/domain"
	"back/pkg/endpoint"
)

// ProcessHandler 客户 Handler
type ProcessHandler struct {
	service *application.ProcessService
}

// NewProcessHandler 创建 Handler
func NewProcessHandler(service *application.ProcessService) *ProcessHandler {
	return &ProcessHandler{service: service}
}

// Create 创建客户
func (h *ProcessHandler) Create(c *gin.Context) {
	var req application.CreateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Get 获取客户
func (h *ProcessHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List 列表
func (h *ProcessHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	list, total, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"list":  list,
	})
}

// Update 更新客户
func (h *ProcessHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete 删除客户
func (h *ProcessHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
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
