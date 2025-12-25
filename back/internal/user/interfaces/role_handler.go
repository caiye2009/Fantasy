package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/pkg/audit"
	"back/pkg/endpoint"
	"back/internal/user/application"
	"back/internal/user/domain"
)

// RoleHandler 职位 Handler
type RoleHandler struct {
	service *application.RoleService
}

// NewRoleHandler 创建 Handler
func NewRoleHandler(service *application.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

// Create 创建职位
// @Summary      创建职位
// @Description  创建新职位
// @Tags         职位管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateRoleRequest true "职位信息"
// @Success      200 {object} application.RoleResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user/roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req application.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Create(&req)
	if err != nil {
		if errors.Is(err, domain.ErrRoleCodeDuplicate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "职位编码已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Update 更新职位
// @Summary      更新职位
// @Description  更新职位信息
// @Tags         职位管理
// @Accept       json
// @Produce      json
// @Param        id path int true "职位ID"
// @Param        request body application.UpdateRoleRequest true "职位信息"
// @Success      200 {object} application.RoleResponse "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      404 {object} map[string]string "职位不存在"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user/roles/{id} [post]
func (h *RoleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Update(uint(id), &req)
	if err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "职位不存在"})
			return
		}
		if errors.Is(err, domain.ErrRoleCodeDuplicate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "职位编码已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Get 获取职位详情
// @Summary      获取职位详情
// @Description  根据职位ID获取职位详细信息
// @Tags         职位管理
// @Accept       json
// @Produce      json
// @Param        id path int true "职位ID"
// @Success      200 {object} application.RoleResponse "获取成功"
// @Failure      404 {object} map[string]string "职位不存在"
// @Security     Bearer
// @Router       /user/roles/{id} [get]
func (h *RoleHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "职位不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List 获取职位列表
// @Summary      获取职位列表
// @Description  获取职位列表（支持状态筛选和分页）
// @Tags         职位管理
// @Accept       json
// @Produce      json
// @Param        status query string false "职位状态（active/inactive）"
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Success      200 {object} application.RoleListResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user/roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	resp, err := h.service.List(statusPtr, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Delete 删除职位（软删除）
// @Summary      删除职位
// @Description  软删除职位
// @Tags         职位管理
// @Accept       json
// @Produce      json
// @Param        id path int true "职位ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      404 {object} map[string]string "职位不存在"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user/roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.Deactivate(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "职位不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetRoutes 获取路由定义
func (h *RoleHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{
			Method:      "POST",
			Path:        "/role",
			Handler:     h.Create,
			Middlewares: []gin.HandlerFunc{audit.Mark("role", "create")},
			Name:        "创建职位",
		},
		{
			Method:      "GET",
			Path:        "/role/:id",
			Handler:     h.Get,
			Middlewares: nil,
			Name:        "获取职位详情",
		},
		{
			Method:      "GET",
			Path:        "/role",
			Handler:     h.List,
			Middlewares: nil,
			Name:        "获取职位列表",
		},
		{
			Method:      "PUT",
			Path:        "/role/:id",
			Handler:     h.Update,
			Middlewares: []gin.HandlerFunc{audit.Mark("role", "update")},
			Name:        "更新职位信息",
		},
		{
			Method:      "DELETE",
			Path:        "/role/:id",
			Handler:     h.Delete,
			Middlewares: []gin.HandlerFunc{audit.Mark("role", "delete")},
			Name:        "删除职位",
		},
	}
}

