package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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
// @Router       /role [post]
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
// @Router       /role/{id} [put]
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
// @Router       /role/{id} [get]
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
// @Router       /role [get]
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

// Deactivate 停用职位
// @Summary      停用职位
// @Description  停用职位（软删除）
// @Tags         职位管理
// @Accept       json
// @Produce      json
// @Param        id path int true "职位ID"
// @Success      200 {object} map[string]string "停用成功"
// @Failure      404 {object} map[string]string "职位不存在"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /role/{id}/deactivate [put]
func (h *RoleHandler) Deactivate(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "停用成功"})
}

// Activate 激活职位
// @Summary      激活职位
// @Description  激活已停用的职位
// @Tags         职位管理
// @Accept       json
// @Produce      json
// @Param        id path int true "职位ID"
// @Success      200 {object} map[string]string "激活成功"
// @Failure      404 {object} map[string]string "职位不存在"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /role/{id}/activate [put]
func (h *RoleHandler) Activate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.Activate(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "职位不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "激活成功"})
}

// RegisterRoleHandlers 注册职位路由
func RegisterRoleHandlers(rg *gin.RouterGroup, service *application.RoleService) {
	handler := NewRoleHandler(service)

	rg.POST("/role", handler.Create)
	rg.GET("/role/:id", handler.Get)
	rg.GET("/role", handler.List)
	rg.PUT("/role/:id", handler.Update)
	rg.PUT("/role/:id/deactivate", handler.Deactivate)
	rg.PUT("/role/:id/activate", handler.Activate)
}

