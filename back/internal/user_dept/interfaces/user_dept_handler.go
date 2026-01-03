package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/user_dept/application"
	"back/internal/user_dept/domain"
	"back/pkg/endpoint"
)

// UserDeptHandler 用户部门 Handler
type UserDeptHandler struct {
	service *application.UserDeptService
}

// NewUserDeptHandler 创建 Handler
func NewUserDeptHandler(service *application.UserDeptService) *UserDeptHandler {
	return &UserDeptHandler{service: service}
}

// Create 创建部门
// @Summary      创建部门
// @Description  创建新的用户部门
// @Tags         用户部门
// @Accept       json
// @Produce      json
// @Param        request body application.CreateUserDeptRequest true "创建部门请求"
// @Success      200 {object} application.UserDeptResponse
// @Failure      400 {object} map[string]string
// @Failure      409 {object} map[string]string "部门代码已存在"
// @Failure      500 {object} map[string]string
// @Router       /user_depts [post]
// @Security     BearerAuth
func (h *UserDeptHandler) Create(c *gin.Context) {
	var req application.CreateUserDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrUserDeptCodeExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "dept_code already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Get 获取单个部门
// @Summary      获取部门详情
// @Description  根据ID获取部门详细信息
// @Tags         用户部门
// @Accept       json
// @Produce      json
// @Param        id path int true "部门ID"
// @Success      200 {object} application.UserDeptResponse
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user_depts/{id} [get]
// @Security     BearerAuth
func (h *UserDeptHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrUserDeptNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_dept not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List 列表
// @Summary      获取部门列表
// @Description  分页获取部门列表
// @Tags         用户部门
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} map[string]interface{} "total: 总数, list: 部门列表"
// @Failure      500 {object} map[string]string
// @Router       /user_depts [get]
// @Security     BearerAuth
func (h *UserDeptHandler) List(c *gin.Context) {
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

// Update 更新部门
// @Summary      更新部门
// @Description  更新部门信息
// @Tags         用户部门
// @Accept       json
// @Produce      json
// @Param        id path int true "部门ID"
// @Param        request body map[string]interface{} true "更新字段"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user_depts/{id} [put]
// @Security     BearerAuth
func (h *UserDeptHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrUserDeptNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_dept not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete 删除部门
// @Summary      删除部门
// @Description  软删除部门
// @Tags         用户部门
// @Accept       json
// @Produce      json
// @Param        id path int true "部门ID"
// @Success      200 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /user_depts/{id} [delete]
// @Security     BearerAuth
func (h *UserDeptHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrUserDeptNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_dept not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

// GetRoutes 获取路由定义
func (h *UserDeptHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/user_depts", Handler: h.Create, Domain: "user_dept", Action: "create"},
		{Method: "GET", Path: "/user_depts/:id", Handler: h.Get, Domain: "user_dept", Action: "get"},
		{Method: "GET", Path: "/user_depts", Handler: h.List, Domain: "user_dept", Action: "list"},
		{Method: "PUT", Path: "/user_depts/:id", Handler: h.Update, Domain: "user_dept", Action: "update"},
		{Method: "DELETE", Path: "/user_depts/:id", Handler: h.Delete, Domain: "user_dept", Action: "delete"},
	}
}
