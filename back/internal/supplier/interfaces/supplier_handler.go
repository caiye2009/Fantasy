package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/pkg/audit"
	"back/pkg/endpoint"
	"back/internal/supplier/application"
	"back/internal/supplier/domain"
)

// SupplierHandler 供应商 Handler
type SupplierHandler struct {
	service *application.SupplierService
}

// NewSupplierHandler 创建 Handler
func NewSupplierHandler(service *application.SupplierService) *SupplierHandler {
	return &SupplierHandler{service: service}
}

// Create 创建供应商
// @Summary      创建供应商
// @Description  创建新的供应商信息
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateSupplierRequest true "供应商信息"
// @Success      200 {object} application.SupplierResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /supplier [post]
func (h *SupplierHandler) Create(c *gin.Context) {
	var req application.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrSupplierNameExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "供应商名称已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Get 获取供应商详情
// @Summary      获取供应商详情
// @Description  根据供应商ID获取供应商详细信息
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Param        id path int true "供应商ID"
// @Success      200 {object} application.SupplierResponse "获取成功"
// @Failure      404 {object} map[string]string "供应商不存在"
// @Security     Bearer
// @Router       /supplier/{id} [get]
func (h *SupplierHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrSupplierNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List 获取供应商列表
// @Summary      获取供应商列表
// @Description  获取所有供应商列表
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} application.SupplierListResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /supplier [get]
func (h *SupplierHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Update 更新供应商信息
// @Summary      更新供应商信息
// @Description  根据供应商ID更新供应商信息
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Param        id path int true "供应商ID"
// @Param        request body application.UpdateSupplierRequest true "更新的供应商信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /supplier/{id} [post]
func (h *SupplierHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		if errors.Is(err, domain.ErrSupplierNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
			return
		}
		if errors.Is(err, domain.ErrSupplierNameExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "供应商名称已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除供应商
// @Summary      删除供应商
// @Description  根据供应商ID删除供应商
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Param        id path int true "供应商ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /supplier/{id} [delete]
func (h *SupplierHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrSupplierNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetRoutes 获取路由定义
func (h *SupplierHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{
			Method:      "POST",
			Path:        "/supplier",
			Handler:     h.Create,
			Middlewares: []gin.HandlerFunc{audit.Mark("supplier", "create")},
			Name:        "创建供应商",
		},
		{
			Method:      "GET",
			Path:        "/supplier/:id",
			Handler:     h.Get,
			Middlewares: nil,
			Name:        "获取供应商详情",
		},
		{
			Method:      "GET",
			Path:        "/supplier",
			Handler:     h.List,
			Middlewares: nil,
			Name:        "获取供应商列表",
		},
		{
			Method:      "PUT",
			Path:        "/supplier/:id",
			Handler:     h.Update,
			Middlewares: []gin.HandlerFunc{audit.Mark("supplier", "update")},
			Name:        "更新供应商信息",
		},
		{
			Method:      "DELETE",
			Path:        "/supplier/:id",
			Handler:     h.Delete,
			Middlewares: []gin.HandlerFunc{audit.Mark("supplier", "delete")},
			Name:        "删除供应商",
		},
	}
}
