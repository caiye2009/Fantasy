package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/user_role/application"
	"back/internal/user_role/domain"
	"back/pkg/endpoint"
)

// UserRoleHandler 用户角色 Handler
type UserRoleHandler struct {
	service *application.UserRoleService
}

// NewUserRoleHandler 创建 Handler
func NewUserRoleHandler(service *application.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{service: service}
}

// Create 创建角色
// @Summary      创建角色
// @Description  创建新的用户角色
// @Tags         用户角色
// @Accept       json
// @Produce      json
// @Param        request body application.CreateUserRoleRequest true "创建角色请求"
// @Success      200 {object} application.UserRoleResponse
// @Failure      400 {object} map[string]string
// @Failure      409 {object} map[string]string "角色代码已存在"
// @Failure      500 {object} map[string]string
// @Router       /user_roles [post]
// @Security     BearerAuth
func (h *UserRoleHandler) Create(c *gin.Context) {
	var req application.CreateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrUserRoleCodeExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "role_code already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Get 获取单个角色
// @Summary      获取角色详情
// @Description  根据ID获取角色详细信息
// @Tags         用户角色
// @Accept       json
// @Produce      json
// @Param        id path int true "角色ID"
// @Success      200 {object} application.UserRoleResponse
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user_roles/{id} [get]
// @Security     BearerAuth
func (h *UserRoleHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrUserRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List 列表
// @Summary      获取角色列表
// @Description  分页获取角色列表
// @Tags         用户角色
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} map[string]interface{} "total: 总数, list: 角色列表"
// @Failure      500 {object} map[string]string
// @Router       /user_roles [get]
// @Security     BearerAuth
func (h *UserRoleHandler) List(c *gin.Context) {
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

// Update 更新角色
// @Summary      更新角色
// @Description  更新角色信息
// @Tags         用户角色
// @Accept       json
// @Produce      json
// @Param        id path int true "角色ID"
// @Param        request body map[string]interface{} true "更新字段"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user_roles/{id} [put]
// @Security     BearerAuth
func (h *UserRoleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrUserRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete 删除角色
// @Summary      删除角色
// @Description  软删除角色
// @Tags         用户角色
// @Accept       json
// @Produce      json
// @Param        id path int true "角色ID"
// @Success      200 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user_roles/{id} [delete]
// @Security     BearerAuth
func (h *UserRoleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrUserRoleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

// GetRoutes 获取路由定义
func (h *UserRoleHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/user_roles", Handler: h.Create, Domain: "user_role", Action: "create"},
		{Method: "GET", Path: "/user_roles/:id", Handler: h.Get, Domain: "user_role", Action: "get"},
		{Method: "GET", Path: "/user_roles", Handler: h.List, Domain: "user_role", Action: "list"},
		{Method: "PUT", Path: "/user_roles/:id", Handler: h.Update, Domain: "user_role", Action: "update"},
		{Method: "DELETE", Path: "/user_roles/:id", Handler: h.Delete, Domain: "user_role", Action: "delete"},
	}
}
