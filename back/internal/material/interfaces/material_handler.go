package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/material/application"
	"back/internal/material/domain"
	"back/pkg/endpoint"
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
	var req application.CreateMaterialRequest
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
func (h *MaterialHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "material not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List 列表
func (h *MaterialHandler) List(c *gin.Context) {
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
func (h *MaterialHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "material not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete 删除客户
func (h *MaterialHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "material not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
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
